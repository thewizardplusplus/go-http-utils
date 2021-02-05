package httputils

import (
	"net/http"
)

// HTTPClient ...
//
// It represents the simplified interface of the http.Client structure
// (the primary method only). It is useful for mocking the latter.
//
// See for details:
// https://www.thegreatcodeadventure.com/mocking-http-requests-in-golang/
//
type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}
