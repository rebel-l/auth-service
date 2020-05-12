package facebook

import (
	"net/http"
	"time"
)

const (
	clientDefaultTimeout = 10 * time.Second
)

// TODO: move to httputils.NewClient
func NewClient() *http.Client {
	return &http.Client{
		Timeout: clientDefaultTimeout,
	}
}
