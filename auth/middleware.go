package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/rebel-l/smis"
)

const ContextKeyUserID = "userID"

var (
	ErrNoTokenValidator = errors.New("token validate should never be nil")

	errInvalidToken = smis.Error{
		Code:       "AUTH002",
		StatusCode: http.StatusUnauthorized,
		External:   "not authorized for this operation",
		Internal:   "authorization failed as token is invalid",
	}

	errExpiredToken = smis.Error{
		Code:       "AUTH003",
		StatusCode: http.StatusUnauthorized,
		External:   "authorization has expired",
		Internal:   "authorization failed as token has expired",
	}
)

type Authenticator interface {
	GetUserID(header http.Header) (uuid.UUID, error)
}

type Middleware struct {
	svc  *smis.Service
	auth Authenticator
}

func NewMiddleware(svc *smis.Service, auth Authenticator) (*Middleware, error) {
	if auth == nil {
		return nil, ErrNoTokenValidator
	}

	return &Middleware{svc: svc, auth: auth}, nil
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log := m.svc.NewLogForRequestID(request.Context())
		resp := smis.Response{Log: log}

		defer func() {
			_ = request.Body.Close()
		}()

		uID, err := m.auth.GetUserID(request.Header)
		if err != nil {
			if errors.Is(err, ErrTokenExpired) {
				resp.WriteJSONError(writer, errExpiredToken.WithDetails(err))

				return
			}

			resp.WriteJSONError(writer, errInvalidToken.WithDetails(err))

			return
		}

		// add user ID to context
		ctx := context.WithValue(request.Context(), ContextKeyUserID, uID)
		request = request.WithContext(ctx)

		// handle next
		next.ServeHTTP(writer, request)
	})
}
