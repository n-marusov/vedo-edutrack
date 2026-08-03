// Package integrations provides the domain model for the integrations bounded
// context (F6): webhook subscriptions/deliveries, event deduplication, and the
// SPARQL query result contract.
//
// The domain is pure — no I/O, no logging, no framework imports. State
// transitions (delivery retry → permanent failure → subscription deactivation)
// are modeled here so the business rules are testable without infrastructure
// (ADR-IMPL.PROCESS.repository-structure §5).
package integrations

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// EventType identifies a webhook event contract. Event names are the stable
// public wire format consumed by subscribers (REQ-FR-api.webhooks.*).
type EventType string

// Known webhook event types.
const (
	EventModuleMastered       EventType = "module.mastered"
	EventPlanDeviated         EventType = "plan.deviated"
	EventRouteRecalculated    EventType = "route.recalculated"
	EventStandardRiskDetected EventType = "standard.risk_detected"
)

// KnownEventTypes lists every valid webhook event type (validation source).
var KnownEventTypes = []EventType{
	EventModuleMastered,
	EventPlanDeviated,
	EventRouteRecalculated,
	EventStandardRiskDetected,
}

// Valid reports whether the event type is part of the known set.
func (t EventType) Valid() bool {
	for _, known := range KnownEventTypes {
		if t == known {
			return true
		}
	}
	return false
}

// String returns the wire representation of the event type.
func (t EventType) String() string { return string(t) }

// DeliveryStatus is the state of a single webhook delivery attempt.
type DeliveryStatus string

// Delivery lifecycle states (ADR-§3 webhook delivery state machine).
const (
	DeliveryPending       DeliveryStatus = "pending"
	DeliverySent          DeliveryStatus = "sent"
	DeliveryFailed        DeliveryStatus = "failed"
	DeliveryPermanentFail DeliveryStatus = "permanent_failure"
)

// Valid reports whether the delivery status is part of the known set.
func (s DeliveryStatus) Valid() bool {
	switch s {
	case DeliveryPending, DeliverySent, DeliveryFailed, DeliveryPermanentFail:
		return true
	default:
		return false
	}
}

// MaxConsecutiveDeliveryFailures is the business rule boundary: a subscription
// is deactivated after this many consecutive failed deliveries.
const MaxConsecutiveDeliveryFailures = 5

// MaxSubscriptionsPerTenant caps the number of active webhook subscriptions a
// tenant may own (REQ-FR-api.webhooks.business-rules).
const MaxSubscriptionsPerTenant = 10

// MinSigningSecretLength is the minimum accepted webhook signing secret length.
const MinSigningSecretLength = 32

// SubscriptionID is a typed webhook subscription identifier.
type SubscriptionID string

// String returns the string form of the subscription id.
func (id SubscriptionID) String() string { return string(id) }

// WebhookSubscription is a subscriber endpoint bound to a set of event types.
type WebhookSubscription struct {
	ID            SubscriptionID
	TenantID      string
	URL           string
	EventTypes    []EventType
	SigningSecret string
	Active        bool
	Failures      int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// WebhookDelivery records one delivery attempt for an outbox event to a
// subscription.
type WebhookDelivery struct {
	ID             string
	SubscriptionID SubscriptionID
	EventID        string
	Attempt        int
	Status         DeliveryStatus
	LastAttemptAt  *time.Time
	HTTPStatus     int
	ResponseBody   string
	Error          string
}

// EventDedup is the idempotency record for webhook receipt: a (source,
// event_id) pair is delivered at most once (REQ-FR-api.webhooks.idempotency,
// ADR-§3).
type EventDedup struct {
	Source  string
	EventID string
}

// SPARQLQueryResult is the read-only SPARQL result contract
// (application/sparql-results+json shape). Truncated is set when the result
// set exceeded the row boundary and was clipped.
type SPARQLQueryResult struct {
	Head      SPARQLResultHead
	Results   SPARQLResultBody
	Truncated bool
}

// SPARQLResultHead carries the result variable names.
type SPARQLResultHead struct {
	Vars []string
}

// SPARQLResultBody carries the result bindings (row × variable → value).
type SPARQLResultBody struct {
	Bindings []map[string]SPARQLBinding
}

// SPARQLBinding is one variable binding (type + value per the SPARQL JSON
// results format).
type SPARQLBinding struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// SubscriptionError is a typed validation error for webhook subscriptions.
type SubscriptionError struct {
	Reason string
}

// Error implements the error interface.
func (e *SubscriptionError) Error() string {
	return fmt.Sprintf("subscription validation failed: %s", e.Reason)
}

// ValidateSubscription validates a webhook subscription create/update payload.
// Rules (REQ-FR-api.webhooks.*):
//   - URL must be absolute https (http allowed for localhost/dev sandbox)
//   - event types must come from the known set
//   - signing secret must be at least MinSigningSecretLength chars when set
func ValidateSubscription(tenantID, rawURL string, eventTypes []EventType, signingSecret string) error {
	if strings.TrimSpace(tenantID) == "" {
		return &SubscriptionError{Reason: "tenant is required"}
	}
	if err := validateWebhookURL(rawURL); err != nil {
		return &SubscriptionError{Reason: err.Error()}
	}
	if len(eventTypes) == 0 {
		return &SubscriptionError{Reason: "at least one event type is required"}
	}
	seen := map[EventType]bool{}
	for _, t := range eventTypes {
		if !t.Valid() {
			return &SubscriptionError{Reason: fmt.Sprintf("unknown event type %q", t)}
		}
		if seen[t] {
			return &SubscriptionError{Reason: fmt.Sprintf("duplicate event type %q", t)}
		}
		seen[t] = true
	}
	if signingSecret != "" && len(signingSecret) < MinSigningSecretLength {
		return &SubscriptionError{Reason: fmt.Sprintf("signing secret must be at least %d chars", MinSigningSecretLength)}
	}
	return nil
}

// validateWebhookURL enforces the subscriber URL policy: https by default,
// http allowed only for localhost (dev sandbox).
func validateWebhookURL(rawURL string) error {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		return fmt.Errorf("URL scheme must be https (or http for localhost)")
	}
	if u.Host == "" {
		return fmt.Errorf("URL must have a host")
	}
	if u.Scheme == "http" && !isLocalhost(u.Hostname()) {
		return fmt.Errorf("http URLs are only allowed for localhost (use https in production)")
	}
	return nil
}

// isLocalhost reports whether a host is a loopback host.
func isLocalhost(host string) bool {
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}

// NextDeliveryAttempt computes the next attempt number for a delivery.
// Returns 1 for a fresh delivery; otherwise attempt+1.
func NextDeliveryAttempt(delivery *WebhookDelivery) int {
	if delivery == nil || delivery.Attempt <= 0 {
		return 1
	}
	return delivery.Attempt + 1
}

// ShouldDeactivate reports whether a subscription crossed the consecutive
// failure boundary after the given delivery result.
func ShouldDeactivate(delivery *WebhookDelivery) bool {
	if delivery == nil {
		return false
	}
	return delivery.Status == DeliveryPermanentFail || delivery.Attempt >= MaxConsecutiveDeliveryFailures
}
