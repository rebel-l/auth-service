package facebook

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/rebel-l/auth-service/facebookapi"
	"github.com/rebel-l/smis"
)

// nolint: gocritic
//var (
//	ErrOnInit = errors.New("facebook init failed")
//)

type facebook struct {
	api facebookapi.API
	db  *sqlx.DB
	svc *smis.Service
}

// Init initialises the facebook endpoints.
func Init(svc *smis.Service, db *sqlx.DB, client facebookapi.Client) error {
	endpoint := &facebook{
		api: facebookapi.New(client),
		db:  db,
		svc: svc,
	}

	_, err := svc.RegisterEndpointToPublicChain(pathLogin, http.MethodPut, endpoint.loginPutHandler)
	if err != nil {
		return fmt.Errorf("failed to init handler with path '%s:%s': %w", http.MethodPut, pathLogin, err)
	}

	return nil
}
