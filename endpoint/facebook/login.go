package facebook

import "net/http"

const (
	pathLogin = "/facebook/login"
)

func (f *facebook) loginPutHandler(writer http.ResponseWriter, request *http.Request) {
	log := f.svc.NewLogForRequestID(request.Context())

	// nolint:godox
	// TODO:

	_, err := writer.Write([]byte("pong"))
	if err != nil {
		log.Errorf("ping failed: %s", err)
	}
}
