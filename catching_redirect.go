package httputils

import (
	"net/http"
)

// CatchingRedirect ...
//
// It's a complete analog of the http.Redirect() function with capturing
// the last error from the Write() method of the http.ResponseWriter interface.
//
func CatchingRedirect(
	writer http.ResponseWriter,
	request *http.Request,
	url string,
	statusCode int,
) error {
	catchingWriter := NewCatchingResponseWriter(writer)
	http.Redirect(catchingWriter, request, url, statusCode)

	return catchingWriter.LastError()
}
