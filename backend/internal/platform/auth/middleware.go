package auth

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// ErrUnauthorized is returned when a request has no or an invalid token.
var ErrUnauthorized = errors.New("unauthorized")

// contextKey is the context key for auth claims.
type contextKey struct{}

// Claims carries the validated token claims injected into the request context.
type Claims struct {
	UserID string
	Roles  []string
	Raw    jwt.Token
}

// Middleware validates Bearer JWTs against the local JWKS (RS256).
//
// Requests without a token receive 401. Public paths (/healthz, /readyz,
// /.well-known/jwks.json) are skipped. Validated claims are injected into
// the request context via WithClaims.
func Middleware(key *rsa.PublicKey, kid string, issuer, audience string) func(http.Handler) http.Handler {
	// Build a static validation set from the local public key. Remote JWKS
	// (JWKS_URL) validation is a post-MVP concern.
	pubJWK, err := publicKeyToJWK(key, kid)
	if err != nil {
		// Fall back to a nil set: every request fails closed with 401.
		pubJWK = nil
	}
	var keySet jwk.Set
	if pubJWK != nil {
		keySet = jwk.NewSet()
		_ = keySet.AddKey(pubJWK)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip public endpoints.
			if isPublicPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			token, err := extractBearer(r)
			if err != nil {
				writeUnauthorized(w, "missing or malformed Authorization header")
				return
			}

			claims, err := validate(token, keySet, issuer, audience)
			if err != nil {
				writeUnauthorized(w, "invalid jwt: "+err.Error())
				return
			}

			next.ServeHTTP(w, r.WithContext(WithClaims(r.Context(), claims)))
		})
	}
}

// isPublicPath reports whether the path is exempt from authentication.
func isPublicPath(path string) bool {
	switch path {
	case "/healthz", "/readyz", "/.well-known/jwks.json", "/metrics":
		return true
	}
	return false
}

// extractBearer pulls the Bearer token from the Authorization header.
func extractBearer(r *http.Request) (string, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return "", ErrUnauthorized
	}
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", ErrUnauthorized
	}
	return parts[1], nil
}

// validate parses and verifies a JWT against the key set, checking
// signature (RS256), issuer, audience and expiry.
func validate(token string, keySet jwk.Set, issuer, audience string) (*Claims, error) {
	parsed, err := jwt.Parse([]byte(token),
		jwt.WithKeySet(keySet),
		jwt.WithValidate(true),
		jwt.WithIssuer(issuer),
		jwt.WithAudience(audience),
	)
	if err != nil {
		return nil, err
	}

	roles := []string{}
	if raw, ok := parsed.Get("roles"); ok {
		if arr, ok := raw.([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok {
					roles = append(roles, s)
				}
			}
		}
	}

	return &Claims{
		UserID: parsed.Subject(),
		Roles:  roles,
		Raw:    parsed,
	}, nil
}

// WithClaims injects validated claims into the context.
func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, contextKey{}, claims)
}

// ClaimsFrom extracts claims from the context (nil if absent).
func ClaimsFrom(ctx context.Context) *Claims {
	claims, _ := ctx.Value(contextKey{}).(*Claims)
	return claims
}

// writeUnauthorized writes a 401 JSON error.
func writeUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("WWW-Authenticate", `Bearer realm="edutrack"`)
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized", "message": msg})
}
