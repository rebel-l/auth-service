package facebook

import (
	"net/http"

	"github.com/rebel-l/smis"
)

const (
	pathLogin = "/facebook/login"
)

var (
	errRequest = smis.Error{
		Code:       "FBL001",
		StatusCode: http.StatusBadRequest,
		External:   "no token received or not parsable",
		Internal:   "facebook login failed, no token received or not parsable",
	}

	errLogin = smis.Error{
		Code:       "FBL002",
		StatusCode: http.StatusInternalServerError,
		External:   "login failed",
		Internal:   "facebook login failed",
	}
)

type loginPayload struct {
	AccessToken string `json:"AccessToken"`
}

func (f *facebook) loginPutHandler(writer http.ResponseWriter, request *http.Request) {
	log := f.svc.NewLogForRequestID(request.Context())
	resp := smis.Response{Log: log}

	defer func() {
		_ = request.Body.Close()
	}()

	payload := &loginPayload{}
	if err := smis.ParseJSONRequestBody(request, payload); err != nil {
		resp.WriteJSONError(writer, errRequest.WithDetails(err))

		return
	}

	log.Infof("access token: %#v", payload)

	user, err := f.api.Me(payload.AccessToken)
	if err != nil {
		resp.WriteJSONError(writer, errLogin.WithDetails(err))

		return
	}

	log.Infof("FB Response Body: %#v", user)

	resp.WriteJSON(writer, http.StatusOK, user)
}
