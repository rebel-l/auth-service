package usermodel_test

import (
    "bytes"
    "errors"
    "io"
    "testing"
    "time"
    "github.com/rebel-l/auth-service/user/usermodel"
    "github.com/rebel-l/go-utils/testingutils"
)

func TestUser_DecodeJSON(t *testing.T) {
    t.Parallel()

    createdAt, _ := time.Parse(time.RFC3339Nano, "2019-12-31T03:36:57.9167778+01:00")
    modifiedAt, _ := time.Parse(time.RFC3339Nano, "2020-01-01T15:44:57.9168378+01:00")

    testCases := []struct {
        name        string
        actual      *usermodel.User
        json        io.Reader
        expected    *usermodel.User
        expectedErr error
    }{
        {
            name: "model is nil",
        },
        {
            name:        "no JSON format",
            actual:      &usermodel.User{},
            json:        bytes.NewReader([]byte("no JSON")),
            expected:    &usermodel.User{},
            expectedErr: usermodel.ErrDecodeJSON,
        },
        {
            name:   "success",
            actual: &usermodel.User{},
            json: bytes.NewReader([]byte(`
                {
    "ID": "b58bc532-4326-476c-84b6-08b09a9a2334",
    "EMail": "ellathomas148@example.org",
    "FirstName": "William",
    "LastName": "Garcia",
    "Password": "HA6XIklcKS4z",
    "ExternalID": "SdOqZ7zkIa72eVJlb3GuR5hBN1YJMd9kQtLmGjZR5RFw3GgGc",
    "Type": "a8dmw",
    "CreatedAt": "2019-12-31T03:36:57.9167778+01:00",
    "ModifiedAt": "2020-01-01T15:44:57.9168378+01:00"
}
            `)),
            expected: &usermodel.User{
    ID: testingutils.UUIDParse(t, "b58bc532-4326-476c-84b6-08b09a9a2334"),
    EMail: "ellathomas148@example.org",
    FirstName: "William",
    LastName: "Garcia",
    Password: "HA6XIklcKS4z",
    ExternalID: "SdOqZ7zkIa72eVJlb3GuR5hBN1YJMd9kQtLmGjZR5RFw3GgGc",
    Type: "a8dmw",
    CreatedAt:  createdAt,
    ModifiedAt: modifiedAt,
},
        },
        {
            name:     "empty json",
            actual:   &usermodel.User{},
            json:     bytes.NewReader([]byte("{}")),
            expected: &usermodel.User{},
        },
    }

    for _, testCase := range testCases {
        testCase := testCase
        t.Run(testCase.name, func(t *testing.T) {
            t.Parallel()

            err := testCase.actual.DecodeJSON(testCase.json)
            if !errors.Is(err, testCase.expectedErr) {
                t.Errorf("expected error '%v' but got '%v'", testCase.expectedErr, err)

                return
            }

            assertUser(t, testCase.expected, testCase.actual)
        })
    }
}

func assertUser(t *testing.T, expected, actual *usermodel.User) {
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
    

    if !expected.CreatedAt.Equal(actual.CreatedAt) {
        t.Errorf("expected created at '%s' but got '%s'", expected.CreatedAt.String(), actual.CreatedAt.String())
    }

    if !expected.ModifiedAt.Equal(actual.ModifiedAt) {
        t.Errorf("expected modified at '%s' but got '%s'", expected.ModifiedAt.String(), actual.ModifiedAt.String())
    }
}
