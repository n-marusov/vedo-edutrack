// Package sparql implements the read-only SPARQL endpoint guard and the
// VEDO Hub SPARQL proxy (F6).
//
// The EduTrack SPARQL surface is read-only by contract
// (REQ-FR-api.sparql.read-only): only SELECT / ASK / DESCRIBE / CONSTRUCT are
// permitted. Mutating queries (INSERT, DELETE, LOAD, CLEAR, CREATE, DROP) are
// rejected before they reach the underlying store.
package sparql

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

// DefaultHubTimeout is the execution timeout for the Hub SPARQL proxy (30s
// per the M4 plan; the caller maps a timeout to 504 Gateway Timeout).
const DefaultHubTimeout = 30 * time.Second

// ErrTimeout is returned when the Hub SPARQL call exceeded the timeout
// (caller maps it to 504 Gateway Timeout).
var ErrTimeout = errors.New("sparql query timed out")

// isTimeout reports whether err is a net/url timeout or context deadline.
func isTimeout(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}
	return errors.Is(err, context.DeadlineExceeded)
}

// Client proxies read-only SPARQL queries to the VEDO Hub SPARQL endpoint.
type Client struct {
	endpoint string
	token    string
	timeout  time.Duration
	http     *http.Client
	logger   *zap.Logger
}

// ClientConfig configures the Hub SPARQL proxy client.
type ClientConfig struct {
	// BaseURL is VEDO_HUB_API_URL (the Hub base URL).
	BaseURL string
	// Path is the Hub SPARQL endpoint path (VEDO_HUB_SPARQL_PATH, default /sparql).
	Path string
	// BearerToken is the optional Hub read-only bearer token.
	BearerToken string
	// Timeout is the execution timeout; ≤0 means DefaultHubTimeout.
	Timeout time.Duration
}

// NewClient builds the Hub SPARQL proxy client.
func NewClient(cfg ClientConfig, logger *zap.Logger) (*Client, error) {
	baseURL := strings.TrimSpace(cfg.BaseURL)
	if baseURL == "" {
		return nil, errors.New("hub base URL is required")
	}
	path := cfg.Path
	if path == "" {
		path = "/sparql"
	}
	endpoint := strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/")

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = DefaultHubTimeout
	}
	if logger == nil {
		logger = zap.NewNop()
	}

	c := &Client{
		endpoint: endpoint,
		token:    cfg.BearerToken,
		timeout:  timeout,
		http:     &http.Client{Timeout: timeout},
		logger:   logger.Named("sparql.hub"),
	}
	c.logger.Info("sparql proxy client ready", zap.String("endpoint", endpoint), zap.Duration("timeout", timeout))
	return c, nil
}

// ResultSet is the decoded SPARQL JSON results document.
type ResultSet struct {
	Head struct {
		Vars []string `json:"vars"`
	} `json:"head"`
	Results struct {
		Bindings []map[string]any `json:"bindings"`
	} `json:"results"`
}

// Query executes a read-only SPARQL query against VEDO Hub and returns the
// decoded result set. A timeout surfaces as an error wrapping ErrTimeout (the
// caller maps it to 504); transport errors wrap err.
func (c *Client) Query(ctx context.Context, query string) (*ResultSet, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	form := url.Values{}
	form.Set("query", query)

	started := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint+"?"+form.Encode(), bytes.NewReader(nil))
	if err != nil {
		return nil, fmt.Errorf("build SPARQL request: %w", err)
	}
	req.Header.Set("Accept", "application/sparql-results+json, application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	c.logger.Debug("query", zap.String("queryLength", fmt.Sprint(len(query))))
	resp, err := c.http.Do(req)
	if err != nil {
		c.logger.Error("query failed", zap.Error(err))
		if isTimeout(err) || ctx.Err() != nil {
			return nil, fmt.Errorf("%w: %v", ErrTimeout, err)
		}
		return nil, fmt.Errorf("execute SPARQL query: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read SPARQL response: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		c.logger.Error("query failed", zap.Int("status", resp.StatusCode), zap.ByteString("body", body))
		return nil, fmt.Errorf("SPARQL endpoint returned status %d", resp.StatusCode)
	}

	var result ResultSet
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode SPARQL results: %w", err)
	}
	c.logger.Debug("response", zap.Duration("duration", time.Since(started)), zap.Int("rows", len(result.Results.Bindings)))
	return &result, nil
}
