package userstore_test

import (
	"testing"

	"github.com/rebel-l/auth-service/user/userstore"
)

func TestFind(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("long running test")
	}

	t.Skip("to be implemented!")
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
