// Package auth provides JWT token creation and validation.
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("token invalid")
	ErrBadSignature = errors.New("bad signature")
)

// Claims holds the JWT payload.
type Claims struct {
	UserID   string `json:"sub"`
	Username string `json:"name"`
}

// TokenManager creates and validates JWT tokens (HS256).
type TokenManager struct {
	secret    []byte
	expireSec int
}

// NewTokenManager creates a new TokenManager with the given secret and expiry.
func NewTokenManager(secret string, expireSec int) *TokenManager {
	return &TokenManager{
		secret:    []byte(secret),
		expireSec: expireSec,
	}
}

// Generate creates a signed JWT token for the given claims.
func (m *TokenManager) Generate(c Claims) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerBytes, _ := json.Marshal(header)
	headerB64 := base64URL(headerBytes)

	payload := map[string]interface{}{
		"sub":  c.UserID,
		"name": c.Username,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Duration(m.expireSec) * time.Second).Unix(),
	}
	payloadBytes, _ := json.Marshal(payload)
	payloadB64 := base64URL(payloadBytes)

	signingInput := headerB64 + "." + payloadB64
	signature := m.sign(signingInput)

	return signingInput + "." + signature, nil
}

// Validate parses and validates a JWT token, returning the claims.
func (m *TokenManager) Validate(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrTokenInvalid
	}

	signingInput := parts[0] + "." + parts[1]
	expectedSig := m.sign(signingInput)
	if parts[2] != expectedSig {
		return nil, ErrBadSignature
	}

	payloadBytes, err := base64URLDecode(parts[1])
	if err != nil {
		return nil, ErrTokenInvalid
	}

	var payload struct {
		Sub  string `json:"sub"`
		Name string `json:"name"`
		Exp  int64  `json:"exp"`
	}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, ErrTokenInvalid
	}

	if time.Now().Unix() > payload.Exp {
		return nil, ErrTokenExpired
	}

	return &Claims{
		UserID:   payload.Sub,
		Username: payload.Name,
	}, nil
}

func (m *TokenManager) sign(input string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(input))
	return base64URL(mac.Sum(nil))
}

func base64URL(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func base64URLDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

// GenerateTestToken is a convenience for tests.
func GenerateTestToken(userID, username string) string {
	m := NewTokenManager("test-secret", 3600)
	token, _ := m.Generate(Claims{UserID: userID, Username: username})
	return token
}

// errorf for internal logging, exposed for test checks if needed.
func errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
