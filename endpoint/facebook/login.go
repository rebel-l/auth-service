package facebook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/rebel-l/smis"
)

const (
	pathLogin = "/facebook/login"
)

type loginPayload struct {
	AccessToken string `json:"AccessToken"`
}

func (f *facebook) loginPutHandler(writer http.ResponseWriter, request *http.Request) {
	log := f.svc.NewLogForRequestID(request.Context())

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(err.Error()))
		log.Errorf("facebook login failed to read request body: %v", err)

		return
	}

	var payload loginPayload
	if err = json.Unmarshal(body, &payload); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(err.Error()))
		log.Errorf("facebook login failed to parse request body: %v", err)
		log.Errorf("facebook login failed to parse request body: %v", err)

		return
	}

	log.Infof("access token: %#v", payload)

	user, err := f.api.Me(payload.AccessToken)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
		log.Errorf("facebook login: %v", err)

		return
	}
	log.Infof("FB Response Body: %#v", user)

	resp := smis.Response{Log: log}
	resp.WriteJSON(writer, http.StatusOK, user)
}
