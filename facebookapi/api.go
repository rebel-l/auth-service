package facebookapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	baseURL       = "https://graph.facebook.com/v7.0"
	fieldsDefault = "id,name,email,first_name,last_name"
)

var (
	ErrRequest       = errors.New("request to facebook failed") // TODO: introduce error struct with code
	ErrParseResponse = errors.New("failed to parse response from facebook")
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

	// TODO: move the following lines to smis.ParseJSONBody()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseResponse, err)
	}

	user := &User{}
	if err = json.Unmarshal(body, user); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseResponse, err)
	}

	return user, nil
}
