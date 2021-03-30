package usermapper_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rebel-l/go-utils/osutils"
	"github.com/rebel-l/go-utils/testingutils"

	"github.com/rebel-l/auth-service/bootstrap"
	"github.com/rebel-l/auth-service/config"
	"github.com/rebel-l/auth-service/user/usermapper"
	"github.com/rebel-l/auth-service/user/usermodel"
	"github.com/rebel-l/auth-service/user/userstore"

	_ "github.com/mattn/go-sqlite3"
)

func setup(t *testing.T, name string) *sqlx.DB {
	t.Helper()

	// 0. init path
	storagePath := filepath.Join(".", "..", "..", "storage", "test_user", name)
	scriptPath := filepath.Join(".", "..", "..", "scripts", "sql", "sqlite")
	conf := &config.Database{
		StoragePath:       &storagePath,
		SchemaScriptsPath: &scriptPath,
	}

	// 1. clean up
	if osutils.FileOrPathExists(conf.GetStoragePath()) {
		if err := os.RemoveAll(conf.GetStoragePath()); err != nil {
			t.Fatalf("failed to cleanup test files: %v", err)
		}
	}

	// 2. init database
	db, err := bootstrap.Database(conf, "0.0.0", false)
	if err != nil {
		t.Fatalf("No error expected: %v", err)
	}

	return db
}

func prepareData(db *sqlx.DB, u *usermodel.User) (*usermodel.User, error) {
	var err error

	u.ID, err = uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate id: %w", err)
	}

	ctx := context.Background()
	q := db.Rebind(`
		INSERT INTO users (id, email, firstname, lastname, password, externalid, type) 
		VALUES (?, ?, ?, ?, ?, ?, ?);
	`)

	_, err = db.ExecContext(ctx, q, u.ID, u.EMail, u.FirstName, u.LastName, u.Password, u.ExternalID, u.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to create data: %w", err)
	}

	us := &userstore.User{}

	q = db.Rebind(`SELECT id, email, firstname, lastname, password, externalid, type, created_at, modified_at FROM users WHERE id = ?`)

	if err := db.GetContext(ctx, us, q, u.ID); err != nil {
		return nil, fmt.Errorf("failed to retrieve created data: %w", err)
	}

	return usermapper.StoreToModel(us), nil
}

func TestMapper_Load(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("long running test")
	}

	// 1. setup
	db := setup(t, "mapperLoad")

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatalf("unable to close database connection: %v", err)
		}
	})

	mapper := usermapper.New(db)

	// 2. test
	testCases := []struct {
		name        string
		id          uuid.UUID
		prepare     *usermodel.User
		expected    *usermodel.User
		expectedErr error
	}{
		{
			name: "success",
			prepare: &usermodel.User{
				EMail:      "matthewanderson358@example.org",
				FirstName:  "Anthony",
				LastName:   "Taylor",
				Password:   "yifLGn",
				ExternalID: "WMA9ywHAsAX3eSzF9ZcutvWlEH4ftVUWL",
				Type:       "9bDOBHM3PpR3HgYcMHZ8p6IoLSsk7OA",
			},
			expected: &usermodel.User{
				EMail:      "matthewanderson358@example.org",
				FirstName:  "Anthony",
				LastName:   "Taylor",
				Password:   "yifLGn",
				ExternalID: "WMA9ywHAsAX3eSzF9ZcutvWlEH4ftVUWL",
				Type:       "9bDOBHM3PpR3HgYcMHZ8p6IoLSsk7OA",
			},
		},
		{
			name:        "user not existing",
			id:          testingutils.UUIDParse(t, "ce49c1b0-1db9-435d-9093-8db98d070520"),
			expectedErr: usermapper.ErrNotFound,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var err error

			if testCase.prepare != nil {
				testCase.prepare, err = prepareData(db, testCase.prepare)
				if err != nil {
					t.Fatalf("failed to prepare data: %v", err)
				}

				testCase.id = testCase.prepare.ID
				testCase.expected.ID = testCase.prepare.ID
			}

			actual, err := mapper.Load(context.Background(), testCase.id)
			if !errors.Is(err, testCase.expectedErr) {
				t.Errorf("expected error '%v' but got '%v'", testCase.expectedErr, err)
			}

			assertUser(t, testCase.expected, actual)
		})
	}
}

func TestMapper_Save(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("long running test")
	}

	// 1. setup
	db := setup(t, "mapperSave")

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatalf("unable to close database connection: %v", err)
		}
	})

	mapper := usermapper.New(db)

	// 2. test
	testCases := []struct {
		name        string
		prepare     *usermodel.User
		actual      *usermodel.User
		expected    *usermodel.User
		expectedErr error
		duplicate   bool
	}{
		{
			name:        "model is nil",
			expectedErr: usermapper.ErrNoData,
		},
		{
			name: "model has no ID",
			actual: &usermodel.User{
				EMail:      "masonjohnson343@example.com",
				FirstName:  "Avery",
				LastName:   "Anderson",
				Password:   "DuETA7NDgFv",
				ExternalID: "AMGt7oFDdSsHMsUl1J7ZP8YvBWaa1SIuMjXT",
				Type:       "8ntIxmkTwNau47N9HyPBhvdNA",
			},
			expected: &usermodel.User{
				EMail:      "masonjohnson343@example.com",
				FirstName:  "Avery",
				LastName:   "Anderson",
				Password:   "DuETA7NDgFv",
				ExternalID: "AMGt7oFDdSsHMsUl1J7ZP8YvBWaa1SIuMjXT",
				Type:       "8ntIxmkTwNau47N9HyPBhvdNA",
			},
		},
		{
			name: "model has ID",
			prepare: &usermodel.User{
				EMail:      "elijahdavis704@test.org",
				FirstName:  "Matthew",
				LastName:   "Taylor",
				Password:   "r8IVnxK5xk",
				ExternalID: "soBzMYj6L44E",
				Type:       "9YCwHmf7K0uSYhHR0hZScdPse605b20pvrROtYSZt",
			},
			actual: &usermodel.User{
				EMail:      "elijahdavis704@test.org",
				FirstName:  "Matthew",
				LastName:   "Taylor",
				Password:   "r8IVnxK5xk",
				ExternalID: "soBzMYj6L44E",
				Type:       "9YCwHmf7K0uSYhHR0hZScdPse605b20pvrROtYSZt",
			},
			expected: &usermodel.User{
				EMail:      "elijahdavis704@test.org",
				FirstName:  "Matthew",
				LastName:   "Taylor",
				Password:   "r8IVnxK5xk",
				ExternalID: "soBzMYj6L44E",
				Type:       "9YCwHmf7K0uSYhHR0hZScdPse605b20pvrROtYSZt",
			},
		},
		{
			name: "update not existing model",
			actual: &usermodel.User{
				ID:         testingutils.UUIDParse(t, "3c3deabb-cd5b-4eb7-9a1a-2d198c518649"),
				EMail:      "aubreygarcia447@example.com",
				FirstName:  "Michael",
				LastName:   "Garcia",
				Password:   "ou2JkkxpxbcoVS63m25TcKdM26KgXZ2cM48E3OVIm8yfg7D1pb",
				ExternalID: "arlamIMQXH80ZnWK7GFysPu45kOj00EQdVoKM9obY9Y",
				Type:       "PgCDnyh5SohsxKCrGUdmR2muHcJw9vgW",
			},
			expectedErr: usermapper.ErrSaveToDB,
		},
		{
			name:      "model is duplicate EMail",
			duplicate: true,
			prepare: &usermodel.User{
				EMail:      "liambrown063@example.net",
				FirstName:  "Noah",
				LastName:   "Garcia",
				Password:   "VUXdzMUzSLQ3dyBe8SBeic8XxWXJ",
				ExternalID: "eZVEAdWHPjqa2pjM",
				Type:       "VS6as5wvfsFitI5Ivw8dUIHmsFGLvDbji5QsVGxiElF1",
			},
			actual: &usermodel.User{
				EMail:      "liambrown063@example.net",
				FirstName:  "Matthew",
				LastName:   "Harris",
				Password:   "979NCHQD3ByX0KHopzNXHwrNeCIU6q5KyO7DQV",
				ExternalID: "Bgm0wy1BLB4VHy281PRIElla0QMTNAeSey8",
				Type:       "zTTma6nn87scc77tT11u",
			},
			expectedErr: usermapper.ErrSaveToDB,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var err error

			if testCase.prepare != nil {
				testCase.prepare, err = prepareData(db, testCase.prepare)
				if err != nil {
					t.Fatalf("failed to prepare data: %v", err)
				}

				if !testCase.duplicate {
					testCase.actual.ID = testCase.prepare.ID
				}
			}

			res, err := mapper.Save(context.Background(), testCase.actual)
			if !errors.Is(err, testCase.expectedErr) {
				t.Errorf("expected error '%v' but got '%v'", testCase.expectedErr, err)
			}

			if res != nil && testCase.expected != nil {
				testCase.expected.ID = res.ID
			}

			assertUser(t, testCase.expected, res)
		})
	}
}

func TestMapper_Delete(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("long running test")
	}

	// 1. setup
	db := setup(t, "mapperDelete")

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatalf("unable to close database connection: %v", err)
		}
	})

	mapper := usermapper.New(db)

	// 2. test
	testCases := []struct {
		name        string
		id          uuid.UUID
		prepare     *usermodel.User
		expectedErr error
	}{
		{
			name: "success",
			prepare: &usermodel.User{
				EMail:      "anthonyanderson681@example.com",
				FirstName:  "Ethan",
				LastName:   "Miller",
				Password:   "ziTcFudZ1JtHxczHco8f",
				ExternalID: "2lEWN3Zgogb3sCm9WsXa3y6CPe",
				Type:       "e062Pami87p9",
			},
		},
		{
			name: "user not existing",
			id:   testingutils.UUIDParse(t, "3c19e139-9c6d-43c5-83e3-9247800f219d"),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var err error

			if testCase.prepare != nil {
				testCase.prepare, err = prepareData(db, testCase.prepare)
				if err != nil {
					t.Fatalf("failed to prepare data: %v", err)
				}

				testCase.id = testCase.prepare.ID
			}

			err = mapper.Delete(context.Background(), testCase.id)
			if !errors.Is(err, testCase.expectedErr) {
				t.Errorf("expected error '%v' but got '%v'", testCase.expectedErr, err)

				return
			}

			if testCase.expectedErr == nil {
				_, err = mapper.Load(context.Background(), testCase.id)
				if !errors.Is(err, usermapper.ErrNotFound) {
					t.Errorf("expected that user was deleted but got error '%v'", err)
				}
			}
		})
	}
}

func TestMapper_StoreToModel(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		actual   *userstore.User
		expected *usermodel.User
	}{
		{
			name:     "store is nil",
			expected: &usermodel.User{},
		},
		{
			name: "store has all attributes set",
			actual: &userstore.User{
				ID:         testingutils.UUIDParse(t, "2d6df424-7811-4415-9d9d-a4fe60d9429d"),
				EMail:      "jaydendavis233@test.org",
				FirstName:  "Mia",
				LastName:   "Thompson",
				Password:   "uPVMlN6BRTPgA",
				ExternalID: "SDZlvKGAbp9oPNFF",
				Type:       "2F4aEOVauUX3",
			},
			expected: &usermodel.User{
				ID:         testingutils.UUIDParse(t, "2d6df424-7811-4415-9d9d-a4fe60d9429d"),
				EMail:      "jaydendavis233@test.org",
				FirstName:  "Mia",
				LastName:   "Thompson",
				Password:   "uPVMlN6BRTPgA",
				ExternalID: "SDZlvKGAbp9oPNFF",
				Type:       "2F4aEOVauUX3",
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			u := usermapper.StoreToModel(testCase.actual)

			assertUser(t, testCase.expected, u)
		})
	}
}

func assertUser(t *testing.T, expected, actual *usermodel.User) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}

	if expected != nil && actual == nil || expected == nil && actual != nil {
		t.Errorf("expected user '%v' but got '%v'", expected, actual)

		return
	}

	if expected.ID != actual.ID {
		t.Errorf("expected ID %s but got %s", expected.ID, actual.ID)
	}

	if expected.EMail != actual.EMail {
		t.Errorf("expected EMail %s but got %s", expected.EMail, actual.EMail)
	}

	if expected.FirstName != actual.FirstName {
		t.Errorf("expected FirstName %s but got %s", expected.FirstName, actual.FirstName)
	}

	if expected.LastName != actual.LastName {
		t.Errorf("expected LastName %s but got %s", expected.LastName, actual.LastName)
	}

	if expected.Password != actual.Password {
		t.Errorf("expected Password %s but got %s", expected.Password, actual.Password)
	}

	if expected.ExternalID != actual.ExternalID {
		t.Errorf("expected ExternalID %s but got %s", expected.ExternalID, actual.ExternalID)
	}

	if expected.Type != actual.Type {
		t.Errorf("expected Type %s but got %s", expected.Type, actual.Type)
	}

	if expected.CreatedAt != actual.CreatedAt && actual.CreatedAt.IsZero() {
		t.Error("created at should be greater than the zero date")
	}

	if expected.ModifiedAt != actual.ModifiedAt && actual.ModifiedAt.IsZero() {
		t.Error("created at should be greater than the zero date")
	}
}
