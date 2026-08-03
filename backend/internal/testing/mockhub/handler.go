package mockhub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vektah/gqlparser/v2/parser"
	"github.com/vektah/gqlparser/v2/validator"
	"go.uber.org/zap"
)

// schema is the loaded GraphQL schema (embedded SDL from schema.graphql).
var schema *ast.Schema

func init() {
	// Load the embedded SDL once at startup.
	schema, _ = gqlparser.LoadSchema(&ast.Source{
		Name:  "schema.graphql",
		Input: mustSchemaSDL(),
	})
}

// Handler is the mock VEDO Hub GraphQL HTTP handler.
type Handler struct {
	ont    *Ontology
	logger *zap.Logger
}

// NewHandler creates the GraphQL handler over the given in-memory ontology.
func NewHandler(ont *Ontology, logger *zap.Logger) http.Handler {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Handler{ont: ont, logger: logger}
}

// ServeHTTP routes POST /graphql (GraphQL execution) and GET /healthz.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/healthz":
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
		return
	case r.Method != http.MethodPost || r.URL.Path != "/graphql":
		http.Error(w, `{"error":"not_found"}`, http.StatusNotFound)
		return
	}

	// Auth: any Authorization: Bearer token is accepted; missing → GraphQL error.
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		h.writeGraphQLError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var req struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeGraphQLError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	if req.Query == "" {
		h.writeGraphQLError(w, http.StatusBadRequest, "query is required")
		return
	}

	result, err := h.execute(r.Context(), req.Query, req.Variables)
	if err != nil {
		h.logger.Error("graphql execution error", zap.Error(err))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}

// execute parses, validates and runs the query against the QueryRoot resolvers.
func (h *Handler) execute(ctx context.Context, query string, variables map[string]interface{}) (map[string]interface{}, error) {
	doc, err := parser.ParseQuery(&ast.Source{Input: query, Name: "request"})
	if err != nil {
		return nil, err
	}

	listErr := validator.ValidateWithRules(schema, doc, nil)
	if len(listErr) > 0 {
		return nil, listErr[0]
	}

	if len(doc.Operations) == 0 {
		return nil, fmt.Errorf("no operations in query")
	}
	op := doc.Operations[0]
	h.logger.Info("graphql operation", zap.String("operationName", op.Name))

	result := map[string]interface{}{}
	fields, err := h.resolveSelectionSet(ctx, op.SelectionSet, op.VariableDefinitions, variables)
	if err != nil {
		result["errors"] = []interface{}{gqlerror.Error{Message: err.Error()}}
		// GraphQL errors are delivered in the 200 response body (errors array),
		// not as an HTTP error — intentionally returning nil here.
		return result, nil //nolint:nilerr // GraphQL contract: errors in body
	}
	result["data"] = fields
	return result, nil
}

// resolveSelectionSet resolves a selection set against the QueryRoot.
func (h *Handler) resolveSelectionSet(ctx context.Context, selSet ast.SelectionSet, varDefs ast.VariableDefinitionList, variables map[string]interface{}) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	for _, sel := range selSet {
		field, ok := sel.(*ast.Field)
		if !ok {
			continue // fragments unsupported in this minimal executor
		}
		args := h.resolveArgs(field.Arguments, varDefs, variables)
		value, err := h.resolveField(ctx, field.Name, args, field.SelectionSet)
		if err != nil {
			return nil, err
		}
		out[field.Alias] = value
	}
	return out, nil
}

// resolveArgs resolves argument values (literals + variable references).
func (h *Handler) resolveArgs(argDefs ast.ArgumentList, varDefs ast.VariableDefinitionList, variables map[string]interface{}) map[string]interface{} {
	args := map[string]interface{}{}
	for _, arg := range argDefs {
		args[arg.Name] = resolveValue(arg.Value, varDefs, variables)
	}
	return args
}

// resolveField dispatches a QueryRoot field to its resolver.
//
//nolint:gocyclo // dispatch switch — one case per schema field (CC bounded by schema size)
func (h *Handler) resolveField(ctx context.Context, name string, args map[string]interface{}, selSet ast.SelectionSet) (interface{}, error) {
	switch name {
	case "classes":
		return h.resolveClasses(args, selSet)
	case "class":
		return h.resolveClass(args)
	case "classTree":
		return h.resolveClassTree(args)
	case "classAncestors":
		return h.resolveClassAncestors(args)
	case "classDescendants":
		return h.resolveClassDescendants(args)
	case "graphNeighborhood":
		return h.resolveGraphNeighborhood(args, selSet)
	case "autocompleteClasses":
		return h.resolveAutocomplete(args)
	case "properties":
		return h.resolveProperties(args, selSet)
	case "property":
		return h.resolveProperty(args)
	case "individuals", "individual":
		// TBox-only ontology: no individuals (data-driven, empty result).
		return h.resolveEmptyIndividuals(name)
	case "_service":
		return map[string]interface{}{"sdl": mustSchemaSDL()}, nil
	default:
		return nil, fmt.Errorf("unknown field %q", name)
	}
}

// resolveClasses lists classes with search + pagination.
func (h *Handler) resolveClasses(args map[string]interface{}, selSet ast.SelectionSet) (interface{}, error) {
	q, _ := args["q"].(string)
	page := intArg(args, "page", 0)
	perPage := intArg(args, "perPage", 20)
	if perPage <= 0 {
		perPage = 20
	}

	ids := h.ont.Autocomplete(q, 0) // 0 = no limit
	total := len(ids)
	start := page * perPage
	if start > total {
		start = total
	}
	end := start + perPage
	if end > total {
		end = total
	}

	items := make([]interface{}, 0, end-start)
	for _, id := range ids[start:end] {
		items = append(items, h.classSummary(id))
	}
	return map[string]interface{}{
		"items":   items,
		"total":   total,
		"page":    page,
		"perPage": perPage,
	}, nil
}

// resolveClass returns a single class (full detail).
func (h *Handler) resolveClass(args map[string]interface{}) (interface{}, error) {
	id, _ := args["classId"].(string)
	c := h.ont.Classes[id]
	if c == nil {
		return nil, nil
	}
	return h.classDetail(c), nil
}

// resolveClassTree returns root classes with lazy children.
func (h *Handler) resolveClassTree(args map[string]interface{}) (interface{}, error) {
	roots := h.ont.Roots()
	out := make([]interface{}, 0, len(roots))
	for _, root := range roots {
		out = append(out, h.treeNode(root))
	}
	return out, nil
}

// resolveClassAncestors returns the breadcrumb (root → class).
func (h *Handler) resolveClassAncestors(args map[string]interface{}) (interface{}, error) {
	id, _ := args["classId"].(string)
	bc := h.ont.Breadcrumb(id)
	out := make([]interface{}, 0, len(bc))
	for _, label := range bc {
		out = append(out, map[string]interface{}{"id": label, "label": label})
	}
	return out, nil
}

// resolveClassDescendants returns descendant tree nodes (maxDepth-bounded).
func (h *Handler) resolveClassDescendants(args map[string]interface{}) (interface{}, error) {
	id, _ := args["classId"].(string)
	maxDepth := intArg(args, "maxDepth", 10)
	children := h.ont.Classes[id].Children
	out := make([]interface{}, 0, len(children))
	for _, child := range children {
		out = append(out, h.treeNodeDepth(child, maxDepth))
	}
	return out, nil
}

// resolveGraphNeighborhood returns the node with its connected nodes/edges.
// TBox-only ontology: no individuals/edges — the node is returned without
// edges (data-driven).
func (h *Handler) resolveGraphNeighborhood(args map[string]interface{}, selSet ast.SelectionSet) (interface{}, error) {
	id, _ := args["classId"].(string)
	c := h.ont.Classes[id]
	if c == nil {
		return nil, fmt.Errorf("class %q not found", id)
	}
	neighbors := []interface{}{}
	for _, p := range c.Parents {
		if pc := h.ont.Classes[p]; pc != nil {
			neighbors = append(neighbors, h.classDetail(pc))
		}
	}
	for _, ch := range c.Children {
		if cc := h.ont.Classes[ch]; cc != nil {
			neighbors = append(neighbors, h.classDetail(cc))
		}
	}
	return map[string]interface{}{
		"node":      h.classDetail(c),
		"neighbors": neighbors,
		"edges":     []interface{}{},
	}, nil
}

// resolveAutocomplete returns class summaries matching the query.
func (h *Handler) resolveAutocomplete(args map[string]interface{}) (interface{}, error) {
	q, _ := args["q"].(string)
	limit := intArg(args, "limit", 20)
	ids := h.ont.Autocomplete(q, limit)
	out := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		out = append(out, h.classSummary(id))
	}
	return out, nil
}

// resolveProperties lists properties (TBox).
func (h *Handler) resolveProperties(args map[string]interface{}, selSet ast.SelectionSet) (interface{}, error) {
	propertyType, _ := args["propertyType"].(string)
	page := intArg(args, "page", 0)
	perPage := intArg(args, "perPage", 20)
	if perPage <= 0 {
		perPage = 20
	}

	var props []*Property
	for _, p := range h.ont.Properties {
		if propertyType == "" || strings.Contains(p.Type, propertyType) {
			props = append(props, p)
		}
	}
	total := len(props)
	start := page * perPage
	if start > total {
		start = total
	}
	end := start + perPage
	if end > total {
		end = total
	}
	items := make([]interface{}, 0, end-start)
	for _, p := range props[start:end] {
		items = append(items, h.propertyDetail(p))
	}
	return map[string]interface{}{
		"items":   items,
		"total":   total,
		"page":    page,
		"perPage": perPage,
	}, nil
}

// resolveProperty returns a single property.
func (h *Handler) resolveProperty(args map[string]interface{}) (interface{}, error) {
	id, _ := args["propertyId"].(string)
	p := h.ont.Properties[id]
	if p == nil {
		return nil, nil
	}
	return h.propertyDetail(p), nil
}

// resolveEmptyIndividuals returns an empty paginated result (TBox-only).
func (h *Handler) resolveEmptyIndividuals(name string) (interface{}, error) {
	return map[string]interface{}{
		"items":   []interface{}{},
		"total":   0,
		"page":    0,
		"perPage": 20,
	}, nil
}

// ── value builders ────────────────────────────────────────────────────────

func (h *Handler) classSummary(id string) interface{} {
	c := h.ont.Classes[id]
	return map[string]interface{}{
		"id":      c.ID,
		"label":   c.Label,
		"comment": c.Comment,
		"parents": c.Parents,
	}
}

func (h *Handler) classDetail(c *Class) interface{} {
	return map[string]interface{}{
		"id":           c.ID,
		"label":        c.Label,
		"comment":      c.Comment,
		"entityType":   "Class",
		"parents":      c.Parents,
		"children":     c.Children,
		"isAbstract":   c.IsAbstract,
		"isDeprecated": c.IsDeprecated,
	}
}

func (h *Handler) treeNode(id string) interface{} {
	c := h.ont.Classes[id]
	children := []interface{}{}
	for _, ch := range c.Children {
		children = append(children, h.treeNode(ch))
	}
	return map[string]interface{}{
		"id":       c.ID,
		"label":    c.Label,
		"children": children,
	}
}

func (h *Handler) treeNodeDepth(id string, maxDepth int) interface{} {
	c := h.ont.Classes[id]
	children := []interface{}{}
	if maxDepth > 0 {
		for _, ch := range c.Children {
			children = append(children, h.treeNodeDepth(ch, maxDepth-1))
		}
	}
	return map[string]interface{}{
		"id":       c.ID,
		"label":    c.Label,
		"children": children,
	}
}

func (h *Handler) propertyDetail(p *Property) interface{} {
	return map[string]interface{}{
		"id":         p.ID,
		"label":      p.Label,
		"comment":    p.Comment,
		"entityType": "Property",
		"domains":    p.Domains,
		"ranges":     p.Ranges,
		"functional": p.Functional,
	}
}

// ── helpers ───────────────────────────────────────────────────────────────

// resolveValue resolves an AST value (literal or variable reference).
func resolveValue(v *ast.Value, varDefs ast.VariableDefinitionList, variables map[string]interface{}) interface{} {
	switch v.Kind {
	case ast.Variable:
		return resolveVariable(v, varDefs, variables)
	case ast.IntValue, ast.FloatValue:
		if n, err := strconv.ParseFloat(v.Raw, 64); err == nil {
			return n
		}
		return v.Raw
	case ast.BooleanValue:
		return v.Raw == "true"
	case ast.NullValue:
		return nil
	case ast.EnumValue:
		return v.Raw
	case ast.ListValue:
		out := make([]interface{}, 0, len(v.Children))
		for _, item := range v.Children {
			out = append(out, resolveValue(item.Value, varDefs, variables))
		}
		return out
	default: // StringValue and others
		return v.Raw
	}
}

// resolveVariable resolves a variable reference, falling back to its declared
// default value when not supplied.
func resolveVariable(v *ast.Value, varDefs ast.VariableDefinitionList, variables map[string]interface{}) interface{} {
	if val, ok := variables[v.Raw]; ok {
		return val
	}
	for _, def := range varDefs {
		if def.Variable == v.Raw && def.DefaultValue != nil {
			return resolveValue(def.DefaultValue, varDefs, variables)
		}
	}
	return nil
}

// intArg reads an int argument with a default.
func intArg(args map[string]interface{}, name string, def int) int {
	if v, ok := args[name]; ok {
		switch n := v.(type) {
		case float64:
			return int(n)
		case int:
			return n
		case string:
			if parsed, err := strconv.Atoi(n); err == nil {
				return parsed
			}
		}
	}
	return def
}

// writeGraphQLError writes a GraphQL-style error response.
func (h *Handler) writeGraphQLError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"errors": []interface{}{map[string]interface{}{"message": message}},
	})
}

// mustSchemaSDL returns the embedded schema SDL (panic if missing).
func mustSchemaSDL() string {
	return embeddedSchemaSDL
}
