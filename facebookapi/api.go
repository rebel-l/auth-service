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

// ErrRequest describes the error happened on requests to facebook.
var ErrRequest = errors.New("request to facebook failed")

// Client defines the interface interacting with the facebook API.
type Client interface {
	Get(url string) (resp *http.Response, err error)
}

// API is the struct communicating with the API of facebook.
type API struct {
	client Client
}

// New returns a new API struct.
func New(client Client) API {
	return API{client: client}
}

// Me is the request to the API of facebook returning the information of the current logged in user.
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
		return nil, fmt.Errorf("failed to parse response from facebook API /me: %w", err)
	}

	return user, nil
}
