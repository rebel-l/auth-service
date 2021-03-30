package usermodel

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "time"
    "github.com/google/uuid"
)

// ErrDecodeJSON occurs if the a string is not in JSON format.
var ErrDecodeJSON = errors.New("failed to decode JSON")

// User represents a model of repository including business logic.
type User struct {
    ID	uuid.UUID	`json:"ID"`
    EMail	string	`json:"EMail"`
    FirstName	string	`json:"FirstName"`
    LastName	string	`json:"LastName"`
    Password	string	`json:"Password"`
    ExternalID	string	`json:"ExternalID"`
    Type	string	`json:"Type"`
    CreatedAt  time.Time `json:"CreatedAt"`
    ModifiedAt time.Time `json:"ModifiedAt"`
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
