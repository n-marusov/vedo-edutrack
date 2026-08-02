// Package logger provides a shared structured logger (zap + OTel bridge).
//
// All bounded contexts use this logger through the platform layer,
// ensuring consistent log format (JSON), correlation IDs (trace_id, span_id),
// and configurable verbosity via LOG_LEVEL.
//
// See ADR-IMPL.PROCESS.development-tooling §4 and
// REQ-NFR-ops.observability.log-level-config.
package logger

import (
	"log"
)

// New creates a new structured logger.
//
// TODO: Wire go.uber.org/zap with otelzap bridge.
// Logs are JSON-formatted and routed to stdout (collected by OTel Collector
// in the docker-compose stack). Trace correlation is automatic via OTel bridge.
func New(level string) error {
	log.Printf("[INFO] [logger.New] creating logger with level %q", level)
	_ = level
	return nil
}
