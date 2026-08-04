package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go.uber.org/zap"
)

// testServer runs an MCP server over in-memory buffers and returns the input
// writer, output reader, and the server.
type testServer struct {
	mu     sync.Mutex
	in     *io.PipeWriter
	out    *bytes.Buffer
	server *Server
}

func newTestServer(t *testing.T, apiKey string, tools []Tool) *testServer {
	t.Helper()
	inR, inW := io.Pipe()
	out := &bytes.Buffer{}
	server := New(Config{APIKey: apiKey, Tools: tools}, inR, out, zap.NewNop())

	// Serve in a goroutine; the test writes to the pipe and reads from out.
	go func() {
		_ = server.Serve(context.Background())
	}()
	t.Cleanup(func() { _ = inW.Close() })

	return &testServer{in: inW, out: out, server: server}
}

// request writes a JSON-RPC request line and returns the response with the
// matching id. When id == 0, any response (e.g. a parse error without an id)
// is accepted.
func (ts *testServer) request(t *testing.T, payload string, id float64) map[string]any {
	t.Helper()
	ts.mu.Lock()
	_, _ = ts.in.Write([]byte(payload + "\n"))
	ts.mu.Unlock()

	// Wait for the response.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		ts.mu.Lock()
		raw := ts.out.String()
		ts.mu.Unlock()
		if resp := findResponse(t, raw, id); resp != nil {
			return resp
		}
		time.Sleep(5 * time.Millisecond)
	}
	ts.mu.Lock()
	raw := ts.out.String()
	ts.mu.Unlock()
	t.Fatalf("no response for id %v within timeout\n%s", id, raw)
	return nil
}

// findResponse scans raw output lines for a response with the matching id
// (any error response when id == 0).
func findResponse(t *testing.T, raw string, id float64) map[string]any {
	t.Helper()
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var resp map[string]any
		if err := json.Unmarshal([]byte(line), &resp); err != nil {
			t.Fatalf("unmarshal response %q: %v", line, err)
		}
		rid, ok := resp["id"].(float64)
		if id == 0 {
			// Parse errors carry no id; accept the first error response.
			if _, hasErr := resp["error"]; hasErr {
				return resp
			}
			continue
		}
		if ok && rid == id {
			return resp
		}
	}
	return nil
}

func echoTools() []Tool {
	return []Tool{
		{
			Name:        "echo",
			Description: "Echoes the message argument.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"message": map[string]any{"type": "string"},
				},
				"required": []string{"message"},
			},
			Handler: func(_ context.Context, args map[string]any) (any, error) {
				return map[string]any{"echo": args["message"]}, nil
			},
		},
	}
}

func TestInitializeReturnsCapabilities(t *testing.T) {
	ts := newTestServer(t, "", echoTools())
	r := ts.request(t, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"clientInfo":{"name":"test","version":"1.0"}}}`, 1)
	if r["id"] != float64(1) {
		t.Fatalf("id = %v, want 1", r["id"])
	}
	if _, hasErr := r["error"]; hasErr {
		t.Fatalf("unexpected error: %v", r["error"])
	}
	result, ok := r["result"].(map[string]any)
	if !ok {
		t.Fatalf("result not an object: %v", r["result"])
	}
	if result["protocolVersion"] != ProtocolVersion {
		t.Fatalf("protocolVersion = %v, want %s", result["protocolVersion"], ProtocolVersion)
	}
	serverInfo, _ := result["serverInfo"].(map[string]any)
	if serverInfo["name"] != "vedo-edutrack-mcp" {
		t.Fatalf("serverInfo.name = %v", serverInfo["name"])
	}
}

func TestToolsListReturnsDefinitions(t *testing.T) {
	ts := newTestServer(t, "", echoTools())
	r := ts.request(t, `{"jsonrpc":"2.0","id":2,"method":"initialize","params":{"clientInfo":{"name":"test","version":"1.0"}}}`, 2)
	_ = r // initialize
	r = ts.request(t, `{"jsonrpc":"2.0","id":3,"method":"tools/list","params":{}}`, 3)
	result, _ := r["result"].(map[string]any)
	tools, _ := result["tools"].([]any)
	if len(tools) != 1 {
		t.Fatalf("tools len = %d, want 1", len(tools))
	}
	first := tools[0].(map[string]any)
	if first["name"] != "echo" || first["description"] == "" {
		t.Fatalf("unexpected tool: %v", first)
	}
	if first["inputSchema"] == nil {
		t.Fatal("expected inputSchema")
	}
}

func TestToolsCallInvokesHandler(t *testing.T) {
	ts := newTestServer(t, "", echoTools())
	r := ts.request(t, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"clientInfo":{"name":"test","version":"1.0"}}}`, 1)
	_ = r
	r = ts.request(t, `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"echo","arguments":{"message":"hello"}}}`, 4)
	result, _ := r["result"].(map[string]any)
	content, _ := result["content"].([]any)
	if len(content) != 1 {
		t.Fatalf("content len = %d, want 1", len(content))
	}
	textBlock := content[0].(map[string]any)
	if textBlock["type"] != "text" {
		t.Fatalf("content type = %v, want text", textBlock["type"])
	}
	if !strings.Contains(textBlock["text"].(string), "hello") {
		t.Fatalf("text = %v, want echo of hello", textBlock["text"])
	}
}

func TestMissingAuthRejected(t *testing.T) {
	ts := newTestServer(t, "secret-key", echoTools())
	// tools/list before initialize → auth error.
	r := ts.request(t, `{"jsonrpc":"2.0","id":9,"method":"tools/list","params":{}}`, 9)
	errObj, ok := r["error"].(map[string]any)
	if !ok {
		t.Fatalf("expected error response, got %v", r)
	}
	if errObj["code"] != float64(-32002) {
		t.Fatalf("error code = %v, want -32002", errObj["code"])
	}
}

func TestInitializeRejectsWrongAPIKey(t *testing.T) {
	ts := newTestServer(t, "secret-key", echoTools())
	r := ts.request(t, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"clientInfo":{"name":"test","version":"1.0"},"capabilities":{"apiKey":"wrong"}}}`, 1)
	if _, hasErr := r["error"]; !hasErr {
		t.Fatalf("expected error for wrong API key, got %v", r)
	}
}

func TestUnknownToolRejected(t *testing.T) {
	ts := newTestServer(t, "", echoTools())
	r := ts.request(t, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"clientInfo":{"name":"test","version":"1.0"}}}`, 1)
	_ = r
	r = ts.request(t, `{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"nope","arguments":{}}}`, 5)
	errObj, ok := r["error"].(map[string]any)
	if !ok {
		t.Fatalf("expected error, got %v", r)
	}
	if errObj["code"] != float64(-32602) {
		t.Fatalf("error code = %v, want -32602", errObj["code"])
	}
}

func TestMalformedJSONReturnsParseError(t *testing.T) {
	ts := newTestServer(t, "", echoTools())
	r := ts.request(t, `{not valid json`, 0)
	errObj, ok := r["error"].(map[string]any)
	if !ok {
		t.Fatalf("expected parse error, got %v", r)
	}
	if errObj["code"] != float64(-32700) {
		t.Fatalf("error code = %v, want -32700", errObj["code"])
	}
}

func TestUnknownMethodRejected(t *testing.T) {
	ts := newTestServer(t, "", echoTools())
	r := ts.request(t, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"clientInfo":{"name":"test","version":"1.0"}}}`, 1)
	_ = r
	r = ts.request(t, `{"jsonrpc":"2.0","id":6,"method":"bogus","params":{}}`, 6)
	errObj, _ := r["error"].(map[string]any)
	if errObj["code"] != float64(-32601) {
		t.Fatalf("error code = %v, want -32601", errObj["code"])
	}
}

func TestStdinCloseTriggersShutdown(t *testing.T) {
	var shutdown atomic.Int32
	inR, inW := io.Pipe()
	out := &bytes.Buffer{}
	server := New(Config{
		Tools: echoTools(),
		OnShutdown: func() {
			shutdown.Add(1)
		},
	}, inR, out, zap.NewNop())

	done := make(chan struct{})
	go func() {
		_ = server.Serve(context.Background())
		close(done)
	}()

	// Write a request, then close stdin → EOF → Serve returns and shutdown fires.
	_, _ = inW.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"clientInfo":{"name":"test","version":"1.0"}}}` + "\n"))
	_ = inW.Close()

	select {
	case <-done:
		// Serve returned after stdin EOF.
	case <-time.After(3 * time.Second):
		t.Fatal("Serve did not return on stdin EOF")
	}
	if shutdown.Load() != 1 {
		t.Fatalf("shutdown calls = %d, want 1", shutdown.Load())
	}
}
