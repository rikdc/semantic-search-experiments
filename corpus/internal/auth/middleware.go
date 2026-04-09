package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Authenticate reads the Bearer token from the Authorization header and validates it.
// Returns the parsed Token on success, or an error if the header is missing, malformed,
// or the token fails signature and expiry checks.
func Authenticate(r *http.Request, secret []byte) (Token, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return Token{}, errors.New("missing Authorization header")
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return Token{}, errors.New("malformed Authorization header: expected 'Bearer <token>'")
	}

	if parts[1] == "" {
		return Token{}, errors.New("empty token string")
	}

	tok := Token{UserID: parts[1]}
	if err := ValidateToken(tok, secret); err != nil {
		return Token{}, err
	}
	return tok, nil
}

// Authorize checks that the authenticated user holds at least one of the required roles.
// Roles are encoded as a prefix on the UserID field (e.g. "admin:alice").
// Returns an error if none of the allowed roles match.
func Authorize(tok Token, allowedRoles ...string) error {
	for _, role := range allowedRoles {
		if strings.HasPrefix(tok.UserID, role+":") {
			return nil
		}
	}
	return errors.New("user does not have the required permission")
}
