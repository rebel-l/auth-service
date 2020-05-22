package facebookapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rebel-l/smis"
)

const (
	baseURL       = "https://graph.facebook.com/v7.0"
	fieldsDefault = "id,name,email,first_name,last_name"
)

var (
	ErrRequest = errors.New("request to facebook failed")
)

type Client interface {
	Get(url string) (resp *http.Response, err error)
}

type API struct {
	client Client
}

func New(client Client) API {
	return API{client: client}
}

func (a API) Me(accessToken string) (*User, error) {
	requestURI := fmt.Sprintf("%s/me?fields=%s&access_token=%s", baseURL, fieldsDefault, accessToken)

	resp, err := a.client.Get(requestURI)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequest, err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status code was :%d", ErrRequest, resp.StatusCode)
	}

	user := &User{}
	if err = smis.ParseJSONResponseBody(resp, user); err != nil {
		return nil, err
	}

	return user, nil
}
