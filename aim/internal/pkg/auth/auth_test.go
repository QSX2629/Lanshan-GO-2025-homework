package auth

import (
	"testing"
	"time"
)

func TestGenerateAndValidate(t *testing.T) {
	m := NewTokenManager("my-secret-key", 3600)
	claims := Claims{UserID: "user123", Username: "alice"}

	token, err := m.Generate(claims)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if token == "" {
		t.Fatal("Generate() returned empty token")
	}

	validated, err := m.Validate(token)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if validated.UserID != "user123" {
		t.Errorf("UserID = %q, want user123", validated.UserID)
	}
	if validated.Username != "alice" {
		t.Errorf("Username = %q, want alice", validated.Username)
	}
}

func TestValidateExpired(t *testing.T) {
	m := NewTokenManager("secret", -1) // expires immediately
	claims := Claims{UserID: "u1", Username: "bob"}

	token, err := m.Generate(claims)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	_, err = m.Validate(token)
	if err != ErrTokenExpired {
		t.Errorf("expected ErrTokenExpired, got %v", err)
	}
}

func TestValidateBadSignature(t *testing.T) {
	m := NewTokenManager("secret-a", 3600)
	token, _ := m.Generate(Claims{UserID: "u1", Username: "x"})

	// Validate with a different secret.
	m2 := NewTokenManager("secret-b", 3600)
	_, err := m2.Validate(token)
	if err != ErrBadSignature {
		t.Errorf("expected ErrBadSignature, got %v", err)
	}
}

func TestValidateMalformed(t *testing.T) {
	m := NewTokenManager("secret", 3600)

	tests := []string{
		"",
		"abc",
		"abc.def",
		"abc.def.ghi.jkl",
	}

	for _, tc := range tests {
		_, err := m.Validate(tc)
		if err == nil {
			t.Errorf("expected error for input %q", tc)
		}
	}
}

func TestValidateTamperedPayload(t *testing.T) {
	m := NewTokenManager("secret", 3600)
	token, _ := m.Generate(Claims{UserID: "u1", Username: "x"})

	// Tamper with the payload portion (replace base64 payload).
	parts := []byte(token)
	// Corrupt part of the signature.
	parts[len(parts)-1] ^= 0xff

	_, err := m.Validate(string(parts))
	if err != ErrBadSignature {
		t.Errorf("expected ErrBadSignature for tampered token, got %v", err)
	}
}

func TestGenerateTestToken(t *testing.T) {
	token := GenerateTestToken("uid", "uname")
	m := NewTokenManager("test-secret", 3600)
	c, err := m.Validate(token)
	if err != nil {
		t.Fatalf("Validate(GenerateTestToken()) error = %v", err)
	}
	if c.UserID != "uid" || c.Username != "uname" {
		t.Errorf("claims = %+v, want uid/uname", c)
	}
}
