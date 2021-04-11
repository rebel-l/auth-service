package usermodel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/rebel-l/auth-service/facebookapi"

	"github.com/google/uuid"
)

// ErrDecodeJSON occurs if the a string is not in JSON format.
var ErrDecodeJSON = errors.New("failed to decode JSON")

// User represents a model of repository including business logic.
type User struct {
	ID         uuid.UUID `json:"ID"`
	EMail      string    `json:"EMail"`
	FirstName  string    `json:"FirstName"`
	LastName   string    `json:"LastName"`
	Password   string    `json:"-"`
	ExternalID string    `json:"ExternalID"`
	Type       string    `json:"Type"`
	CreatedAt  time.Time `json:"CreatedAt"`
	ModifiedAt time.Time `json:"ModifiedAt"`
}

// NewFromFacebook returns a new user model out of a facebook user.
func NewFromFacebook(fb *facebookapi.User) *User {
	return &User{
		EMail:      fb.EMail,
		FirstName:  fb.FirstName,
		LastName:   fb.LastName,
		ExternalID: fb.ID,
		Type:       TypeFacebook,
	}
}

// DecodeJSON converts JSON data to struct.
func (r *User) DecodeJSON(reader io.Reader) error {
	if r == nil {
		return nil
	}

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(r); err != nil {
		return fmt.Errorf("%w: %v", ErrDecodeJSON, err)
	}

	return nil
}
