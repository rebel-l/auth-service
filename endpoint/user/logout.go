package user

import (
	"net/http"

	"github.com/rebel-l/smis"
)

const (
	pathLogout = "/logout"
)

var errLogout = smis.Error{
	Code:       "AUTH001",
	StatusCode: http.StatusInternalServerError,
	External:   "logout failed",
	Internal:   "logout failed",
}

func (u *user) logoutHandler(writer http.ResponseWriter, request *http.Request) {
	log := u.svc.NewLogForRequestID(request.Context())
	resp := smis.Response{Log: log}

	defer func() {
		_ = request.Body.Close()
	}()

	if err := u.tokenManger.DeleteTokens(request); err != nil {
		resp.WriteJSONError(writer, errLogout.WithDetails(err))
	}

	writer.WriteHeader(http.StatusNoContent)
}
