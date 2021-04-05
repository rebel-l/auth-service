package auth

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// ErrNoJWT indicates that the JWT is not defined.
var ErrNoJWT = errors.New("JWT should not be empty")

// Token represents a token including the signed JWT token.
type Token struct {
	ID      uuid.UUID
	JWT     string
	Expires int64
}

// NewToken returns a new token struct.
func NewToken(id uuid.UUID, expires int64, token *jwt.Token, secret string) (*Token, error) {
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
