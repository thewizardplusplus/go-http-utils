package httputils

import (
	"net/http"
)

// HTTPClient ...
type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}
