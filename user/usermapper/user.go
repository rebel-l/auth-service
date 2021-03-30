package usermapper

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    "github.com/rebel-l/auth-service/user/usermodel"
    "github.com/rebel-l/auth-service/user/userstore"
    "github.com/rebel-l/go-utils/uuidutils"
)

var (
    // ErrLoadFromDB occurs if something went wrong on loading.
    ErrLoadFromDB = errors.New("failed to load user from database")

    // ErrNoData occurs if given model is nil.
    ErrNoData = errors.New("user is nil")

    // ErrSaveToDB occurs if something went wrong on saving.
    ErrSaveToDB = errors.New("failed to save user to database")

    // ErrDeleteFromDB occurs if something went wrong on deleting.
    ErrDeleteFromDB = errors.New("failed to delete user from database")

    // ErrNotFound occurs if record doesn't exist in database.
    ErrNotFound = errors.New("user was not found")
)

// Mapper provides methods to load and persist user models.
type Mapper struct {
    db *sqlx.DB
}

// New returns a new mapper.
func New(db *sqlx.DB) *Mapper {
    return &Mapper{db: db}
}

// Load returns a user model loaded from database by ID.
func (m *Mapper) Load(ctx context.Context, id uuid.UUID) (*usermodel.User, error) {
    s := &userstore.User{ID: id}

    if err := s.Read(ctx, m.db); errors.Is(err, sql.ErrNoRows) {
        return nil, ErrNotFound
    } else if err != nil {
        return nil, fmt.Errorf("%w: %v", ErrLoadFromDB, err)
    }

    return StoreToModel(s), nil
}

// Save persists (create or update) the model and returns the changed data (id, createdAt or modifiedAt).
func (m *Mapper) Save(ctx context.Context, model *usermodel.User) (*usermodel.User, error) {
    if model == nil {
        return nil, ErrNoData
    }

    s := modelToStore(model)

    if uuidutils.IsEmpty(model.ID) {
        if err := s.Create(ctx, m.db); err != nil {
            return nil, fmt.Errorf("%w: %v", ErrSaveToDB, err)
        }
    } else {
        if err := s.Update(ctx, m.db); err != nil {
            return nil, fmt.Errorf("%w: %v", ErrSaveToDB, err)
        }
    }

    model = StoreToModel(s)

    return model, nil
}

// Delete removes a model from database by ID.
func (m *Mapper) Delete(ctx context.Context, id uuid.UUID) error {
    s := &userstore.User{ID: id}
    if err := s.Delete(ctx, m.db); err != nil {
        return fmt.Errorf("%w: %v", ErrDeleteFromDB, err)
    }

    return nil
}

// StoreToModel returns a model based on the given store object. It maps all properties from store to model.
func StoreToModel(s *userstore.User) *usermodel.User {
    if s == nil {
        return &usermodel.User{}
    }

    return &usermodel.User{
        ID: s.ID,
        EMail: s.EMail,
        FirstName: s.FirstName,
        LastName: s.LastName,
        Password: s.Password,
        ExternalID: s.ExternalID,
        Type: s.Type,
        CreatedAt:  s.CreatedAt,
        ModifiedAt: s.ModifiedAt,
    }
}

// modelToStore returns a store based on the given model object. It maps all properties from model to store.
func modelToStore(m *usermodel.User) *userstore.User {
    return &userstore.User{
        ID: m.ID,
        EMail: m.EMail,
        FirstName: m.FirstName,
        LastName: m.LastName,
        Password: m.Password,
        ExternalID: m.ExternalID,
        Type: m.Type,
        CreatedAt:  m.CreatedAt,
        ModifiedAt: m.ModifiedAt,
    }
}
