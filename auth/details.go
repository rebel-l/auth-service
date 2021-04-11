package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/rebel-l/go-utils/uuidutils"

	"github.com/google/uuid"
)

const (
	headerKeyAuth = "Authorization"
)

// ErrInvalidData indicates that the JWT is not defined.
var ErrInvalidData = errors.New("data for details are invalid")

// Details represents a token including the signed JWT token.
type Details struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	Expires time.Time
}

// NewToken returns a new token struct.
func NewToken(id uuid.UUID, userID uuid.UUID, expires time.Time) (*Details, error) {
	if uuidutils.IsEmpty(id) {
		return nil, fmt.Errorf("%w: tokenID", ErrInvalidData)
	}

	if uuidutils.IsEmpty(userID) {
		return nil, fmt.Errorf("%w: userID", ErrInvalidData)
	}

	if expires.IsZero() || expires.Before(time.Now()) {
		return nil, fmt.Errorf("%w: expiration time", ErrInvalidData)
	}

	return &Details{
		ID:      id,
		UserID:  userID,
		Expires: expires,
	}, nil
}
