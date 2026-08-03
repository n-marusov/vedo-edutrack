package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// checkHealth implements the `health` subcommand — a self-contained liveness
// probe used by the container HEALTHCHECK (distroless images have no curl/wget):
//
//	HEALTHCHECK CMD ["/usr/local/bin/vedo-edutrack", "health"]
//
// It GETs http://127.0.0.1:<PORT>/healthz and exits 0 on HTTP 200, 1 otherwise.
func checkHealth() error {
	port := 8080
	if raw := os.Getenv("PORT"); raw != "" {
		if p, err := strconv.Atoi(raw); err == nil {
			port = p
		}
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://127.0.0.1:" + strconv.Itoa(port) + "/healthz")
	if err != nil {
		return fmt.Errorf("health probe failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, resp.Body)
		return fmt.Errorf("health probe: unexpected status %d", resp.StatusCode)
	}

	_, _ = io.Copy(io.Discard, resp.Body)
	return nil
}
