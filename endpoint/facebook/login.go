package facebook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	pathLogin = "/facebook/login"
)

type loginPayload struct {
	// nolint:godox
	AccessToken string `json:"accessToken"` // TODO: unify key in frontend/backend to CamelCase
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

		return
	}

	log.Infof("access token: %#v", payload)

	_, err = writer.Write([]byte("pong"))
	if err != nil {
		log.Errorf("ping failed: %s", err)
	}
}
