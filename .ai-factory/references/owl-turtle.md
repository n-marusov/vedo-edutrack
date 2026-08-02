# OWL 2 in Turtle (RDF 1.1) Reference

> Source:
> - https://www.w3.org/TR/turtle/ (RDF 1.1 Turtle, W3C Rec 2014-02-25)
> - https://www.w3.org/TR/owl2-syntax/ (OWL 2 Structural Specification & Functional-Style Syntax, 2nd ed., W3C Rec 2012-12-11)
> - https://www.w3.org/TR/owl2-mapping-to-rdf/ (OWL 2 Mapping to RDF Graphs, 2nd ed., W3C Rec 2012-12-11)
>
> Created: 2026-08-02

## Overview

This reference is a practical grammar + mapping handbook for writing **OWL 2 DL ontologies serialized in Turtle** (`.ttl`), the exact format used by `traceability.ttl`. It combines:

1. The complete **Turtle** concrete syntax (directives, triples, literals, blank nodes, collections, escapes).
2. The authoritative **OWL 2 → RDF triple mapping** — the exact triples every OWL 2 construct produces, so a TBox can be written by hand without parser errors or DL-typing violations.

Rule of thumb: *Turtle answers "how do I write this syntactically?", the mapping answers "what triples does this OWL axiom need?"*. A `.ttl` ontology is valid iff it parses as Turtle AND the resulting RDF graph is a valid OWL 2 DL ontology (typing constraints + global restrictions).

## Core Concepts

- **RDF graph** = set of triples `(subject, predicate, object)`. Turtle is a text syntax for an RDF graph.
- **RDF terms**: IRIs, literals (lexical form + optional language tag OR datatype IRI), blank nodes.
- **OWL 2 ontology** = ontology IRI (+ optional version IRI) + annotations + axioms. In RDF it is a graph whose header node has `rdf:type owl:Ontology`.
- **TBox** = class/property definitions and axioms (schema). **ABox** = instances/assertions. `traceability.ttl` is TBox-only.
- **OWL 2 DL typing constraints**: every IRI used as class / object property / data property / annotation property must be declared with the right type; no IRI may be both class and datatype, nor two kinds of property. `owl:TransitiveProperty`, `owl:FunctionalProperty` etc. *imply* `owl:ObjectProperty` / `owl:DatatypeProperty` typing when parsed (OWL 1 DL compat, mapping Table 6).
- **Literal default datatype**: a literal without `^^` and without `@lang` has datatype `xsd:string`. With `@lang` it has datatype `rdf:langString`.
- **`a` keyword** = `rdf:type` (the only allowed abbreviation).

## Turtle Syntax (RDF 1.1)

### Document structure

```
turtleDoc   ::= statement*
statement   ::= directive | triples '.'
directive   ::= @prefix ID <IRI> . | @base <IRI> . | PREFIX pn: <IRI> | BASE <IRI>
triples     ::= subject predicateObjectList | blankNodePropertyList predicateObjectList?
predicateObjectList ::= verb objectList (';' (verb objectList)?)*
objectList  ::= object (',' object)*
```

- Every **triple statement ends with `.`**; `.` also terminates `@prefix`/`@base` directives.
- **`;`** repeats the subject: `s p1 o1 ; p2 o2 .` ≡ two triples with subject `s`.
- **`,`** repeats subject+predicate: `s p o1 , o2 .` ≡ two triples.
- Multiple `rdf:type` on one subject: `tr:x a owl:ObjectProperty, owl:TransitiveProperty .`

### Directives

| Form | Case | Trailing dot | Example |
|------|------|--------------|---------|
| `@prefix` | case-**sensitive** | **required** | `@prefix tr: <https://vedo-edutrack.dev/traceability#> .` |
| `@base` | case-sensitive | **required** | `@base <http://example.org/> .` |
| `PREFIX` / `BASE` | case-insensitive | **forbidden** | `PREFIX tr: <https://.../>` |

- `@prefix` must be followed by whitespace then `label:` (label may be empty → default prefix `:`).
- Prefix labels: letters/digits/`-`/`.` etc. (PN_PREFIX); local part (after `:`) may contain digits at start, `.`, `-`, `%`-escapes, reserved-char escapes (`\~.\-!$&'()*+,;=/?#@%_`).
- `PNAME_LN` local parts allow leading digits: `tr:hasArtifactId`, `:F1` are fine. `tr:` alone (no local part) is also a valid IRI term (the namespace itself).

### IRIs

- Absolute: `<https://vedo-edutrack.dev/traceability#Artifact>`.
- Prefixed: `tr:Artifact` = prefix IRI + local part.
- Relative IRIs `<#foo>` resolve against `@base` (avoid in ontologies; use absolute or prefixed).
- `a` in predicate position = `rdf:type`.
- IRIREF cannot contain raw `<>"{}|^\` or control chars / spaces; use `\uXXXX` / `\UXXXXXXXX` escapes if needed.
- `#` inside `<...>` is part of the IRI, **not** a comment.

### Literals

```
RDFLiteral ::= String (LANGTAG | '^^' iri)?      -- no whitespace between parts!
String      ::= "..." | '...' | """...""" | '''...'''
LANGTAG     ::= @[a-zA-Z]+ ('-'[a-zA-Z0-9]+)*
```

- `"label"@en` — language tag; `"text"^^xsd:string` — datatype; `"text"` — plain (xsd:string).
- **No whitespace allowed** between the closing quote and `@lang` or `^^datatype`: `"x" @en` and `"x" ^^ xsd:string` are **syntax errors**.
- `@en` and `^^` are mutually exclusive on the same literal.
- Triple-quoted `"""..."""` literals may span lines and contain embedded quotes (but not `"""`).
- Plain quotes may not contain `LF`/`CR` or the delimiter char; use escapes.

**Numbers (shorthand, case matters):**

| Datatype | Shorthand | Regex | Full form |
|----------|-----------|-------|-----------|
| `xsd:integer` | `-5` | `[+-]?[0-9]+` | `"-5"^^xsd:integer` |
| `xsd:decimal` | `-5.0` | `[+-]?[0-9]*\.[0-9]+` | `"-5.0"^^xsd:decimal` |
| `xsd:double` | `4.2E9` | mantissa + `[eE][+-]?[0-9]+` | `"4.2E9"^^xsd:double` |
| `xsd:boolean` | `true` / `false` | — | `"true"^^xsd:boolean` |

- Booleans are **case-sensitive** in Turtle (`True` is invalid).
- No whitespace between sign and digits: `- 5` invalid.
- Cardinality values in OWL restrictions canonically use `"n"^^xsd:nonNegativeInteger`; bare integers are accepted by parsers (the reverse mapping matches any literal whose value is a non-negative integer).

### Blank nodes

- Labeled: `_:node1` — `_:` followed by letters/digits/`_`/`-`/`.` (`.` not first/last). Same label = same node in the document.
- Anonymous property list (the workhorse for OWL restrictions):
  ```
  tr:Artifact rdfs:subClassOf [ a owl:Restriction ; owl:onProperty tr:hasArtifactId ; owl:cardinality 1 ] .
  ```
  `[ ... ]` in subject or object position creates a fresh blank node.
- Anonymous node in subject position: `[] a owl:AllDisjointClasses ; owl:members ( ... ) .`

### Collections (RDF lists)

- `( term1 term2 ... )` in subject/object position = an `rdf:first`/`rdf:rest` list chain ending `rdf:nil`; `()` = `rdf:nil`.
- Used by OWL for `owl:members`, `owl:oneOf`, `owl:unionOf`, `owl:intersectionOf`, `owl:withRestrictions`, `owl:hasKey`, `owl:propertyChainAxiom`.
- Lists can be nested and can contain blank-node property lists.

### Comments and whitespace

- `#` starts a comment outside IRIs/strings, to end of line (or EOF). Comments are treated as whitespace.
- Whitespace = space, tab, CR, LF. Any amount between tokens; required where two tokens would otherwise merge.
- Encoding is always **UTF-8**; media type `text/turtle`, extension `.ttl`.

### Escape sequences

| Kind | Sequences | Allowed in |
|------|-----------|------------|
| Numeric | `\uXXXX` (4 hex), `\UXXXXXXXX` (8 hex) | IRIs, strings |
| String | `\t \b \n \r \f \" \' \\` | strings only |
| Reserved char | `\` + `~.-!$&'()*+,;=/?#@%_` | local names only |

## OWL 2 → Turtle Mapping (authoritative triple patterns)

Notation: `*:x` IRI, `_:x` blank node, `lt` literal. "Main node" of an expression is what the containing triple points to.

### Ontology header

| Construct | Triples |
|-----------|---------|
| Ontology with IRI | `<ontoIRI> a owl:Ontology .` (+ optional `<ontoIRI> owl:versionIRI <verIRI> .`) |
| Imports | `<ontoIRI> owl:imports <other> .` |
| Ontology annotations | `<ontoIRI> rdfs:label "..."@en .` etc. (plain triples on the ontology IRI) |

`traceability.ttl` uses: `tr: a owl:Ontology ; rdfs:label "..."@en ; rdfs:comment "..."@en ; owl:versionInfo "0.1.0" .`

### Declarations (typing)

| Declaration | Triple |
|-------------|--------|
| Class( C ) | `C a owl:Class .` |
| Datatype( DT ) | `DT a rdfs:Datatype .` |
| ObjectProperty( OP ) | `OP a owl:ObjectProperty .` |
| DataProperty( DP ) | `DP a owl:DatatypeProperty .` |
| AnnotationProperty( AP ) | `AP a owl:AnnotationProperty .` |
| NamedIndividual( a ) | `a a owl:NamedIndividual .` |

Property characteristics **imply** the base declaration:
`P a owl:TransitiveProperty .` ⇒ also ObjectProperty; `P a owl:FunctionalProperty .` ⇒ ObjectProperty **or** DatatypeProperty depending on usage context (parser resolves via other declarations / usage).

### Class axioms

| Axiom | Triples |
|-------|---------|
| SubClassOf( C1 C2 ) | `C1 rdfs:subClassOf C2 .` |
| EquivalentClasses( C1 ... Cn ) | `C1 owl:equivalentClass C2 . ... C(n-1) owl:equivalentClass Cn .` |
| DisjointClasses( C1 C2 ) | `C1 owl:disjointWith C2 .` |
| DisjointClasses( C1 ... Cn ), n>2 | `_:x a owl:AllDisjointClasses . _:x owl:members ( C1 ... Cn ) .` |
| DisjointUnion( C CE1 ... CEn ) | `C owl:disjointUnionOf ( CE1 ... CEn ) .` |
| ClassAssertion( CE a ) | `a rdf:type CE .` (predicate is the class expression) |

### Object property axioms

| Axiom | Triples |
|-------|---------|
| SubObjectPropertyOf( P1 P2 ) | `P1 rdfs:subPropertyOf P2 .` |
| InverseObjectProperties( P1 P2 ) | `P1 owl:inverseOf P2 .` |
| ObjectPropertyDomain( P CE ) | `P rdfs:domain CE .` |
| ObjectPropertyRange( P CE ) | `P rdfs:range CE .` |
| FunctionalObjectProperty( P ) | `P a owl:FunctionalProperty .` |
| InverseFunctionalObjectProperty( P ) | `P a owl:InverseFunctionalProperty .` |
| Reflexive / Irreflexive / Symmetric / Asymmetric / Transitive | `P a owl:ReflexiveProperty .` / `owl:IrreflexiveProperty` / `owl:SymmetricProperty` / `owl:AsymmetricProperty` / `owl:TransitiveProperty` |
| SubObjectPropertyOf( chain, P ) | `P owl:propertyChainAxiom ( P1 ... Pn ) .` |
| DisjointObjectProperties( P1 P2 ) | `P1 owl:propertyDisjointWith P2 .` |
| DisjointObjectProperties( n>2 ) | `_:x a owl:AllDisjointProperties . _:x owl:members ( ... ) .` |
| ObjectInverseOf( P ) | `_:x owl:inverseOf P .` |
| ObjectPropertyAssertion( P a1 a2 ) | `a1 P a2 .` |

### Data property axioms

| Axiom | Triples |
|-------|---------|
| SubDataPropertyOf( P1 P2 ) | `P1 rdfs:subPropertyOf P2 .` |
| DataPropertyDomain( P CE ) | `P rdfs:domain CE .` |
| DataPropertyRange( P DR ) | `P rdfs:range DR .` |
| FunctionalDataProperty( P ) | `P a owl:FunctionalProperty .` |
| DataPropertyAssertion( P a lt ) | `a P lt .` |

### Restrictions (class expressions, blank-node based)

Common skeleton (unqualified):
```
_:x a owl:Restriction ;
    owl:onProperty P ;
    owl:someValuesFrom C .        # or owl:allValuesFrom / owl:hasValue
```
Cardinality (unqualified): `owl:cardinality "1"^^xsd:nonNegativeInteger`, `owl:minCardinality`, `owl:maxCardinality`.
Qualified cardinality adds a qualifier triple:
- `owl:minQualifiedCardinality` / `owl:maxQualifiedCardinality` / `owl:qualifiedCardinality` + `owl:onClass C` (object) or `owl:onDataRange DR` (data).

| Functional form | Triples |
|-----------------|---------|
| ObjectSomeValuesFrom( OPE CE ) | `_:x a owl:Restriction ; owl:onProperty OPE ; owl:someValuesFrom CE .` |
| ObjectAllValuesFrom( OPE CE ) | `... owl:allValuesFrom CE .` |
| ObjectHasValue( OPE a ) | `... owl:hasValue a .` |
| ObjectHasSelf( OPE ) | `... owl:hasSelf "true"^^xsd:boolean .` |
| ObjectMinCardinality( n OPE ) | `... owl:minCardinality "n"^^xsd:nonNegativeInteger .` |
| ObjectMaxCardinality( n OPE ) | `... owl:maxCardinality "n"^^xsd:nonNegativeInteger .` |
| ObjectExactCardinality( n OPE ) | `... owl:cardinality "n"^^xsd:nonNegativeInteger .` |
| DataSomeValuesFrom( DPE DR ) | `... owl:someValuesFrom DR .` (DR is rdfs:Datatype or restriction node) |
| DataHasValue( DPE lt ) | `... owl:hasValue lt .` |

Usage in subClassOf (as in traceability.ttl):
```
tr:Artifact rdfs:subClassOf
  [ a owl:Restriction ; owl:onProperty tr:hasArtifactId ; owl:cardinality 1 ] .
```

### Assertions about individuals

| Axiom | Triples |
|-------|---------|
| SameIndividual( a1 ... an ) | `a1 owl:sameAs a2 . ... a(n-1) owl:sameAs an .` |
| DifferentIndividuals( a1 a2 ) | `a1 owl:differentFrom a2 .` |
| DifferentIndividuals( a1 ... an ), n>2 | `_:x a owl:AllDifferent . _:x owl:members ( a1 ... an ) .` |
| NegativeObjectPropertyAssertion | `_:x a owl:NegativePropertyAssertion ; owl:sourceIndividual a1 ; owl:assertionProperty P ; owl:targetIndividual a2 .` |
| NegativeDataPropertyAssertion | same with `owl:targetValue lt .` |

### Annotations (non-logical)

- Plain: `IRI rdfs:label "..."@en .`, `IRI rdfs:comment "..."@en .`, `IRI owl:versionInfo "..." .`
- Built-in annotation properties: `rdfs:label`, `rdfs:comment`, `rdfs:seeAlso`, `rdfs:isDefinedBy`, `owl:versionInfo`, `owl:deprecated`, `owl:priorVersion`, `owl:backwardCompatibleWith`, `owl:incompatibleWith`.
- Annotated axioms are reified: keep the main triple, plus
  ```
  _:x a owl:Axiom ;
      owl:annotatedSource <s> ;
      owl:annotatedProperty <p> ;
      owl:annotatedTarget <o> ;
      rdfs:comment "..." .
  ```
  (For blank-node axioms — AllDisjointClasses etc. — annotations attach directly to the existing blank node; no reification needed.)

### RDF lists used by OWL

`( x1 ... xn )` expands to `_:b0 rdf:first x1 ; rdf:rest _:b1 . ... _:bn rdf:rest rdf:nil .`
Required list-valued properties: `owl:members`, `owl:distinctMembers` (OWL 1), `owl:oneOf`, `owl:unionOf`, `owl:intersectionOf`, `owl:withRestrictions`, `owl:hasKey`, `owl:propertyChainAxiom`, `owl:disjointUnionOf`.

## OWL 2 DL constraints that matter for traceability.ttl

- **Ontology header**: exactly one `x rdf:type owl:Ontology` node must be matched (with optional `owl:versionIRI`); the ontology IRI must not be from the reserved vocabulary.
- **Typing (Section 5.8.1)**: object / data / annotation property sets are pairwise disjoint; class vs datatype disjoint. Declare every entity you use. `tr:BusinessFunction` / `tr:BusinessProblem` are plain classes — fine even without `rdfs:subClassOf tr:Artifact`.
- **Simple roles (Section 11.2)**: transitive properties (`owl:TransitiveProperty`), their super-properties, and chain-defined properties must **not** appear in `ObjectMin/Max/ExactCardinality`, `ObjectHasSelf`, `FunctionalObjectProperty`, `InverseFunctionalObjectProperty`, `Irreflexive`, `Asymmetric`, or `DisjointObjectProperties`. `traceability.ttl` is safe: `tr:traceableTo` is transitive but only used as super-property and never in restrictions; `owl:cardinality 1` is on the *datatype* property `tr:hasArtifactId` (data properties have no such restriction).
- **Anonymous individuals**: cannot occur in `SameIndividual`, `DifferentIndividuals`, negative assertions, `ObjectOneOf`, `ObjectHasValue`; the anonymous-individual assertion graph must be a forest.
- **Datatypes**: OWL 2 DL allows only the built-in datatype map (`xsd:*` listed in Section 4), `rdfs:Literal`, or datatypes defined via `DatatypeDefinition`. `xsd:string`, `xsd:dateTime`, `xsd:nonNegativeInteger`, `xsd:boolean` are all fine.
- **Declaration consistency** is optional (an ontology MAY use undeclared individuals), but declared types must never contradict usage.
- Axiom sets must not contain structurally equivalent duplicates (sets, not bags).

## Validation checklist (run before committing generated .ttl)

1. **Turtle parse** — `riot --validate traceability.ttl` (Apache Jena), or `python -c "import rdflib; rdflib.Graph().parse('traceability.ttl', format='turtle')"`.
2. **OWL structural check** — `robot validate --input traceability.ttl` (ROBOT), or load into Protégé and run the reasoner (HermiT/ELK).
3. **Manual spot-checks:**
   - Every statement ends with `.`; `;` between predicate lists; `,` between objects.
   - `@prefix ... .` has the trailing dot; no trailing dot after `PREFIX ...` form.
   - No `"text" @en` or `"text" ^^xsd:string` whitespace mistakes.
   - `owl:members` lists have ≥ 2 entries (n>2 rules) and are well-formed `rdf:first/rest/nil` chains.
   - All IRIs referenced in axioms are declared (`a owl:Class`, `a owl:ObjectProperty`, `a owl:DatatypeProperty`).
   - Cardinality literals are `"n"^^xsd:nonNegativeInteger` (bare integers also accepted by most parsers, but canonical form is safer).
   - A restriction blank node contains `owl:onProperty` plus exactly one of `someValuesFrom` / `allValuesFrom` / `hasValue` / `cardinality` / `minCardinality` / `maxCardinality` / qualified variants.
   - Property characteristics and domains/ranges stay on the correct property-kind (object vs data).

## Best Practices

1. Write `@prefix` (not `PREFIX`) for maximal tool compatibility, with trailing dots.
2. Declare every class and property before/alongside its axioms — helps parsers disambiguate object vs data properties (critical for `owl:FunctionalProperty` on `tr:hasArtifactId`).
3. Keep one IRI scheme per namespace: `tr:` prefix ends with `#` (or `/`) and local parts are CamelCase / dotted identifiers.
4. Put `rdfs:label` and `rdfs:comment` (with `@en` tags) on every class and property — cheap and machine-readable documentation.
5. Use the `[ a owl:Restriction ; ... ]` inline blank-node form for subClassOf constraints; it stays readable and is the canonical mapping.
6. For disjointness of >2 classes always use the single `_:x a owl:AllDisjointClasses ; owl:members (...)` form — one axiom, clearer than pairwise `owl:disjointWith`.
7. Keep the file UTF-8; no BOM; LF or CRLF both parse, but pick one.
8. Group the file logically (header → classes → properties → axioms) with section comments; comments are free (treated as whitespace).
9. Version the ontology via `owl:versionInfo` (and optionally `owl:versionIRI`).
10. After editing, run the validation checklist — cheap, catches 99% of generation errors.

## Common Pitfalls

- **Missing trailing `.`** after `@prefix`/`@base` — the single most common Turtle error.
- **`"text" @en`** — whitespace before `@`/`^^` makes the whole statement a parse error.
- **`true`/`True`/`TRUE`** — booleans are lowercase-only in Turtle (`TrUe` invalid).
- **Using `rdfs:Class` instead of `owl:Class`** for declarations — `owl:Class` is the OWL 2 DL declaration; `rdfs:Class` alone is not recognized as a class declaration by the mapping (though OWL 1 compat removes redundant `rdfs:Class` triples when `owl:Class` is present).
- **Declaring a datatype with `a owl:Class`** — datatypes must be `rdfs:Datatype`; class/datatype typing is disjoint in OWL 2 DL.
- **`owl:AllDisjointClasses` with 1 member** — reverse mapping requires n ≥ 2; pairwise `owl:disjointWith` is for exactly 2.
- **Typing `owl:FunctionalProperty` on a property that is used with both IRI and literal objects** — the parser must be able to pick object vs data property; ambiguous usage can make the ontology non-DL.
- **Using a transitive property in a cardinality restriction** — violates the simple-role restriction (undecidable in DL); `traceability.ttl` correctly avoids it.
- **Literal without datatype but intended as date** — `"2026-08-02"` is `xsd:string`, not `xsd:dateTime`; write `"2026-08-02T00:00:00Z"^^xsd:dateTime`.
- **Unbalanced quotes/brackets in multiline literals** — `"""..."""` must not contain `"""`; check for accidental `"""` inside.
- **Reusing one blank-node label for unrelated purposes** — within a document the same `_:x` label is the *same* node; use `[]` for one-off nodes.
- **Non-ASCII / spaces inside `<IRI>`** — raw spaces, `<`, `>`, `"`, `{`, `}`, `|`, `^`, `\` are forbidden in IRIREF; escape with `\uXXXX` or use prefixed names.
- **`owl:versionIRI` without ontology IRI** — a version IRI is only allowed when an ontology IRI is present.

## Version Notes

- Turtle: RDF 1.1 (2014) added `PREFIX`/`BASE` (SPARQL-style) alongside `@prefix`/`@base`; both are normative. RDF 1.1 literals: no more "plain literal" type — `"x"@en` is `rdf:langString`; `"x"` is `xsd:string`.
- OWL 2 (2012, 2nd ed.) is the current Recommendation; OWL 1 serializations (`owl:distinctMembers`, `owl:DataRange`, reification with `owl:subject/owl:predicate/owl:object`) are recognized for backwards compatibility but should not be produced.
- OWL 2 annotations reification uses `owl:annotatedSource` / `owl:annotatedProperty` / `owl:annotatedTarget` (renamed from OWL 1's `owl:subject/object/predicate` in 2009).
