package facebook

import (
	"net/http"

	"github.com/rebel-l/auth-service/auth"
	"github.com/rebel-l/auth-service/user/usermapper"
	"github.com/rebel-l/auth-service/user/usermodel"

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

type loginRequestBody struct {
	AccessToken string `json:"AccessToken"`
}

type loginResponseBody struct {
	FirstName string            `json:"FirstName"`
	Tokens    map[string]string `json:"Tokens"`
}

func newLoginResponseBody(m *usermodel.User, tokens map[string]*auth.Token) *loginResponseBody {
	resp := &loginResponseBody{
		Tokens: make(map[string]string),
	}

	if m != nil {
		resp.FirstName = m.FirstName
	}

	for k, v := range tokens {
		if v != nil {
			resp.Tokens[k] = v.JWT
		}
	}

	return resp
}

func (f *facebook) loginPutHandler(writer http.ResponseWriter, request *http.Request) {
	log := f.svc.NewLogForRequestID(request.Context())
	resp := smis.Response{Log: log}

	defer func() {
		_ = request.Body.Close()
	}()

	// get details from facebook API
	fbPayload := &loginRequestBody{}
	if err := smis.ParseJSONRequestBody(request, fbPayload); err != nil {
		resp.WriteJSONError(writer, errRequest.WithDetails(err))

		return
	}

	fbUser, err := f.api.Me(fbPayload.AccessToken)
	if err != nil {
		resp.WriteJSONError(writer, errLogin.WithDetails(err))

		return
	}

	// ensure user is in database
	model := usermodel.NewFromFacebook(fbUser)
	mapper := usermapper.New(f.db)

	model, err = mapper.SaveByEmail(request.Context(), model)
	if err != nil {
		resp.WriteJSONError(writer, errLogin.WithDetails(err))

		return
	}

	tokens, err := f.tokenManger.GenerateTokens(model)
	if err != nil {
		resp.WriteJSONError(writer, errLogin.WithDetails(err))
	}

	resp.WriteJSON(writer, http.StatusOK, newLoginResponseBody(model, tokens))
}
