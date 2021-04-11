package facebook

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/rebel-l/smis"

	"github.com/rebel-l/auth-service/auth"
	"github.com/rebel-l/auth-service/facebookapi"
)

// nolint: gocritic
//var (
//	ErrOnInit = errors.New("facebook init failed")
//)

type facebook struct {
	api         facebookapi.API
	db          *sqlx.DB
	svc         *smis.Service
	tokenManger auth.TokenGenerator
}

// Init initialises the facebook endpoints.
func Init(svc *smis.Service, db *sqlx.DB, tokenManager auth.TokenGenerator, client facebookapi.Client) error {
	endpoint := &facebook{
		api:         facebookapi.New(client),
		db:          db,
		svc:         svc,
		tokenManger: tokenManager,
	}

	_, err := svc.RegisterEndpointToPublicChain(pathLogin, http.MethodPut, endpoint.loginPutHandler)
	if err != nil {
		return fmt.Errorf("failed to init handler with path '%s:%s': %w", http.MethodPut, pathLogin, err)
	}

	return nil
}
