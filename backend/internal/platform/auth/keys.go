// Package auth provides local JWT authentication with a self-signed RS256
// key pair for development environments.
//
// Keycloak is post-MVP (ADR-DES.SECURITY.rbac-model); local RS256/JWKS is the
// MVP path. The public key is served at /.well-known/jwks.json; tokens are
// signed with the private key and validated via JWKS.
package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

// GenerateKeyPair creates a new 2048-bit RSA key pair.
func GenerateKeyPair() (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("generate RSA key: %w", err)
	}
	return key, nil
}

// LoadOrCreateKey loads the private key from JWKS_PRIVATE_KEY_PATH, or
// generates a new one and persists it if the file does not exist.
func LoadOrCreateKey(path string) (*rsa.PrivateKey, error) {
	if path != "" {
		if key, err := loadKeyFile(path); err == nil {
			return key, nil
		}
	}

	key, err := GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	if path != "" {
		if err := saveKeyFile(path, key); err != nil {
			return nil, fmt.Errorf("persist generated key: %w", err)
		}
	}
	return key, nil
}

// loadKeyFile reads and parses a PEM-encoded RSA private key.
func loadKeyFile(path string) (*rsa.PrivateKey, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, fmt.Errorf("no PEM block in %s", path)
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse PKCS1 key: %w", err)
	}
	return key, nil
}

// saveKeyFile persists the private key as PEM (0600) and creates parent dirs.
func saveKeyFile(path string, key *rsa.PrivateKey) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}
	raw := pem.EncodeToMemory(block)
	return os.WriteFile(path, raw, 0o600)
}

// MarshalJWKS serializes the public key as a JWKS JSON document
// (kid, kty=RSA, n, e, alg=RS256).
func MarshalJWKS(publicKey *rsa.PublicKey, kid string) ([]byte, error) {
	n := base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(intToBytes(publicKey.E))

	doc := map[string]interface{}{
		"keys": []interface{}{
			map[string]interface{}{
				"kid": kid,
				"kty": "RSA",
				"n":   n,
				"e":   e,
				"alg": "RS256",
				"use": "sig",
			},
		},
	}
	return json.Marshal(doc)
}

// intToBytes converts a small int to a big-endian byte slice.
func intToBytes(v int) []byte {
	if v == 0 {
		return []byte{0}
	}
	var out []byte
	for v > 0 {
		out = append([]byte{byte(v & 0xff)}, out...)
		v >>= 8
	}
	return out
}

// publicKeyToJWK converts an RSA public key to a jwx JWK (used by the
// middleware's validation set).
func publicKeyToJWK(publicKey *rsa.PublicKey, kid string) (jwk.Key, error) {
	key, err := jwk.FromRaw(publicKey)
	if err != nil {
		return nil, fmt.Errorf("convert key to JWK: %w", err)
	}
	if err := key.Set(jwk.KeyIDKey, kid); err != nil {
		return nil, fmt.Errorf("set kid: %w", err)
	}
	if err := key.Set(jwk.AlgorithmKey, "RS256"); err != nil {
		return nil, fmt.Errorf("set alg: %w", err)
	}
	return key, nil
}
