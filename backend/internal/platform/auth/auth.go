package auth

import (
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

// DefaultKeyPath is the default location for the dev private key.
const DefaultKeyPath = "./tmp/jwt_key.pem"

// Auth bundles the local JWT infrastructure: the RSA key pair, JWKS
// serialization, token issuance and the validation middleware.
type Auth struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	Kid        string
	Issuer     string
	Audience   string
	Logger     *zap.Logger
}

// New creates an Auth from the configured key path (auto-generates if missing).
func New(keyPath, issuer, audience string, logger *zap.Logger) (*Auth, error) {
	if keyPath == "" {
		keyPath = envOrDefault("JWKS_PRIVATE_KEY_PATH", DefaultKeyPath)
	}
	key, err := LoadOrCreateKey(keyPath)
	if err != nil {
		return nil, err
	}
	if issuer == "" {
		issuer = "vedo-edutrack"
	}
	if audience == "" {
		audience = "vedo-edutrack"
	}
	return &Auth{
		PrivateKey: key,
		PublicKey:  &key.PublicKey,
		Kid:        "dev-key",
		Issuer:     issuer,
		Audience:   audience,
		Logger:     logger,
	}, nil
}

// NewNop creates an Auth without persistence (tests).
func NewNop(logger *zap.Logger) *Auth {
	a, _ := New("", "vedo-edutrack", "vedo-edutrack", logger)
	return a
}

// JWKSHandler serves the public key as a JWKS document.
func (a *Auth) JWKSHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		doc, err := MarshalJWKS(a.PublicKey, a.Kid)
		if err != nil {
			http.Error(w, `{"error":"jwks_serialization"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(doc)
	}
}

// TokenHandler issues a dev JWT from {"user_id", "roles"}.
// WARNs when used with ENV=production.
func (a *Auth) TokenHandler() http.HandlerFunc {
	type tokenRequest struct {
		UserID string   `json:"user_id"`
		Roles  []string `json:"roles"`
	}
	type tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("APP_ENV") == "production" {
			a.Logger.Warn("dev token endpoint used with ENV=production")
		}
		var req tokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid_request"}`, http.StatusBadRequest)
			return
		}
		if req.UserID == "" {
			http.Error(w, `{"error":"user_id_required"}`, http.StatusBadRequest)
			return
		}
		if req.Roles == nil {
			req.Roles = []string{"learner"}
		}

		ttl := 24 * time.Hour
		token, err := IssueToken(a.PrivateKey, req.UserID, req.Roles, a.Issuer, a.Audience, ttl)
		if err != nil {
			a.Logger.Error("token issuance failed", zap.Error(err))
			http.Error(w, `{"error":"token_issuance_failed"}`, http.StatusInternalServerError)
			return
		}

		a.Logger.Info("dev token issued", zap.String("user_id", req.UserID))
		writeJSONResponse(w, http.StatusOK, tokenResponse{
			AccessToken: token,
			TokenType:   "Bearer",
			ExpiresIn:   int64(ttl.Seconds()),
		})
	}
}

// Middleware returns the JWT validation middleware.
func (a *Auth) Middleware() func(http.Handler) http.Handler {
	return Middleware(a.PublicKey, a.Kid, a.Issuer, a.Audience)
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func writeJSONResponse(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
