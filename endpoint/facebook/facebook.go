package facebook

import (
	"fmt"
	"net/http"

	"github.com/rebel-l/auth-service/facebookapi"

	"github.com/rebel-l/smis"
)

//var (
//	ErrOnInit = errors.New("facebook init failed")
//)

type facebook struct {
	svc *smis.Service
	api facebookapi.API
}

// Init initialises the facebook endpoints.
func Init(svc *smis.Service, client facebookapi.Client) error {
	endpoint := &facebook{
		svc: svc,
		api: facebookapi.New(client),
	}

	_, err := svc.RegisterEndpointToPublicChain(pathLogin, http.MethodPut, endpoint.loginPutHandler)
	if err != nil {
		return fmt.Errorf("failed to init handler with path '%s:%s': %w", http.MethodPut, pathLogin, err)
	}

	//_, err = svc.RegisterEndpointToRestictedChain(pathLogin, http.MethodPut, endpoint.loginPutHandler)
	//_, err = svc.RegisterEndpointToPublicChain(pathLogin, http.MethodOptions, endpoint.loginPutHandler)

	return err
}
