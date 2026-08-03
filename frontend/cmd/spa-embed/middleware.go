package main

import (
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// --- Rate Limiter (token bucket) ---

type tokenBucket struct {
	rate     int       // tokens per second
	burst    int       // max burst
	tokens   float64   // current tokens
	lastFill time.Time // last token refill
	mu       sync.Mutex
}

func newTokenBucket(rate, burst int) *tokenBucket {
	return &tokenBucket{
		rate:     rate,
		burst:    burst,
		tokens:   float64(burst),
		lastFill: time.Now(),
	}
}

func (tb *tokenBucket) allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastFill).Seconds()
	tb.tokens += elapsed * float64(tb.rate)
	if tb.tokens > float64(tb.burst) {
		tb.tokens = float64(tb.burst)
	}
	tb.lastFill = now

	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}
	return false
}

// rateLimitMiddleware limits request rate using a token bucket.
func rateLimitMiddleware(rate, burst int) middleware {
	bucket := newTokenBucket(rate, burst)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !bucket.allow() {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// middleware wraps an http.Handler with additional behaviour.
type middleware func(http.Handler) http.Handler

// chain applies middlewares in order (outermost first).
func chain(h http.Handler, mws ...middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// recoveryMiddleware catches panics and returns 500.
func recoveryMiddleware(logger *slog.Logger) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered",
						"error", rec,
						"path", r.URL.Path,
						"method", r.Method,
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// requestIDMiddleware injects a unique request ID via context and X-Request-ID header.
func requestIDMiddleware(logger *slog.Logger) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = newRequestID()
			}
			w.Header().Set("X-Request-ID", reqID)
			next.ServeHTTP(w, r)
		})
	}
}

// loggingMiddleware logs every request with structured fields.
func loggingMiddleware(logger *slog.Logger) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(wrapped, r)
			logger.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrapped.statusCode,
				"duration_ms", time.Since(start).Milliseconds(),
				"remote_addr", r.RemoteAddr,
			)
		})
	}
}

// securityHeadersMiddleware sets OWASP-recommended security headers.
func securityHeadersMiddleware() middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Strict-Transport-Security", "max-age=63072000")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			next.ServeHTTP(w, r)
		})
	}
}

// cacheMiddleware adds Cache-Control headers based on path.
//   - /assets/*  → immutable, max-age=31536000 (1 year)
//   - everything else → no-cache (ETag/revalidation)
func cacheMiddleware() middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(r.URL.Path) >= 8 && r.URL.Path[:8] == "/assets/" {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			} else {
				w.Header().Set("Cache-Control", "no-cache")
			}
			next.ServeHTTP(w, r)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
