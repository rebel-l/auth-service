package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/rebel-l/auth-service/auth"
	"github.com/rebel-l/smis"
)

const ContextKeyUserID contextType = "userID"

var (
	ErrNoTokenValidator = errors.New("token validate should never be nil")

	ErrInvalidToken = smis.Error{
		Code:       "AUTH002",
		StatusCode: http.StatusUnauthorized,
		External:   "not authorized for this operation",
		Internal:   "authorization failed as token is invalid",
	}

	ErrExpiredToken = smis.Error{
		Code:       "AUTH003",
		StatusCode: http.StatusUnauthorized,
		External:   "authorization has expired",
		Internal:   "authorization failed as token has expired",
	}
)

type contextType string

type Authenticator interface {
	GetUserID(header http.Header) (uuid.UUID, error)
}

type Auth struct {
	svc  *smis.Service
	auth Authenticator
}

func NewAuth(svc *smis.Service, auth Authenticator) (*Auth, error) {
	if auth == nil {
		return nil, ErrNoTokenValidator
	}

	return &Auth{svc: svc, auth: auth}, nil
}

func (m *Auth) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log := m.svc.NewLogForRequestID(request.Context())
		resp := smis.Response{Log: log}

		defer func() {
			_ = request.Body.Close()
		}()

		uID, err := m.auth.GetUserID(request.Header)
		if err != nil {
			if errors.Is(err, auth.ErrTokenExpired) {
				resp.WriteJSONError(writer, ErrExpiredToken.WithDetails(err))

				return
			}

			resp.WriteJSONError(writer, ErrInvalidToken.WithDetails(err))

			return
		}

		// add user ID to context
		ctx := context.WithValue(request.Context(), ContextKeyUserID, uID)
		request = request.WithContext(ctx)

		// handle next
		next.ServeHTTP(writer, request)
	})
}
