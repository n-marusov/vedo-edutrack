// Package hub implements the VEDO Hub GraphQL adapter for ontology reads.
package hub

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	ontology "vedo-edutrack/backend/internal/modules/ontologyport/domain"
	platformconfig "vedo-edutrack/backend/internal/platform/config"
)

const (
	defaultGraphQLPath = "/graphql"
	defaultTimeout     = 3 * time.Second
)

// Config configures the VEDO Hub GraphQL client.
type Config struct {
	BaseURL     string
	GraphQLPath string
	BearerToken string
	Timeout     time.Duration
}

// Client is a read-only GraphQL client for the VEDO Hub ontology service.
type Client struct {
	endpoint string
	token    string
	timeout  time.Duration
	http     *http.Client
	logger   *zap.Logger
}

// NewClientFromConfig creates a Hub GraphQL client from the shared platform configuration.
func NewClientFromConfig(cfg *platformconfig.Config, logger *zap.Logger) (*Client, error) {
	if cfg == nil {
		return nil, errors.New("platform config is required")
	}
	return NewClient(Config{
		BaseURL:     cfg.HubAPIURL,
		GraphQLPath: cfg.HubGraphQLPath,
		BearerToken: cfg.HubBearerToken,
		Timeout:     defaultTimeout,
	}, logger)
}

// NewClient creates a Hub GraphQL client with the F0 timeout boundary enforced.
func NewClient(cfg Config, logger *zap.Logger) (*Client, error) {
	baseURL := strings.TrimSpace(cfg.BaseURL)
	if baseURL == "" {
		return nil, errors.New("hub base URL is required")
	}

	path := cfg.GraphQLPath
	if path == "" {
		path = defaultGraphQLPath
	}
	endpoint := strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/")

	timeout := cfg.Timeout
	if timeout <= 0 || timeout > defaultTimeout {
		timeout = defaultTimeout
	}
	if logger == nil {
		logger = zap.NewNop()
	}

	client := &Client{
		endpoint: endpoint,
		token:    cfg.BearerToken,
		timeout:  timeout,
		http:     &http.Client{Timeout: timeout},
		logger:   logger.Named("ontologyport.hub"),
	}
	client.logger.Info("connected to Hub", zap.String("url", baseURL), zap.String("endpoint", endpoint), zap.Duration("timeout", timeout))
	return client, nil
}

type graphQLRequest struct {
	OperationName string         `json:"operationName,omitempty"`
	Query         string         `json:"query"`
	Variables     map[string]any `json:"variables,omitempty"`
}

type graphQLError struct {
	Message string         `json:"message"`
	Path    []any          `json:"path,omitempty"`
	Ext     map[string]any `json:"extensions,omitempty"`
}

func (c *Client) execute(ctx context.Context, operation, query string, variables map[string]any, out any) error {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	started := time.Now()
	payload, err := json.Marshal(graphQLRequest{OperationName: operation, Query: query, Variables: variables})
	if err != nil {
		return fmt.Errorf("marshal GraphQL request %s: %w", operation, err)
	}

	c.logger.Debug("GraphQL query", zap.String("operation", operation), zap.Any("variables", variables))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("build GraphQL request %s: %w", operation, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		c.logger.Error("query failed", zap.String("operation", operation), zap.Error(err))
		return fmt.Errorf("execute GraphQL operation %s: %w", operation, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read GraphQL response %s: %w", operation, err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		c.logger.Error("query failed", zap.String("operation", operation), zap.Int("status", resp.StatusCode), zap.ByteString("body", body))
		return fmt.Errorf("GraphQL operation %s failed with status %d", operation, resp.StatusCode)
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decode GraphQL response %s: %w", operation, err)
	}
	c.logger.Debug("GraphQL response", zap.String("operation", operation), zap.Duration("duration", time.Since(started)))
	return nil
}

func graphQLErrors(operation string, errs []graphQLError) error {
	if len(errs) == 0 {
		return nil
	}
	messages := make([]string, 0, len(errs))
	for _, e := range errs {
		messages = append(messages, e.Message)
	}
	return fmt.Errorf("GraphQL operation %s returned errors: %s", operation, strings.Join(messages, "; "))
}

func metadataFromAny(raw map[string]any) map[string]string {
	if len(raw) == 0 {
		return nil
	}
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		switch value := v.(type) {
		case string:
			out[k] = value
		case nil:
			continue
		default:
			out[k] = fmt.Sprint(value)
		}
	}
	return out
}

func validateLink(link ontology.OntologyLink) error {
	if link.SourceModuleID == "" || link.TargetModuleID == "" {
		return fmt.Errorf("ontology link requires source and target ids: %+v", link)
	}
	if err := link.Type.Validate(); err != nil {
		return err
	}
	return nil
}
