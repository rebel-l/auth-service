package userstore

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// ErrFindFail indicates that find failed.
var ErrFindFail = errors.New("failed to find users")

// UserCollection is a collection of zero to many users.
// nolint: godox
type UserCollection []*User // TODO: generator

// First returns the first user of the collection. Returns nil if there aren't any users in the collection.
func (c UserCollection) First() *User {
	if len(c) > 0 {
		return c[0]
	}

	return nil
}

// Find returns a collection of users depending on the where statement. Ensure to escape parameters in where-clause
// with question marks (?) and provide the values in the correct order as last parameters.
// nolint: godox
// TODO: generator.
func Find(ctx context.Context, db *sqlx.DB, where string, args ...interface{}) (UserCollection, error) {
	var collection UserCollection

	q := db.Rebind(selectWhere(where))

	if err := db.SelectContext(ctx, &collection, q, args...); err != nil {
		// nolint: godox
		return nil, fmt.Errorf("%w: %v", ErrFindFail, err) // TODO: in generator replace 'users'
	}

	return collection, nil
}
