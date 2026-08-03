// Package circuitbreaker implements a circuit breaker for outbound VEDO Hub
// calls (F6, ADR-DES.API.communication-patterns §5).
//
// States: closed → open → half-open → closed.
//   - closed:   calls pass through; N consecutive failures open the circuit
//   - open:     calls fail fast (ErrOpen) for the timeout window
//   - half-open: a probe call is allowed; success closes, failure re-opens
//
// The breaker is per-service (e.g. the SPARQL proxy) with configurable
// thresholds. Timeout boundary: ≤3s per ADR-§5 (enforced by the caller's
// context; the breaker itself does not perform I/O).
package circuitbreaker

import (
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"
)

// ErrOpen is returned when the circuit is open (fail fast, no attempt).
var ErrOpen = errors.New("circuit breaker open: service unavailable")

// State is the breaker state.
type State string

// Breaker states.
const (
	StateClosed   State = "closed"
	StateOpen     State = "open"
	StateHalfOpen State = "half-open"
)

// Config tunes the breaker thresholds.
type Config struct {
	// FailureThreshold is the consecutive-failure count that opens the circuit.
	FailureThreshold int
	// Timeout is how long the circuit stays open before a half-open probe.
	Timeout time.Duration
	// SuccessThreshold is the consecutive-success count (half-open) that
	// closes the circuit.
	SuccessThreshold int
}

// DefaultConfig returns the ADR-§5 defaults: 5 failures, 30s timeout, 3
// successes to close.
func DefaultConfig() Config {
	return Config{
		FailureThreshold: 5,
		Timeout:          30 * time.Second,
		SuccessThreshold: 3,
	}
}

// Breaker is a concurrency-safe circuit breaker.
type Breaker struct {
	mu        sync.Mutex
	config    Config
	state     State
	failures  int
	successes int
	openedAt  time.Time
	logger    *zap.Logger
}

// New builds a breaker with the given config (defaults applied for zero fields).
func New(cfg Config, logger *zap.Logger) *Breaker {
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = DefaultConfig().FailureThreshold
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = DefaultConfig().Timeout
	}
	if cfg.SuccessThreshold <= 0 {
		cfg.SuccessThreshold = DefaultConfig().SuccessThreshold
	}
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Breaker{config: cfg, state: StateClosed, logger: logger.Named("circuitbreaker")}
}

// State returns the current state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.effectiveState(time.Now())
}

// Allow reports whether a call may proceed. When the circuit is open, it
// fails fast with ErrOpen.
func (b *Breaker) Allow() (bool, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	state := b.effectiveState(time.Now())
	if state == StateOpen {
		return false, ErrOpen
	}
	return true, nil
}

// effectiveState computes the state considering the open timeout: an open
// circuit transitions to half-open after the timeout elapses.
func (b *Breaker) effectiveState(now time.Time) State {
	if b.state == StateOpen && now.Sub(b.openedAt) >= b.config.Timeout {
		b.state = StateHalfOpen
		b.successes = 0
		b.logger.Info("CircuitHalfOpen", zap.String("service", "outbound"))
	}
	return b.state
}

// Success records a successful call. In half-open, consecutive successes
// close the circuit.
func (b *Breaker) Success() {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case StateHalfOpen:
		b.successes++
		b.failures = 0
		if b.successes >= b.config.SuccessThreshold {
			b.state = StateClosed
			b.logger.Info("CircuitClosed", zap.String("service", "outbound"))
		}
	case StateClosed:
		b.failures = 0
	}
}

// Failure records a failed call. In closed state, consecutive failures open
// the circuit; in half-open, a single failure re-opens it.
func (b *Breaker) Failure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case StateClosed:
		b.failures++
		if b.failures >= b.config.FailureThreshold {
			b.state = StateOpen
			b.openedAt = time.Now()
			b.logger.Warn("CircuitOpened", zap.String("service", "outbound"), zap.Int("failures", b.failures))
		}
	case StateHalfOpen:
		b.state = StateOpen
		b.openedAt = time.Now()
		b.logger.Warn("CircuitOpened", zap.String("service", "outbound"), zap.Int("failures", 1))
	}
}

// Reset forces the circuit back to closed (used on config reload / recovery).
func (b *Breaker) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.state = StateClosed
	b.failures = 0
	b.successes = 0
}

// Execute runs fn guarded by the breaker: it fails fast with ErrOpen when the
// circuit is open, and records Success/Failure based on the returned error.
func (b *Breaker) Execute(fn func() error) error {
	allowed, err := b.Allow()
	if err != nil {
		return err
	}
	if !allowed {
		return ErrOpen
	}
	callErr := fn()
	if callErr != nil {
		b.Failure()
		return callErr
	}
	b.Success()
	return nil
}
