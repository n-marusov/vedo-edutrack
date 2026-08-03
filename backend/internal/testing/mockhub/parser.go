// Package mockhub implements a minimal in-memory VEDO Hub stand-in for
// dev/test/CI (ADR-DES.INFRA.mock-hub-strategy, M0.3 T21–T25).
//
// It parses an arbitrary .ttl ontology (mini-Turtle subset) and serves it via
// a GraphQL endpoint (gqlparser/v2 resolvers, T22). Test-only tooling — not
// part of the product binary.
package mockhub

import (
	"fmt"
	"io"
	"strings"
)

// TurtleParser parses the mini-Turtle subset needed for the known TBox shape:
// prefixes, `s a owl:Class | owl:ObjectProperty | owl:DatatypeProperty`,
// rdfs:label / rdfs:comment / rdfs:subClassOf / rdfs:domain / rdfs:range,
// owl:FunctionalProperty. Comments and blank nodes are ignored.
type TurtleParser struct {
	ont      *Ontology
	prefixes map[string]string
	line     int
}

// Parse reads a Turtle document and builds the in-memory Ontology.
// Errors carry the source line number.
func Parse(reader io.Reader) (*Ontology, error) {
	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read ontology: %w", err)
	}

	p := &TurtleParser{
		ont:      NewOntology(),
		prefixes: map[string]string{},
	}
	if err := p.parse(string(raw)); err != nil {
		return nil, err
	}
	return p.ont, nil
}

func (p *TurtleParser) parse(src string) error {
	// Strip comments and split into statements (semicolon-terminated).
	cleaned := p.stripComments(src)
	statements := splitStatements(cleaned)

	for _, stmt := range statements {
		p.line++
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		// @prefix declarations.
		if strings.HasPrefix(stmt, "@prefix") {
			if err := p.parsePrefix(stmt); err != nil {
				return err
			}
			continue
		}
		if err := p.parseStatement(stmt); err != nil {
			return fmt.Errorf("line %d: %w", p.line, err)
		}
	}
	return nil
}

// stripComments removes `#` comments, ignoring `#` inside <IRI> spans and
// quoted literals.
func (p *TurtleParser) stripComments(src string) string {
	var sb strings.Builder
	inIri := false
	inQuote := false
	for _, line := range strings.Split(src, "\n") {
		for _, r := range line {
			switch {
			case r == '<':
				inIri = true
				sb.WriteRune(r)
			case r == '>':
				inIri = false
				sb.WriteRune(r)
			case r == '"':
				inQuote = !inQuote
				sb.WriteRune(r)
			case r == '#' && !inIri && !inQuote:
				// comment to end of line
				sb.WriteString("\n")
				goto nextLine
			default:
				sb.WriteRune(r)
			}
		}
		sb.WriteString("\n")
	nextLine:
		inIri = false
		inQuote = false
	}
	return sb.String()
}

// parsePrefix handles `@prefix pfx: <iri> .`
func (p *TurtleParser) parsePrefix(stmt string) error {
	rest := strings.TrimPrefix(stmt, "@prefix")
	rest = strings.TrimSuffix(strings.TrimSpace(rest), ".")
	parts := strings.SplitN(rest, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("line %d: malformed prefix %q", p.line, stmt)
	}
	name := strings.TrimSpace(parts[0])
	iri := strings.TrimSpace(parts[1])
	iri = strings.Trim(iri, "<> ")
	p.prefixes[name] = iri
	return nil
}

// parseStatement handles triples and `;`-separated predicate lists.
// Supported forms:
//
//	<s> a owl:Class .
//	<s> rdfs:label "text" .
//	<s> rdfs:subClassOf <t> .
//	<s> rdfs:domain <t> . / rdfs:range <t> .
//	<s> owl:FunctionalProperty true .
func (p *TurtleParser) parseStatement(stmt string) error {
	// Skip statements containing blank nodes (e.g. owl:Restriction `[ ... ]`)
	// — the plan requires ignoring blank nodes.
	if strings.Contains(stmt, "[") || strings.Contains(stmt, "]") {
		return nil
	}

	// Split subject | predicate-object list.
	subjEnd := indexAny(stmt, " \t")
	if subjEnd < 0 {
		return fmt.Errorf("line %d: expected subject in %q", p.line, stmt)
	}
	subject := strings.TrimSpace(stmt[:subjEnd])
	rest := strings.TrimSpace(stmt[subjEnd:])
	rest = strings.TrimSuffix(rest, ".")

	// Predicate-object pairs separated by `;` (and `,` for object lists —
	// simplified: handle `;` only).
	pairs := splitByTopLevel(rest, ';')
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		if err := p.parsePredicateObject(subject, pair); err != nil {
			return fmt.Errorf("line %d: %w", p.line, err)
		}
	}
	return nil
}

// parsePredicateObject handles a single `predicate object` pair.
//
//nolint:gocyclo // predicate dispatch switch — bounded by the supported Turtle subset
func (p *TurtleParser) parsePredicateObject(subject, pair string) error {
	predEnd := strings.IndexAny(pair, " \t")
	if predEnd < 0 {
		return fmt.Errorf("expected predicate in %q", pair)
	}
	predicate := strings.TrimSpace(pair[:predEnd])
	obj := strings.TrimSpace(pair[predEnd:])

	subj, err := p.resolve(subject)
	if err != nil {
		return err
	}
	pred, err := p.resolve(predicate)
	if err != nil {
		return err
	}

	switch pred {
	case "http://www.w3.org/1999/02/22-rdf-syntax-ns#type":
		return p.applyType(subj, obj)
	case "http://www.w3.org/2000/01/rdf-schema#label":
		p.setLabel(subj, obj)
	case "http://www.w3.org/2000/01/rdf-schema#comment":
		p.setComment(subj, obj)
	case "http://www.w3.org/2000/01/rdf-schema#subClassOf":
		return p.applySubClassOf(subj, obj)
	case "http://www.w3.org/2000/01/rdf-schema#domain":
		return p.applyDomainRange(subj, obj, true)
	case "http://www.w3.org/2000/01/rdf-schema#range":
		return p.applyDomainRange(subj, obj, false)
	case "http://www.w3.org/2002/07/owl#FunctionalProperty":
		if pr := p.ont.Properties[subj]; pr != nil {
			pr.Functional = true
		}
	}
	return nil
}

// applyType handles `s a owl:Class | owl:ObjectProperty | owl:DatatypeProperty`.
func (p *TurtleParser) applyType(subj, obj string) error {
	typ, err := p.resolve(obj)
	if err != nil {
		return err
	}
	switch typ {
	case "http://www.w3.org/2002/07/owl#Class":
		p.ont.Classes[subj] = &Class{ID: subj}
	case "http://www.w3.org/2002/07/owl#ObjectProperty", "http://www.w3.org/2002/07/owl#DatatypeProperty":
		p.ont.Properties[subj] = &Property{ID: subj, Type: typ}
	}
	return nil
}

// applySubClassOf wires parent/child relationships between classes.
func (p *TurtleParser) applySubClassOf(subj, obj string) error {
	parent, err := p.resolve(obj)
	if err != nil {
		return err
	}
	if c := p.ont.Classes[subj]; c != nil {
		c.Parents = append(c.Parents, parent)
	}
	if parentClass := p.ont.Classes[parent]; parentClass != nil {
		parentClass.Children = append(parentClass.Children, subj)
	}
	return nil
}

// applyDomainRange sets rdfs:domain (isDomain=true) or rdfs:range on a property.
func (p *TurtleParser) applyDomainRange(subj, obj string, isDomain bool) error {
	iri, err := p.resolve(obj)
	if err != nil {
		return err
	}
	pr := p.ont.Properties[subj]
	if pr == nil {
		return nil
	}
	if isDomain {
		pr.Domains = append(pr.Domains, iri)
	} else {
		pr.Ranges = append(pr.Ranges, iri)
	}
	return nil
}

// setLabel assigns rdfs:label to a class or property.
func (p *TurtleParser) setLabel(subj, obj string) {
	label := trimQuoted(obj)
	if c := p.ont.Classes[subj]; c != nil {
		c.Label = label
	}
	if pr := p.ont.Properties[subj]; pr != nil {
		pr.Label = label
	}
}

// setComment assigns rdfs:comment to a class or property.
func (p *TurtleParser) setComment(subj, obj string) {
	comment := trimQuoted(obj)
	if c := p.ont.Classes[subj]; c != nil {
		c.Comment = comment
	}
	if pr := p.ont.Properties[subj]; pr != nil {
		pr.Comment = comment
	}
}

// resolve expands a prefixed name (`pfx:name` or `<iri>`) to a full IRI.
// Turtle's `a` is the rdf:type shorthand.
func (p *TurtleParser) resolve(token string) (string, error) {
	token = strings.TrimSpace(token)
	if token == "a" {
		return "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", nil
	}
	if strings.HasPrefix(token, "<") && strings.HasSuffix(token, ">") {
		return strings.Trim(token, "<>"), nil
	}
	if strings.HasPrefix(token, "\"") {
		return "", fmt.Errorf("unexpected literal in IRI position: %q", token)
	}
	parts := strings.SplitN(token, ":", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("cannot resolve %q (expected <iri> or pfx:name)", token)
	}
	prefix, ok := p.prefixes[parts[0]]
	if !ok {
		return "", fmt.Errorf("unknown prefix %q in %q", parts[0], token)
	}
	return prefix + parts[1], nil
}

// trimQuoted removes surrounding double quotes from a literal and strips a
// trailing language tag (@en, @ru…).
func trimQuoted(s string) string {
	s = strings.TrimSpace(s)
	// Strip a language tag first: "text"@en → "text"
	if idx := strings.LastIndex(s, "\"@"); idx > 0 {
		s = s[:idx+1]
	}
	// Strip surrounding quotes: "text" → text
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

// splitStatements splits a Turtle document into top-level statements
// (each ending with `.`). Simplified: splits on `.` outside <> and quotes.
func splitStatements(src string) []string {
	var out []string
	var sb strings.Builder
	inIri := false
	inQuote := false
	for _, r := range src {
		switch {
		case r == '<':
			inIri = true
			sb.WriteRune(r)
		case r == '>':
			inIri = false
			sb.WriteRune(r)
		case r == '"':
			inQuote = !inQuote
			sb.WriteRune(r)
		case r == '.' && !inIri && !inQuote:
			out = append(out, sb.String())
			sb.Reset()
		default:
			sb.WriteRune(r)
		}
	}
	return out
}

// splitByTopLevel splits a predicate-object list on a separator outside <>,
// quotes and []-bracketed blank nodes.
func splitByTopLevel(s string, sep rune) []string {
	var out []string
	var sb strings.Builder
	var state splitState
	for _, r := range s {
		state.step(r)
		if r == sep && !state.inside() {
			out = append(out, sb.String())
			sb.Reset()
			continue
		}
		sb.WriteRune(r)
	}
	if strings.TrimSpace(sb.String()) != "" {
		out = append(out, sb.String())
	}
	return out
}

// splitState tracks whether we are inside an IRI, a quoted literal or a
// []-bracketed blank node (for top-level separator detection).
type splitState struct {
	inIri     bool
	inQuote   bool
	inBracket int
}

func (s *splitState) step(r rune) {
	switch r {
	case '<':
		s.inIri = true
	case '>':
		s.inIri = false
	case '"':
		s.inQuote = !s.inQuote
	case '[':
		s.inBracket++
	case ']':
		if s.inBracket > 0 {
			s.inBracket--
		}
	}
}

func (s *splitState) inside() bool {
	return s.inIri || s.inQuote || s.inBracket > 0
}

// indexAny returns the index of the first whitespace character.
func indexAny(s string, chars string) int {
	for i, r := range s {
		for _, c := range chars {
			if r == c {
				return i
			}
		}
	}
	return -1
}
