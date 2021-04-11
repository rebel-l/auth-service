package user

import (
	"errors"
	"net/http"

	"github.com/rebel-l/auth-service/endpoint/middleware"

	"github.com/rebel-l/auth-service/auth"

	"github.com/rebel-l/smis"
)

const pathRefresh = "/refresh"

var errRequest = smis.Error{
	Code:       "AUTH004",
	StatusCode: http.StatusBadRequest,
	External:   "no token received or not parsable",
	Internal:   "refresh token failed, no token received or not parsable",
}

type refreshRequestBody struct {
	RefreshToken string `json:"RefreshToken"`
}

type refreshResponseBody struct {
	Tokens map[string]string `json:"Tokens"`
}

func (u *user) refreshHandler(writer http.ResponseWriter, request *http.Request) {
	log := u.svc.NewLogForRequestID(request.Context())
	resp := smis.Response{Log: log}

	defer func() {
		_ = request.Body.Close()
	}()

	requestBody := &refreshRequestBody{}
	if err := smis.ParseJSONRequestBody(request, requestBody); err != nil {
		resp.WriteJSONError(writer, errRequest.WithDetails(err))

		return
	}

	tokens, err := u.tokenManger.RefreshTokens(requestBody.RefreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrTokenExpired) {
			resp.WriteJSONError(writer, middleware.ErrExpiredToken.WithDetails(err))

			return
		}

		resp.WriteJSONError(writer, middleware.ErrInvalidToken.WithDetails(err))

		return
	}

	resp.WriteJSON(writer, http.StatusOK, &refreshResponseBody{Tokens: tokens})
}
