package circuitbreaker

import (
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestBreakerClosedAllowsCalls(t *testing.T) {
	b := New(DefaultConfig(), zap.NewNop())
	if b.State() != StateClosed {
		t.Fatalf("initial state = %s, want closed", b.State())
	}
	allowed, err := b.Allow()
	if !allowed || err != nil {
		t.Fatalf("expected closed circuit to allow, got allowed=%t err=%v", allowed, err)
	}
}

func TestBreakerOpensAfterFailureThreshold(t *testing.T) {
	b := New(Config{FailureThreshold: 3, Timeout: time.Minute, SuccessThreshold: 2}, zap.NewNop())
	for i := 0; i < 3; i++ {
		b.Failure()
	}
	if b.State() != StateOpen {
		t.Fatalf("state = %s, want open", b.State())
	}
	allowed, err := b.Allow()
	if allowed || !errors.Is(err, ErrOpen) {
		t.Fatalf("expected open circuit to fail fast, got allowed=%t err=%v", allowed, err)
	}
}

func TestBreakerRecoversToHalfOpenThenClosed(t *testing.T) {
	b := New(Config{FailureThreshold: 2, Timeout: 50 * time.Millisecond, SuccessThreshold: 2}, zap.NewNop())
	b.Failure()
	b.Failure()
	if b.State() != StateOpen {
		t.Fatalf("state = %s, want open", b.State())
	}

	// After the timeout the circuit transitions to half-open on next check.
	time.Sleep(60 * time.Millisecond)
	allowed, err := b.Allow()
	if !allowed || err != nil {
		t.Fatalf("expected half-open probe allowed, got allowed=%t err=%v", allowed, err)
	}
	if b.State() != StateHalfOpen {
		t.Fatalf("state = %s, want half-open", b.State())
	}

	// Consecutive successes close the circuit.
	b.Success()
	b.Success()
	if b.State() != StateClosed {
		t.Fatalf("state = %s, want closed after success threshold", b.State())
	}
}

func TestBreakerHalfOpenFailureReopens(t *testing.T) {
	b := New(Config{FailureThreshold: 2, Timeout: 20 * time.Millisecond, SuccessThreshold: 2}, zap.NewNop())
	b.Failure()
	b.Failure()
	time.Sleep(30 * time.Millisecond)
	if _, err := b.Allow(); err != nil {
		t.Fatalf("expected half-open probe, got %v", err)
	}
	b.Failure() // one failure in half-open re-opens
	if b.State() != StateOpen {
		t.Fatalf("state = %s, want open after half-open failure", b.State())
	}
}

func TestBreakerExecuteGuardsAndRecords(t *testing.T) {
	b := New(Config{FailureThreshold: 2, Timeout: time.Minute, SuccessThreshold: 1}, zap.NewNop())
	b.Failure()
	b.Failure()

	// Circuit open: Execute fails fast without invoking fn.
	invoked := false
	err := b.Execute(func() error {
		invoked = true
		return nil
	})
	if !errors.Is(err, ErrOpen) {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
	if invoked {
		t.Fatal("expected fn not invoked while circuit open")
	}
}

func TestBreakerExecuteRecordsSuccess(t *testing.T) {
	b := New(DefaultConfig(), zap.NewNop())
	err := b.Execute(func() error { return nil })
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if b.State() != StateClosed {
		t.Fatalf("state = %s, want closed", b.State())
	}
}

func TestBreakerExecuteRecordsFailure(t *testing.T) {
	b := New(Config{FailureThreshold: 1, Timeout: time.Minute, SuccessThreshold: 1}, zap.NewNop())
	sentinel := errors.New("hub down")
	err := b.Execute(func() error { return sentinel })
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if b.State() != StateOpen {
		t.Fatalf("state = %s, want open after failure", b.State())
	}
}

func TestBreakerReset(t *testing.T) {
	b := New(Config{FailureThreshold: 1, Timeout: time.Minute, SuccessThreshold: 1}, zap.NewNop())
	b.Failure()
	if b.State() != StateOpen {
		t.Fatalf("state = %s, want open", b.State())
	}
	b.Reset()
	if b.State() != StateClosed {
		t.Fatalf("state = %s, want closed after reset", b.State())
	}
}
