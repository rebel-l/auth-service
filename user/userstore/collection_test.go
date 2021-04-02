package userstore_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/rebel-l/auth-service/user/userstore"
)

func TestFind(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("long running test")
	}

	// 1. setup
	db := setup(t, "storeFind")

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatalf("unable to close database connection: %v", err)
		}
	})

	// 2. test
	testCases := []struct {
		name        string
		prepare     userstore.UserCollection
		where       string
		args        []interface{}
		expected    userstore.UserCollection
		expectedErr error
	}{
		{
			name: "no matching entry",
			prepare: userstore.UserCollection{
				&userstore.User{
					EMail:     "c.schumann@gmx.de",
					FirstName: "Clara",
					LastName:  "Schumann",
					Type:      "standard",
				},
			},
			where: "id = ?",
			args:  []interface{}{"abc123"},
		},
		{
			name:        "wrong where",
			where:       "syntax",
			args:        []interface{}{"wrong"},
			expectedErr: userstore.ErrFindFail,
		},
		{
			name: "one matching entry",
			prepare: userstore.UserCollection{
				&userstore.User{
					EMail:     "f.chopin@gmx.fr",
					FirstName: "Frédéric",
					LastName:  "Chopin",
					Type:      "standard",
				},
			},
			where: "lastname = ?",
			args:  []interface{}{"Chopin"},
			expected: userstore.UserCollection{
				&userstore.User{
					EMail:     "f.chopin@gmx.fr",
					FirstName: "Frédéric",
					LastName:  "Chopin",
					Type:      "standard",
				},
			},
		},
		{
			name: "two matching entries",
			prepare: userstore.UserCollection{
				&userstore.User{
					EMail:     "f.liszt@gmx.de",
					FirstName: "Franz",
					LastName:  "Liszt",
					Type:      "standard",
				},
				&userstore.User{
					EMail:     "f.schuber@t-online.de",
					FirstName: "Franz",
					LastName:  "Schubert",
					Type:      "standard",
				},
			},
			where: "firstname = ?",
			args:  []interface{}{"Franz"},
			expected: userstore.UserCollection{
				&userstore.User{
					EMail:     "f.liszt@gmx.de",
					FirstName: "Franz",
					LastName:  "Liszt",
					Type:      "standard",
				},
				&userstore.User{
					EMail:     "f.schuber@t-online.de",
					FirstName: "Franz",
					LastName:  "Schubert",
					Type:      "standard",
				},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			if len(testCase.prepare) > 0 {
				prepare(t, db, testCase.prepare, testCase.expected)
			}

			actual, err := userstore.Find(context.Background(), db, testCase.where, testCase.args...)
			if !errors.Is(err, testCase.expectedErr) {
				t.Fatalf("expected error %v but got %v", testCase.expectedErr, err)
			}

			assertCollection(t, testCase.expected, actual)
		})
	}
}

func prepare(t *testing.T, db *sqlx.DB, users, expectedUsers userstore.UserCollection) {
	t.Helper()

	for _, user := range users {
		if err := user.Create(context.Background(), db); err != nil {
			t.Fatalf("failed to prepare data: %v", err)
		}

		for _, expectedUser := range expectedUsers {
			if user.EMail == expectedUser.EMail {
				expectedUser.ID = user.ID
			}
		}
	}

	return
}

func assertCollection(t *testing.T, expected, actual userstore.UserCollection) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Fatalf("expected %d entries in collection but got %d", len(expected), len(actual))
	}

	for k, _ := range expected {
		assertUser(t, expected[k], actual[k])
	}
}

func TestUserCollection_First(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		collection userstore.UserCollection
		expected   *userstore.User
	}{
		{
			name: "user collection is nil",
		},
		{
			name:       "empty user collection",
			collection: userstore.UserCollection{},
		},
		{
			name: "one entry in user collection",
			collection: userstore.UserCollection{
				&userstore.User{
					EMail: "l.v.beethoven@t-online.de",
				},
			},
			expected: &userstore.User{
				EMail: "l.v.beethoven@t-online.de",
			},
		},
		{
			name: "two entries in user collection",
			collection: userstore.UserCollection{
				&userstore.User{
					EMail: "j.s.bach@t-online.de",
				},
				&userstore.User{
					EMail: "r.wagner@gmx.de",
				},
			},
			expected: &userstore.User{
				EMail: "j.s.bach@t-online.de",
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actual := testCase.collection.First()

			if testCase.expected == nil && actual == nil {
				// nothing else needs to be checked
				return
			}

			if (testCase.expected != nil && actual == nil) || (testCase.expected == nil && actual != nil) {
				t.Errorf("expected %v but got %v", testCase.expected, actual)

				return
			}

			if testCase.expected.EMail != actual.EMail {
				t.Errorf("expected email %s but got %s", testCase.expected.EMail, actual.EMail)
			}
		})
	}
}
