package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const (
	headerKeyAuth = "Authorization"
)

// ErrNoJWT indicates that the JWT is not defined.
var ErrNoJWT = errors.New("JWT should not be empty")

// Token represents a token including the signed JWT token.
type Token struct {
	ID      uuid.UUID
	JWT     string
	Expires time.Time
}

// NewToken returns a new token struct.
func NewToken(id uuid.UUID, expires time.Time, token *jwt.Token, secret string) (*Token, error) {
	if token == nil {
		return nil, ErrNoJWT
	}

	t, err := token.SignedString([]byte(secret)) // TODO: extend secret with a random salt
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &Token{
		ID:      id,
		JWT:     t,
		Expires: expires,
	}, nil
}

func extractToken(header http.Header) string {
	token := header.Get(headerKeyAuth)

	parts := strings.Split(token, " ")
	if len(parts) > 1 {
		return parts[1]
	}

	return ""
}
