package userstore

import (
    "context"
    "errors"
    "time"
    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    "github.com/rebel-l/go-utils/uuidutils"
)

var (
    // ErrIDMissing will be thrown if an ID is expected but not set.
    ErrIDMissing = errors.New("id is mandatory for this operation")

    // ErrCreatingID will be thrown if creating an ID failed.
    ErrCreatingID = errors.New("id creation failed")

    // ErrIDIsSet will be thrown if no ID is expected but already set.
    ErrIDIsSet = errors.New("id should be not set for this operation, use update instead")

    // ErrDataMissing will be thrown if mandatory data is not set.
    ErrDataMissing = errors.New("no data or mandatory data missing")
)

// User represents the user in the database.
type User struct {
    ID	uuid.UUID	`db:"id"`
    EMail	string	`db:"email"`
    FirstName	string	`db:"firstname"`
    LastName	string	`db:"lastname"`
    Password	string	`db:"password"`
    ExternalID	string	`db:"externalid"`
    Type	string	`db:"type"`
    CreatedAt  time.Time `db:"created_at"`
    ModifiedAt time.Time `db:"modified_at"`
}

// Create creates current object in the database.
func (u *User) Create(ctx context.Context, db *sqlx.DB) error {
    if !u.IsValid() {
        return ErrDataMissing
    }

    if !uuidutils.IsEmpty(u.ID) {
        return ErrIDIsSet
    }

    var err error

    u.ID, err = uuid.NewRandom()
    if err != nil {
        return fmt.Errorf("%w: %v", ErrCreatingID, err)
    }

    q := db.Rebind(`
		INSERT INTO users (id, email, firstname, lastname, password, externalid, type) 
		VALUES (?, ?, ?, ?, ?, ?, ?);
	`)

    _, err = db.ExecContext(ctx, q, u.ID, u.EMail, u.FirstName, u.LastName, u.Password, u.ExternalID, u.Type)
    if err != nil {
        return fmt.Errorf("failed to create: %w", err)
    }

    return u.Read(ctx, db)
}

// Read sets the user from database by given ID.
func (u *User) Read(ctx context.Context, db *sqlx.DB) error {
    if u == nil || uuidutils.IsEmpty(u.ID) {
        return ErrIDMissing
    }

    q := db.Rebind(`
        SELECT id, email, firstname, lastname, password, externalid, type, created_at, modified_at
        FROM users
        WHERE id = ?;
    `)
    if err := db.GetContext(ctx, u, q, u.ID); err != nil {
        return fmt.Errorf("failed to read: %w", err)
    }

    return nil
}

// Update changes the current object on the database by ID.
func (u *User) Update(ctx context.Context, db *sqlx.DB) error {
    if !u.IsValid() {
        return ErrDataMissing
    }

    if uuidutils.IsEmpty(u.ID) {
        return ErrIDMissing
    }

    q := db.Rebind(`
		UPDATE users 
		SET email = ?, firstname = ?, lastname = ?, password = ?, externalid = ?, type = ? 
		WHERE id = ?;
	`)

    if _, err := db.ExecContext(ctx, q, u.EMail, u.FirstName, u.LastName, u.Password, u.ExternalID, u.Type, u.ID); err != nil {
        return fmt.Errorf("failed to update: %w", err)
    }

    return u.Read(ctx, db)
}

// Delete removes the current object from database by its ID.
func (u *User) Delete(ctx context.Context, db *sqlx.DB) error {
    if u == nil || uuidutils.IsEmpty(u.ID) {
        return ErrIDMissing
    }

    q := db.Rebind(`
        DELETE FROM users
        WHERE id = ?
    `)

    if _, err := db.ExecContext(ctx, q, u.ID); err != nil {
        return fmt.Errorf("failed to delete: %w", err)
    }

    return nil
}

// IsValid returns true if all mandatory fields are set.
func (u *User) IsValid() bool {
    if u == nil || u.EMail == "" || u.FirstName == "" || u.LastName == "" || u.Type == "" {
        return false
    }

    return true
}
