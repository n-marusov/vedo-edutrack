package mcp

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics for the MCP server (F6). Registered lazily via promauto.
var (
	McpToolCallsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mcp_tool_calls_total",
			Help: "Total MCP tool calls, by tool name.",
		},
		[]string{"tool"},
	)
)
