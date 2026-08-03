package integrations

import (
	"strings"
	"testing"
	"time"
)

func TestEventTypeEnum(t *testing.T) {
	known := []EventType{EventModuleMastered, EventPlanDeviated, EventRouteRecalculated, EventStandardRiskDetected}
	for _, tt := range known {
		if !tt.Valid() {
			t.Fatalf("expected %q to be a valid event type", tt)
		}
		if tt.String() != string(tt) {
			t.Fatalf("String() = %q, want %q", tt.String(), string(tt))
		}
	}
	if EventType("made.up.event").Valid() {
		t.Fatal("expected unknown event type to be invalid")
	}
	if !strings.Contains(KnownEventTypes[0].String(), ".") {
		t.Fatal("event types must use the dotted wire format (module.mastered)")
	}
}

func TestDeliveryStatusEnum(t *testing.T) {
	known := []DeliveryStatus{DeliveryPending, DeliverySent, DeliveryFailed, DeliveryPermanentFail}
	for _, s := range known {
		if !s.Valid() {
			t.Fatalf("expected %q to be a valid delivery status", s)
		}
	}
	if DeliveryStatus("lost").Valid() {
		t.Fatal("expected unknown delivery status to be invalid")
	}
}

func TestValidateSubscriptionRequiresTenant(t *testing.T) {
	err := ValidateSubscription("", "https://example.test/hook", []EventType{EventModuleMastered}, strings.Repeat("s", 32))
	if err == nil {
		t.Fatal("expected tenant required error")
	}
	if !strings.Contains(err.Error(), "tenant") {
		t.Fatalf("expected tenant reason, got %v", err)
	}
}

func TestValidateSubscriptionURLPolicy(t *testing.T) {
	cases := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{name: "https allowed", url: "https://hooks.example.test/v1", wantErr: false},
		{name: "localhost http allowed", url: "http://localhost:9000/hook", wantErr: false},
		{name: "127.0.0.1 http allowed", url: "http://127.0.0.1:9000/hook", wantErr: false},
		{name: "plain http rejected", url: "http://insecure.example.test/hook", wantErr: true},
		{name: "missing scheme rejected", url: "hooks.example.test/v1", wantErr: true},
		{name: "missing host rejected", url: "https:///path", wantErr: true},
		{name: "ftp rejected", url: "ftp://example.test/hook", wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateSubscription("tenant-1", tc.url, []EventType{EventModuleMastered}, strings.Repeat("s", 32))
			if (err != nil) != tc.wantErr {
				t.Fatalf("url=%q err=%v, wantErr=%t", tc.url, err, tc.wantErr)
			}
		})
	}
}

func TestValidateSubscriptionEventTypes(t *testing.T) {
	if err := ValidateSubscription("t1", "https://example.test/hook", nil, strings.Repeat("s", 32)); err == nil {
		t.Fatal("expected error for empty event types")
	}
	if err := ValidateSubscription("t1", "https://example.test/hook", []EventType{"bogus"}, strings.Repeat("s", 32)); err == nil {
		t.Fatal("expected error for unknown event type")
	}
	dupErr := ValidateSubscription("t1", "https://example.test/hook",
		[]EventType{EventModuleMastered, EventModuleMastered}, strings.Repeat("s", 32))
	if dupErr == nil || !strings.Contains(dupErr.Error(), "duplicate") {
		t.Fatalf("expected duplicate event type error, got %v", dupErr)
	}
	if err := ValidateSubscription("t1", "https://example.test/hook",
		[]EventType{EventModuleMastered, EventPlanDeviated}, strings.Repeat("s", 32)); err != nil {
		t.Fatalf("expected valid subscription, got %v", err)
	}
}

func TestValidateSubscriptionSigningSecretMinLength(t *testing.T) {
	short := strings.Repeat("s", MinSigningSecretLength-1)
	err := ValidateSubscription("t1", "https://example.test/hook", []EventType{EventModuleMastered}, short)
	if err == nil || !strings.Contains(err.Error(), "signing secret") {
		t.Fatalf("expected signing secret length error, got %v", err)
	}
	long := strings.Repeat("s", MinSigningSecretLength)
	if err := ValidateSubscription("t1", "https://example.test/hook", []EventType{EventModuleMastered}, long); err != nil {
		t.Fatalf("expected valid signing secret, got %v", err)
	}
	// Empty secret is allowed at creation (a random one is generated server-side).
	if err := ValidateSubscription("t1", "https://example.test/hook", []EventType{EventModuleMastered}, ""); err != nil {
		t.Fatalf("expected empty signing secret allowed, got %v", err)
	}
}

func TestNextDeliveryAttempt(t *testing.T) {
	if got := NextDeliveryAttempt(nil); got != 1 {
		t.Fatalf("nil delivery: attempt=%d, want 1", got)
	}
	if got := NextDeliveryAttempt(&WebhookDelivery{Attempt: 0}); got != 1 {
		t.Fatalf("fresh delivery: attempt=%d, want 1", got)
	}
	if got := NextDeliveryAttempt(&WebhookDelivery{Attempt: 4}); got != 5 {
		t.Fatalf("retry delivery: attempt=%d, want 5", got)
	}
}

func TestShouldDeactivateOnConsecutiveFailures(t *testing.T) {
	if ShouldDeactivate(&WebhookDelivery{Status: DeliveryFailed, Attempt: 3}) {
		t.Fatal("expected no deactivation before the failure boundary")
	}
	if !ShouldDeactivate(&WebhookDelivery{Status: DeliveryFailed, Attempt: MaxConsecutiveDeliveryFailures}) {
		t.Fatalf("expected deactivation at %d consecutive failures", MaxConsecutiveDeliveryFailures)
	}
	if !ShouldDeactivate(&WebhookDelivery{Status: DeliveryPermanentFail, Attempt: 1}) {
		t.Fatal("expected deactivation on permanent failure")
	}
	if ShouldDeactivate(nil) {
		t.Fatal("expected no deactivation for nil delivery")
	}
}

func TestWebhookDeliveryLifecycleFields(t *testing.T) {
	now := time.Now().UTC()
	d := WebhookDelivery{
		ID:             "d1",
		SubscriptionID: SubscriptionID("sub-1"),
		EventID:        "evt-1",
		Attempt:        1,
		Status:         DeliveryPending,
		LastAttemptAt:  &now,
		HTTPStatus:     200,
		ResponseBody:   "{}",
	}
	if !d.Status.Valid() {
		t.Fatalf("expected pending to be valid, got %q", d.Status)
	}
	if d.HTTPStatus != 200 {
		t.Fatalf("unexpected http status: %d", d.HTTPStatus)
	}
}
