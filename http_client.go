package httputils

import (
	"net/http"
)

// HTTPClient ...
//
// It represents the simplified interface of the http.Client structure
// (the primary method only). It is useful for mocking the latter.
//
type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}
