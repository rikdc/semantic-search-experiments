package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"
)

var (
	ErrTokenExpired = errors.New("token has expired")
	ErrTokenInvalid = errors.New("token signature is invalid")
)

// Token holds the claims and HMAC signature for an authenticated session.
type Token struct {
	UserID    string
	ExpiresAt time.Time
	Signature []byte
}

// ValidateToken checks that the token signature is correct and the token has not expired.
// Returns ErrTokenExpired or ErrTokenInvalid if validation fails.
func ValidateToken(tok Token, secret []byte) error {
	if time.Now().After(tok.ExpiresAt) {
		return ErrTokenExpired
	}

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(tok.UserID))
	mac.Write([]byte(tok.ExpiresAt.String()))
	expected := mac.Sum(nil)

	if !hmac.Equal(tok.Signature, expected) {
		return ErrTokenInvalid
	}
	return nil
}

// RefreshToken issues a new token with a fresh expiry, provided the existing token is valid.
// The original token must pass ValidateToken before a refresh is granted.
func RefreshToken(tok Token, secret []byte, ttl time.Duration) (Token, error) {
	if err := ValidateToken(tok, secret); err != nil {
		return Token{}, fmt.Errorf("refresh denied: %w", err)
	}

	fresh := Token{
		UserID:    tok.UserID,
		ExpiresAt: time.Now().Add(ttl),
	}

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(fresh.UserID))
	mac.Write([]byte(fresh.ExpiresAt.String()))
	fresh.Signature = mac.Sum(nil)

	return fresh, nil
}

// RevokeToken marks a token as invalid by clearing its signature.
// The returned token will fail ValidateToken immediately.
func RevokeToken(tok Token) Token {
	tok.Signature = nil
	return tok
}
