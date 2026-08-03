// Package ratelimit implements token-bucket API rate limiting (F6).
//
// Limits are enforced per API key / user with a configurable rate (e.g.
// 10 req/min for SPARQL, 100 req/min for REST) and a small burst allowance.
// Exceeded requests receive 429 with a Retry-After header
// (REQ-NFR-security.compliance.owasp-top10: rate limiting).
package ratelimit

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
)

// bucket is one token bucket (per limiter key).
type bucket struct {
	tokens   float64
	last     time.Time
	capacity float64
	rate     float64 // tokens per second
}

// refill adds tokens up to capacity based on elapsed time.
func (b *bucket) refill(now time.Time) {
	elapsed := now.Sub(b.last).Seconds()
	if elapsed <= 0 {
		return
	}
	b.tokens += elapsed * b.rate
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
	b.last = now
}

// take consumes one token; reports whether the request is allowed.
func (b *bucket) take(now time.Time) bool {
	b.refill(now)
	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}

// retryAfter returns the seconds until at least one token is available.
func (b *bucket) retryAfter(now time.Time) int {
	if b.rate <= 0 {
		return 1
	}
	needed := 1 - b.tokens
	seconds := int(needed / b.rate)
	if seconds < 1 {
		seconds = 1
	}
	return seconds
}

// Limiter is a concurrency-safe token bucket limiter keyed by API key/user.
type Limiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	rate     float64 // tokens per second
	capacity float64 // burst allowance
	logger   *zap.Logger
}

// NewLimiter builds a token bucket limiter with the given rate (tokens per
// second) and burst capacity. rate=0 disables limiting (pass-through).
func NewLimiter(ratePerSecond float64, burst int, logger *zap.Logger) *Limiter {
	if logger == nil {
		logger = zap.NewNop()
	}
	if burst < 1 {
		burst = 1
	}
	return &Limiter{
		buckets:  map[string]*bucket{},
		rate:     ratePerSecond,
		capacity: float64(burst),
		logger:   logger.Named("ratelimit"),
	}
}

// Allow reports whether a request from the given key may proceed. When
// denied, retryAfter returns the seconds until the next token.
func (l *Limiter) Allow(key string) (allowed bool, retryAfter int) {
	if l.rate <= 0 {
		return true, 0
	}
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()

	b, ok := l.buckets[key]
	if !ok {
		b = &bucket{capacity: l.capacity, rate: l.rate, last: now, tokens: l.capacity}
		l.buckets[key] = b
	}
	if b.take(now) {
		return true, 0
	}
	return false, b.retryAfter(now)
}

// Middleware returns a chi-style middleware enforcing the limiter for the
// given key extractor. Denied requests receive 429 + Retry-After.
func (l *Limiter) Middleware(keyFn func(r *http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFn(r)
			if key == "" {
				key = remoteAddr(r)
			}
			allowed, retryAfter := l.Allow(key)
			if !allowed {
				l.logger.Warn("RateLimited",
					zap.String("userID", key),
					zap.String("endpoint", r.URL.Path),
					zap.Int("retryAfter", retryAfter),
				)
				w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"error":"rate_limited","message":"too many requests"}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// remoteAddr returns the request's remote address (fallback key).
func remoteAddr(r *http.Request) string {
	host := r.RemoteAddr
	if i := len(host); i > 0 && host[0] == '[' {
		if j := indexByte(host, ']'); j >= 0 {
			return host[1:j]
		}
	}
	for i := 0; i < len(host); i++ {
		if host[i] == ':' {
			return host[:i]
		}
	}
	return host
}

// indexByte finds the first occurrence of c in s (-1 when absent).
func indexByte(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}
