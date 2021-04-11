package auth

import (
	"errors"
	"net/http"

	"github.com/rebel-l/smis"
)

var (
	ErrNoTokenValidator = errors.New("token validate should never be nil")
	errInvalidToken     = smis.Error{
		Code:       "AUTH002",
		StatusCode: http.StatusUnauthorized,
		External:   "not authorization for this operation",
		Internal:   "authorization failed as token is invalid",
	}
)

type Middleware struct {
	svc       *smis.Service
	validator TokenValidator
}

func NewMiddleware(svc *smis.Service, validator TokenValidator) (*Middleware, error) {
	if validator == nil {
		return nil, ErrNoTokenValidator
	}

	return &Middleware{svc: svc, validator: validator}, nil
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log := m.svc.NewLogForRequestID(request.Context())
		resp := smis.Response{Log: log}

		defer func() {
			_ = request.Body.Close()
		}()

		if !m.validator.IsAccessTokenValid(request.Header) {
			resp.WriteJSONError(writer, errInvalidToken)

			return
		}

		// handle next
		next.ServeHTTP(writer, request)
	})
}
