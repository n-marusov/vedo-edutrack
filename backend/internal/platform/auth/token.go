package auth

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// TokenClaims is the JWT payload for dev-issued tokens.
type TokenClaims struct {
	UserID string   `json:"sub"`
	Roles  []string `json:"roles"`
	Issuer string   `json:"iss"`
	Aud    string   `json:"aud"`
	Expiry time.Time
}

// IssueToken signs a JWT with the given claims using RS256.
//
// The token carries sub, roles, iss, aud, iat, exp (default 24h).
func IssueToken(key *rsa.PrivateKey, userID string, roles []string, issuer, audience string, ttl time.Duration) (string, error) {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	now := time.Now()

	tok, err := jwt.NewBuilder().
		Subject(userID).
		Issuer(issuer).
		Audience([]string{audience}).
		IssuedAt(now).
		Expiration(now.Add(ttl)).
		Claim("roles", roles).
		Build()
	if err != nil {
		return "", fmt.Errorf("build token: %w", err)
	}

	headers := jws.NewHeaders()
	if err := headers.Set(jws.KeyIDKey, "dev-key"); err != nil {
		return "", fmt.Errorf("set kid header: %w", err)
	}

	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256, key, jws.WithProtectedHeaders(headers)))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return string(signed), nil
}
