package config

import (
	"testing"
)

func TestLoadHubGraphQLSettings(t *testing.T) {
	t.Setenv("VEDO_HUB_API_URL", "http://hub.example:8081")
	t.Setenv("VEDO_HUB_GRAPHQL_PATH", "/custom/graphql")
	t.Setenv("VEDO_HUB_BEARER_TOKEN", "secret-token")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.HubAPIURL != "http://hub.example:8081" {
		t.Errorf("HubAPIURL = %q, want http://hub.example:8081", cfg.HubAPIURL)
	}
	if cfg.HubGraphQLPath != "/custom/graphql" {
		t.Errorf("HubGraphQLPath = %q, want /custom/graphql", cfg.HubGraphQLPath)
	}
	if cfg.HubBearerToken != "secret-token" {
		t.Errorf("HubBearerToken = %q, want secret-token", cfg.HubBearerToken)
	}
}

func TestLoadHubGraphQLDefaults(t *testing.T) {
	// Ensure no VEDO_HUB_* overrides leak from the environment.
	t.Setenv("VEDO_HUB_API_URL", "")
	t.Setenv("VEDO_HUB_GRAPHQL_PATH", "")
	t.Setenv("VEDO_HUB_BEARER_TOKEN", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.HubAPIURL != "http://localhost:8081" {
		t.Errorf("HubAPIURL = %q, want default http://localhost:8081", cfg.HubAPIURL)
	}
	if cfg.HubGraphQLPath != "/graphql" {
		t.Errorf("HubGraphQLPath = %q, want default /graphql", cfg.HubGraphQLPath)
	}
	if cfg.HubBearerToken != "" {
		t.Errorf("HubBearerToken = %q, want empty default", cfg.HubBearerToken)
	}
}
