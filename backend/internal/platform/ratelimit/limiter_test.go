package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestLimiterAllowsUpToCapacity(t *testing.T) {
	l := NewLimiter(0.1, 3, zap.NewNop()) // 3 burst tokens, refill 1 per 10s
	for i := 0; i < 3; i++ {
		allowed, retryAfter := l.Allow("user-1")
		if !allowed {
			t.Fatalf("request %d should be allowed (burst capacity 3)", i+1)
		}
		if retryAfter != 0 {
			t.Fatalf("expected no retry-after while allowed, got %d", retryAfter)
		}
	}
	allowed, retryAfter := l.Allow("user-1")
	if allowed {
		t.Fatal("expected denial after burst exhausted")
	}
	if retryAfter < 1 {
		t.Fatalf("expected positive retry-after, got %d", retryAfter)
	}
}

func TestLimiterRefillsOverTime(t *testing.T) {
	l := NewLimiter(1.0, 1, zap.NewNop()) // 1 token/s, burst 1
	if allowed, _ := l.Allow("user-1"); !allowed {
		t.Fatal("expected first request allowed")
	}
	if allowed, _ := l.Allow("user-1"); allowed {
		t.Fatal("expected immediate second request denied")
	}
	time.Sleep(1100 * time.Millisecond)
	if allowed, _ := l.Allow("user-1"); !allowed {
		t.Fatal("expected request allowed after refill")
	}
}

func TestLimiterIsolatesKeys(t *testing.T) {
	l := NewLimiter(0.1, 1, zap.NewNop())
	if allowed, _ := l.Allow("user-a"); !allowed {
		t.Fatal("expected user-a allowed")
	}
	if allowed, _ := l.Allow("user-b"); !allowed {
		t.Fatal("expected user-b allowed (separate bucket)")
	}
	if allowed, _ := l.Allow("user-a"); allowed {
		t.Fatal("expected user-a denied (bucket exhausted)")
	}
}

func TestLimiterDisabledWhenRateZero(t *testing.T) {
	l := NewLimiter(0, 1, zap.NewNop())
	for i := 0; i < 100; i++ {
		if allowed, _ := l.Allow("user-1"); !allowed {
			t.Fatal("expected pass-through when rate is 0")
		}
	}
}

func TestRateLimitMiddlewareReturns429WithRetryAfter(t *testing.T) {
	l := NewLimiter(0, 1, zap.NewNop()) // pass-through guard for setup
	_ = l

	// Enforce: rate so low every request after the first is denied.
	strict := NewLimiter(0.000001, 1, zap.NewNop())
	handler := strict.Middleware(func(r *http.Request) string { return "key-1" })(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First request allowed.
	req := httptest.NewRequest(http.MethodGet, "/sparql", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("first request: status=%d, want 200", rec.Code)
	}

	// Second request denied with 429 + Retry-After.
	req2 := httptest.NewRequest(http.MethodGet, "/sparql", nil)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("second request: status=%d, want 429", rec2.Code)
	}
	if rec2.Header().Get("Retry-After") == "" {
		t.Fatal("expected Retry-After header on 429")
	}
}

func TestRemoteAddrFallbackKey(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.5:54321"
	if got := remoteAddr(req); got != "10.0.0.5" {
		t.Fatalf("remoteAddr = %q, want 10.0.0.5", got)
	}
}
