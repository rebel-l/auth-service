package user

import (
	"fmt"
	"net/http"

	"github.com/rebel-l/auth-service/auth"
	"github.com/rebel-l/smis"
)

type user struct {
	svc         *smis.Service
	tokenManger auth.TokenManager
}

// Init initialises the user endpoints.
func Init(svc *smis.Service, tokenManager auth.TokenManager) error {
	endpoint := &user{
		svc:         svc,
		tokenManger: tokenManager,
	}

	_, err := svc.RegisterEndpointToRestictedChain(pathLogout, http.MethodDelete, endpoint.logoutHandler)
	if err != nil {
		return fmt.Errorf("failed to init handler with path '%s:%s': %w", http.MethodPut, pathLogout, err)
	}

	_, err = svc.RegisterEndpointToPublicChain(pathRefresh, http.MethodPost, endpoint.refreshHandler)
	if err != nil {
		return fmt.Errorf("failed to init handler with path '%s:%s': %w", http.MethodPost, pathRefresh, err)
	}

	return nil
}
