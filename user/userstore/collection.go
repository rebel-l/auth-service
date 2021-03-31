package userstore

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserCollection []*User // TODO: generator

func (c UserCollection) First() *User {
	if len(c) > 0 {
		return c[0]
	}

	return nil
}

func Find(ctx context.Context, db *sqlx.DB, where string, args ...interface{}) (UserCollection, error) { // TODO: generator
	var collection UserCollection

	q := fmt.Sprintf(`
		SELECT id, email, firstname, lastname, password, externalid, type, created_at, modified_at
        FROM users
        WHERE %s;
	`, where) // TODO: have statement at a central place (constant)

	if err := db.SelectContext(ctx, &collection, db.Rebind(q), args...); err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err) // TODO: in generator replace 'users'
	}

	return collection, nil
}
