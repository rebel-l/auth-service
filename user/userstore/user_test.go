package userstore_test

import (
    "context"
    "database/sql"
    "errors"
    "os"
    "path/filepath"
    "testing"
    "time"
    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    "github.com/rebel-l/auth-service/bootstrap"
    "github.com/rebel-l/auth-service/config"
    "github.com/rebel-l/auth-service/user/userstore"
    "github.com/rebel-l/go-utils/osutils"
    "github.com/rebel-l/go-utils/testingutils"
    "github.com/rebel-l/go-utils/uuidutils"

    _ "github.com/mattn/go-sqlite3"
)

func setup(t *testing.T, name string) *sqlx.DB {
    t.Helper()

    // 0. init path
    storagePath := filepath.Join(".", "..", "..", "storage", "test_user", name)

    // nolint: godox
    // TODO: change that it works with other dialects like postgres
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

func TestUser_Create(t *testing.T) {
    t.Parallel()

    if testing.Short() {
        t.Skip("long running test")
    }

    // 1. setup
    db := setup(t, "storeCreate")

    t.Cleanup(func() {
        if err := db.Close(); err != nil {
            t.Fatalf("unable to close database connection: %v", err)
        }
    })

    // 2. test
    testCases := []struct {
        name        string
        actual      *userstore.User
        expected    *userstore.User
        expectedErr error
    }{
        {
            name:        "user is nil",
            expectedErr: userstore.ErrDataMissing,
        },
        {
            name:        "user has password only",
            actual:      &userstore.User{
        Password: "JC7jCCcBRAb1xrDx2pze1N8LST7TkGqEbzjO8Z5ahwWM",
},
            expectedErr: userstore.ErrDataMissing,
        },
        {
            name:        "user has externalid only",
            actual:      &userstore.User{
        ExternalID: "9Lgq64Fx5lBnW8U3LvDghqMrsvfcTSrQ72b3QL",
},
            expectedErr: userstore.ErrDataMissing,
        },
        {
            name:        "user has id",
            actual:      &userstore.User{
        ID: testingutils.UUIDParse(t, "d2272a69-8efc-4ba6-b89f-5e67e5daa64d"),
        EMail: "ethanthomas576@example.org",
        FirstName: "Isabella",
        LastName: "Davis",
        Password: "OgkikSK36aQRR0fkRDw9Qg8tJE6exXGphpLn0AZEWmt",
        ExternalID: "gSB78qRRGiOJMShStMQGh",
        Type: "NrAEicMCGVbz4IL3d",
},
            expectedErr: userstore.ErrIDIsSet,
        },
        {
            name:        "user has all fields set",
            actual:      &userstore.User{
        EMail: "davidmiller080@example.com",
        FirstName: "Emily",
        LastName: "Garcia",
        Password: "g4qFw25CfMHzXD9vdpkyLVnsueEh",
        ExternalID: "0svOKSNzystuZcjkd0FGP7r6U1",
        Type: "FK4PJ0fbjv5GUtNOAlZCwwgBVXHlfJpHmnt56JplUdOoMhFys",
},
            expected: &userstore.User{
        EMail: "davidmiller080@example.com",
        FirstName: "Emily",
        LastName: "Garcia",
        Password: "g4qFw25CfMHzXD9vdpkyLVnsueEh",
        ExternalID: "0svOKSNzystuZcjkd0FGP7r6U1",
        Type: "FK4PJ0fbjv5GUtNOAlZCwwgBVXHlfJpHmnt56JplUdOoMhFys",
},
        },
        {
            name:        "user has only mandatory fields set",
            actual:      &userstore.User{
        EMail: "ellawilliams556@example.com",
        FirstName: "Noah",
        LastName: "Thomas",
        Type: "aNdUbBWkAN0lqw69X9HWRh8ScUDk",
},
            expected: &userstore.User{
        EMail: "ellawilliams556@example.com",
        FirstName: "Noah",
        LastName: "Thomas",
        Type: "aNdUbBWkAN0lqw69X9HWRh8ScUDk",
},
        },
    }

    for _, testCase := range testCases {
        testCase := testCase
        t.Run(testCase.name, func(t *testing.T) {
            t.Parallel()

            err := testCase.actual.Create(context.Background(), db)
            testingutils.ErrorsCheck(t, testCase.expectedErr, err)

            if testCase.expectedErr == nil {
                testCase.expected.ID = testCase.actual.ID
                assertUser(t, testCase.expected, testCase.actual)
            }
        })
    }
}

func TestUser_Read(t *testing.T) {
    t.Parallel()

    if testing.Short() {
        t.Skip("long running test")
    }

    // 1. setup
    db := setup(t, "storeRead")

    t.Cleanup(func() {
        if err := db.Close(); err != nil {
            t.Fatalf("unable to close database connection: %v", err)
        }
    })

    // 2. test
    testCases := []struct {
        name        string
        prepare     *userstore.User
        expected    *userstore.User
        expectedErr error
    }{
        {
            name:        "user is nil",
            expectedErr: userstore.ErrIDMissing,
        },
        {
            name:        "ID not set",
            expectedErr: userstore.ErrIDMissing,
        },
        {
            name:        "success",
            prepare:      &userstore.User{
        EMail: "josephtaylor432@test.net",
        FirstName: "Chloe",
        LastName: "Garcia",
        Password: "2Jbn1oLkDB5Tum96Is7ULZ8aU0dxVXwFNZjChTmzuUtF",
        ExternalID: "tyl1pNQZU1IGE5O7Co83aerR0IVf4od5t8JPJaUabj",
        Type: "ZblvkyIfkAQV2MNymLQQPCt71irSkp8wJSJb2",
},
            expected: &userstore.User{
        EMail: "josephtaylor432@test.net",
        FirstName: "Chloe",
        LastName: "Garcia",
        Password: "2Jbn1oLkDB5Tum96Is7ULZ8aU0dxVXwFNZjChTmzuUtF",
        ExternalID: "tyl1pNQZU1IGE5O7Co83aerR0IVf4od5t8JPJaUabj",
        Type: "ZblvkyIfkAQV2MNymLQQPCt71irSkp8wJSJb2",
},
        },
        {
            name:        "not existing",
            prepare:      &userstore.User{
        ID: testingutils.UUIDParse(t, "c90d5cc8-ae33-48f9-8507-4c4fe25583ca"),
},
            expectedErr: sql.ErrNoRows,
        },
    }

    for _, testCase := range testCases {
        testCase := testCase
        t.Run(testCase.name, func(t *testing.T) {
            t.Parallel()

            var id uuid.UUID
            if testCase.prepare != nil {
                if testCase.prepare.IsValid() {
                    err := testCase.prepare.Create(context.Background(), db)
                    if err != nil {
                        t.Errorf("preparation failed: %v", err)

                        return
                    }
                }
                id = testCase.prepare.ID
            }

            actual := &userstore.User{ID: id}
            err := actual.Read(context.Background(), db)
            testingutils.ErrorsCheck(t, testCase.expectedErr, err)

            if testCase.expectedErr == nil {
                testCase.expected.ID = actual.ID
                assertUser(t, testCase.expected, actual)
            }
        })
    }
}

func assertUser(t *testing.T, expected, actual *userstore.User) {
    t.Helper()

    if expected == nil && actual == nil {
        return
    }

    if expected != nil && actual == nil || expected == nil && actual != nil {
        t.Errorf("expected '%v' but got '%v'", expected, actual)

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
    

    if actual.CreatedAt.IsZero() {
        t.Error("created at should be greater than the zero date")
    }

    if actual.ModifiedAt.IsZero() {
        t.Error("modified at should be greater than the zero date")
    }
}

func TestUser_Update(t *testing.T) {
    t.Parallel()

    if testing.Short() {
        t.Skip("long running test")
    }

    // 1. setup
    db := setup(t, "storeUpdate")

    t.Cleanup(func() {
        if err := db.Close(); err != nil {
            t.Fatalf("unable to close database connection: %v", err)
        }
    })

    // 2. test
    testCases := []struct {
        name        string
        prepare     *userstore.User
        actual      *userstore.User
        expected    *userstore.User
        expectedErr error
    }{
        {
            name:        "user is nil",
            expectedErr: userstore.ErrDataMissing,
        },
        {
            name:        "user has password only",
            actual:      &userstore.User{
        Password: "fFNOftoPGq7WDsfX5kJDJX2r2hRXT7G",
},
                expectedErr: userstore.ErrDataMissing,
        },
        {
            name:        "user has externalid only",
            actual:      &userstore.User{
        ExternalID: "RGGXvnF3VXSR3atRfS4dEOy0tvTm199K5gIzxZVZ6d1rt",
},
                expectedErr: userstore.ErrDataMissing,
        },
        {
            name:        "user has no id",
            actual:      &userstore.User{
        EMail: "liammartinez317@example.com",
        FirstName: "Liam",
        LastName: "Anderson",
        Password: "9NnVfXYNxuhGXx5Gadl7xR6ZsAT",
        ExternalID: "Q9MmMFalI5sWkVXkB631WqGfugI3dtTn1THYIcsOGtRx",
        Type: "2gglbad7WjGA3gWjqMIWqNKVJFpDS",
},
                expectedErr: userstore.ErrIDMissing,
        },
        {
            name:        "not existing",
            actual:      &userstore.User{
        ID: testingutils.UUIDParse(t, "b96e0183-f13a-480e-a9b7-f00b08d95db9"),
        EMail: "michaelbrown612@test.net",
        FirstName: "Emma",
        LastName: "Jones",
        Type: "NtCxocHkMK2MY7I7sCIJIr5mlULdywmNXvTBX1B3l1il8fchw",
},
                expectedErr: sql.ErrNoRows,
        },
        {
            name:        "user has all fields set",
            actual:      &userstore.User{
        EMail: "chloemartin770@example.com",
        FirstName: "Aiden",
        LastName: "Anderson",
        Password: "25PHCATTNiZTjH",
        ExternalID: "bb3EEy7SCMRK1mcaL649tUcqSSpN",
        Type: "0uk58SU8yeIRsKYwqhbkjHotvmng",
},
                prepare: &userstore.User{
        EMail: "williammartinez142@test.org",
        FirstName: "Ethan",
        LastName: "Taylor",
        Password: "1GWfcATx7ywLXawY2JrM73p",
        ExternalID: "6gZPrS4SJLlKcD5y",
        Type: "ULsFqTdIGOBx4qkFyxuZUMRAc1r8gGl88C6",
},
                expected: &userstore.User{
        EMail: "chloemartin770@example.com",
        FirstName: "Aiden",
        LastName: "Anderson",
        Password: "25PHCATTNiZTjH",
        ExternalID: "bb3EEy7SCMRK1mcaL649tUcqSSpN",
        Type: "0uk58SU8yeIRsKYwqhbkjHotvmng",
},
        },
        {
            name:        "user has only mandatory fields set",
            actual:      &userstore.User{
        EMail: "williamwilson403@test.com",
        FirstName: "Addison",
        LastName: "White",
        Type: "dlAHJsgYA3bpC793VWCJDCdnYBSeZe3mNLz6KADzBWNSejYr",
},
                prepare: &userstore.User{
        EMail: "sophiaharris782@test.net",
        FirstName: "Andrew",
        LastName: "Smith",
        Type: "WaxmDjywFlnXw8eV771JqBqmKD6P7zghkE0kn7vOj8rg7sKas",
},
                expected: &userstore.User{
        EMail: "williamwilson403@test.com",
        FirstName: "Addison",
        LastName: "White",
        Type: "dlAHJsgYA3bpC793VWCJDCdnYBSeZe3mNLz6KADzBWNSejYr",
},
        },
    }

    for _, testCase := range testCases {
        testCase := testCase
        t.Run(testCase.name, func(t *testing.T) {
            t.Parallel()

            if testCase.prepare != nil {
                _ = testCase.prepare.Create(context.Background(), db)
                time.Sleep(1 * time.Second)
                testCase.actual.ID = testCase.prepare.ID
            }

            err := testCase.actual.Update(context.Background(), db)
            testingutils.ErrorsCheck(t, testCase.expectedErr, err)

            if testCase.expectedErr == nil {
                testCase.expected.ID = testCase.actual.ID
                assertUser(t, testCase.expected, testCase.actual)
            }

            if testCase.prepare != nil && testCase.actual != nil {
                if testCase.prepare.CreatedAt != testCase.actual.CreatedAt {
                    t.Errorf(
                        "expected created at '%s' but got '%s'",
                        testCase.prepare.CreatedAt.String(),
                        testCase.actual.CreatedAt.String(),
                    )
                }

                if testCase.prepare.ModifiedAt.After(testCase.actual.ModifiedAt) {
                    t.Errorf(
                        "expected modified at '%s' to be before but got '%s'",
                        testCase.prepare.ModifiedAt.String(),
                        testCase.actual.ModifiedAt.String(),
                    )
                }
            }
        })
    }
}

func TestUser_Delete(t *testing.T) {
    t.Parallel()

    if testing.Short() {
        t.Skip("long running test")
    }

    // 1. setup
    db := setup(t, "storeDelete")

    t.Cleanup(func() {
        if err := db.Close(); err != nil {
            t.Fatalf("unable to close database connection: %v", err)
        }
    })

    // 2. test
    testCases := []struct {
        name        string
        prepare     *userstore.User
        expectedErr error
    }{
        {
            name:        "user is nil",
            expectedErr: userstore.ErrIDMissing,
        },
        {
            name:        "user has no ID",
            expectedErr: userstore.ErrIDMissing,
        },
        {
            name: "success",
            prepare: &userstore.User{
        EMail: "danieltaylor778@test.org",
        FirstName: "Isabella",
        LastName: "Robinson",
        Password: "nChz5flr6knF7E4bLjhDWXd80L9zB1JNZpyjpvj",
        ExternalID: "Y9AHCfdj9SjUH2o8Zg7fdd66eJdkcaRW2",
        Type: "NvezPH7bE6TJWMNlj0rzjF",
},
        },
        {
            name: "not existing",
            prepare: &userstore.User{
        ID: testingutils.UUIDParse(t, "4f4f7d92-14e9-4e0b-a2a5-94021efeffd4"),
},
        },
    }

    for _, testCase := range testCases {
        testCase := testCase
        t.Run(testCase.name, func(t *testing.T) {
            t.Parallel()

            var id uuid.UUID
            if testCase.prepare != nil {
                if testCase.prepare.IsValid() {
                    err := testCase.prepare.Create(context.Background(), db)
                    if err != nil {
                        t.Errorf("preparation failed: %v", err)

                        return
                    }
                }
                id = testCase.prepare.ID
            }

            actual := &userstore.User{ID: id}
            err := actual.Delete(context.Background(), db)
            testingutils.ErrorsCheck(t, testCase.expectedErr, err)

            if !uuidutils.IsEmpty(id) {
                err := actual.Read(context.Background(), db)
                if !errors.Is(err, sql.ErrNoRows) {
                    t.Errorf("expected error '%v' after deletion but got '%v'", sql.ErrNoRows, err)
                }
            }
        })
    }
}

func TestUser_IsValid(t *testing.T) {
    t.Parallel()

    testCases := []struct {
        name     string
        actual   *userstore.User
        expected bool
    }{
        {
            name:     "user is nil",
            expected: false,
        },
        {
            name:     "user has id only",
            actual:   &userstore.User{
        ID: testingutils.UUIDParse(t, "50623127-adda-4a30-8707-f3e8bf43529c"),
},
            expected: false,
        },
        {
            name:     "user has email only",
            actual:   &userstore.User{
        EMail: "oliviajones045@example.org",
},
            expected: false,
        },
        {
            name:     "user has firstname only",
            actual:   &userstore.User{
        FirstName: "Elijah",
},
            expected: false,
        },
        {
            name:     "user has lastname only",
            actual:   &userstore.User{
        LastName: "Martin",
},
            expected: false,
        },
        {
            name:     "user has password only",
            actual:   &userstore.User{
        Password: "yJqaO7kK0dF3LMqzrqaLxhCDHrBeU9VpAIBGkR2GSN8qO",
},
            expected: false,
        },
        {
            name:     "user has externalid only",
            actual:   &userstore.User{
        ExternalID: "4bniDHJPniINfGN3Oe6hu",
},
            expected: false,
        },
        {
            name:     "user has type only",
            actual:   &userstore.User{
        Type: "9aZLvn60Ju60jyLtDM5x6MVrN",
},
            expected: false,
        },
        {
            name:     "mandatory fields only",
            actual:   &userstore.User{
        EMail: "masonmartin766@test.org",
        FirstName: "Jayden",
        LastName: "Wilson",
        Type: "BYeVpplvpSfRoQnR",
},
            expected: true,
        },
        {
            name:     "mandatory fields with id",
            actual:   &userstore.User{
        ID: testingutils.UUIDParse(t, "24115a47-805b-4cb3-841e-2a37544aeea9"),
        EMail: "josephdavis771@test.net",
        FirstName: "William",
        LastName: "Brown",
        Type: "Oekw1WROisy63iSKSl8jKspgYGTqmr5HS",
},
            expected: true,
        },
        {
            name:     "all fields",
            actual:   &userstore.User{
        ID: testingutils.UUIDParse(t, "7afcced6-e823-4243-8ebe-f9926c47cbee"),
        EMail: "emmadavis786@example.net",
        FirstName: "William",
        LastName: "Harris",
        Password: "Y0NBc6c9O",
        ExternalID: "V0QHVXtv2ZSo7WIo0KozTQHbfTXUvYyVT6",
        Type: "n0dz4CiFAdAvZw6ZxG2zfYkq2991oPU43CcZG",
},
            expected: true,
        },
        {
            name:     "all fields without id",
            actual:   &userstore.User{
        EMail: "jaydenharris072@test.org",
        FirstName: "Ava",
        LastName: "Anderson",
        Password: "eKkJCvnEbBioH",
        ExternalID: "AigBX9azkHVvRfbC6JRLAc8x",
        Type: "YOsXRyh08VWB9A",
},
            expected: true,
        },
    }

    for _, testCase := range testCases {
        testCase := testCase
        t.Run(testCase.name, func(t *testing.T) {
            t.Parallel()

            res := testCase.actual.IsValid()
            if testCase.expected != res {
                t.Errorf("expected %t but got %t", testCase.expected, res)
            }
        })
    }
}
