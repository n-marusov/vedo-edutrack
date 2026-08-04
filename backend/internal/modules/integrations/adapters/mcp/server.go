// Package mcp implements the MCP (Model Context Protocol) server over stdio
// (F6.6, ADR-DES.API.cli-interface: MCP runs as `vedo-edutrack mcp` over
// stdio).
//
// The server speaks JSON-RPC 2.0 over stdin/stdout: `initialize`,
// `tools/list`, `tools/call`. It is read-oriented by design (ADR-§7): tools
// expose route/progress/coverage/gap/resource reads to AI agents. Logs go to
// stderr — stdout is reserved for the protocol.
//
// The adapter is tool-agnostic: concrete tool handlers (backed by the
// application layers of other bounded contexts) are injected at the
// composition root (cli/mcp.go), keeping the module boundary intact.
package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"

	"go.uber.org/zap"
)

// Protocol version supported by this server.
const ProtocolVersion = "2025-06-18"

// Tool is one MCP tool definition plus its handler.
type Tool struct {
	// Name is the stable tool id (e.g. "get_route").
	Name string
	// Description is shown to the agent.
	Description string
	// InputSchema is a JSON Schema for the tool arguments.
	InputSchema map[string]any
	// Handler invokes the tool with parsed arguments.
	Handler func(ctx context.Context, args map[string]any) (any, error)
}

// Server is the MCP JSON-RPC server over a stdio pair.
type Server struct {
	in       io.Reader
	out      io.Writer
	tools    map[string]Tool
	apiKey   string
	logger   *zap.Logger
	mu       sync.Mutex
	writeMu  sync.Mutex
	client   string
	shutdown func()
	closed   bool
}

// Config configures the MCP server.
type Config struct {
	// APIKey is the required init key ("" disables auth).
	APIKey string
	// Tools is the tool registry.
	Tools []Tool
	// OnShutdown is invoked when stdin closes (graceful shutdown hook).
	OnShutdown func()
}

// New builds an MCP server over the given stdio pair.
func New(cfg Config, in io.Reader, out io.Writer, logger *zap.Logger) *Server {
	if logger == nil {
		logger = zap.NewNop()
	}
	tools := map[string]Tool{}
	for _, t := range cfg.Tools {
		tools[t.Name] = t
	}
	return &Server{
		in:       in,
		out:      out,
		tools:    tools,
		apiKey:   cfg.APIKey,
		logger:   logger.Named("mcp"),
		shutdown: cfg.OnShutdown,
	}
}

// Serve reads JSON-RPC requests from stdin until EOF, then shuts down.
func (s *Server) Serve(ctx context.Context) error {
	scanner := bufio.NewScanner(s.in)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	s.logger.Info("session started")

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		if err := s.handleLine(ctx, line); err != nil {
			s.logger.Error("handle request", zap.Error(err))
		}
	}
	if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
		s.logger.Error("read stdin", zap.Error(err))
	}
	s.logger.Info("session ended (stdin closed)")
	s.triggerShutdown()
	return nil
}

// handleLine parses one JSON-RPC request and writes the response.
func (s *Server) handleLine(ctx context.Context, line string) error {
	var req rpcRequest
	if err := json.Unmarshal([]byte(line), &req); err != nil {
		s.writeError(nil, -32700, "Parse error: "+err.Error())
		return err
	}
	if req.Method == "" {
		s.writeError(&req.ID, -32600, "Invalid Request: missing method")
		return nil
	}

	// Auth: every request except `initialize` requires a valid session key.
	if req.Method != "initialize" && !s.authenticated() {
		s.writeError(&req.ID, -32002, "Not authenticated: call initialize with the API key first")
		return nil
	}

	switch req.Method {
	case "initialize":
		s.handleInitialize(&req)
	case "tools/list":
		s.handleToolsList(&req)
	case "tools/call":
		s.handleToolsCall(ctx, &req)
	default:
		s.writeError(&req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method))
	}
	return nil
}

// authenticated reports whether the session passed initialize with a valid key.
func (s *Server) authenticated() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.client != ""
}

// handleInitialize processes the MCP initialize handshake (auth gate).
func (s *Server) handleInitialize(req *rpcRequest) {
	var params struct {
		ClientInfo struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"clientInfo"`
		ProtocolVersion string `json:"protocolVersion"`
		Capabilities    struct {
			APIKey string `json:"apiKey"`
		} `json:"capabilities"`
	}
	_ = unmarshalParams(req, &params)

	if s.apiKey != "" && params.Capabilities.APIKey != s.apiKey {
		s.writeError(&req.ID, -32002, "Invalid API key")
		return
	}

	s.mu.Lock()
	s.client = params.ClientInfo.Name + " " + params.ClientInfo.Version
	s.mu.Unlock()
	s.logger.Info("session initialized", zap.String("client", s.client))

	s.writeResult(&req.ID, map[string]any{
		"protocolVersion": ProtocolVersion,
		"capabilities": map[string]any{
			"tools": map[string]any{},
		},
		"serverInfo": map[string]any{
			"name":    "vedo-edutrack-mcp",
			"version": "0.1.0",
		},
	})
}

// handleToolsList returns the tool registry with JSON schemas.
func (s *Server) handleToolsList(req *rpcRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()
	tools := make([]map[string]any, 0, len(s.tools))
	for _, t := range s.tools {
		schema := t.InputSchema
		if schema == nil {
			schema = map[string]any{"type": "object", "properties": map[string]any{}}
		}
		tools = append(tools, map[string]any{
			"name":        t.Name,
			"description": t.Description,
			"inputSchema": schema,
		})
	}
	s.writeResult(&req.ID, map[string]any{"tools": tools})
}

// handleToolsCall invokes a tool by name.
func (s *Server) handleToolsCall(ctx context.Context, req *rpcRequest) {
	var params struct {
		Name string         `json:"name"`
		Args map[string]any `json:"arguments"`
	}
	_ = unmarshalParams(req, &params)

	tool, ok := s.tools[params.Name]
	if !ok {
		s.writeError(&req.ID, -32602, fmt.Sprintf("Unknown tool: %s", params.Name))
		return
	}
	if params.Args == nil {
		params.Args = map[string]any{}
	}

	s.logger.Debug("ToolCalled", zap.String("tool", params.Name), zap.Any("args", params.Args))
	result, err := tool.Handler(ctx, params.Args)
	if err != nil {
		s.logger.Error("ToolError", zap.String("tool", params.Name), zap.Error(err))
		s.writeError(&req.ID, -32603, err.Error())
		return
	}
	s.writeResult(&req.ID, map[string]any{
		"content": []map[string]any{
			{"type": "text", "text": jsonString(result)},
		},
	})
}

// triggerShutdown invokes the shutdown hook once.
func (s *Server) triggerShutdown() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	s.closed = true
	if s.shutdown != nil {
		s.shutdown()
	}
}

// rpcRequest is a JSON-RPC 2.0 request envelope.
type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

// writeResult writes a JSON-RPC success response.
func (s *Server) writeResult(id *json.RawMessage, result any) {
	resp := map[string]any{
		"jsonrpc": "2.0",
		"result":  result,
	}
	if id != nil {
		resp["id"] = id
	}
	s.write(resp)
}

// writeError writes a JSON-RPC error response.
func (s *Server) writeError(id *json.RawMessage, code int, message string) {
	resp := map[string]any{
		"jsonrpc": "2.0",
		"error": map[string]any{
			"code":    code,
			"message": message,
		},
	}
	if id != nil {
		resp["id"] = id
	}
	s.write(resp)
}

// write marshals and writes a response line to stdout.
func (s *Server) write(payload map[string]any) {
	data, err := json.Marshal(payload)
	if err != nil {
		s.logger.Error("marshal response", zap.Error(err))
		return
	}
	s.writeMu.Lock()
	_, _ = s.out.Write(append(data, '\n'))
	s.writeMu.Unlock()
}

// unmarshalParams decodes raw params into dst (no-op when params absent).
func unmarshalParams(req *rpcRequest, dst any) error {
	if len(req.Params) == 0 {
		return nil
	}
	return json.Unmarshal(req.Params, dst)
}

// jsonString renders a value as JSON text for the text content block.
func jsonString(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf(`{"error":%q}`, err.Error())
	}
	return string(data)
}
