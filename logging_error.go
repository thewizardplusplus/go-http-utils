package httputils

import (
	"net/http"

	"github.com/go-log/log"
)

// LoggingError ...
//
// It's a complete analog of the http.Error() function with the additional
// logging of the error. This function also accepts an error object instead of
// an error string.
//
func LoggingError(
	logger log.Logger,
	writer http.ResponseWriter,
	err error,
	statusCode int,
) {
	errMessage := err.Error()
	logger.Log(errMessage)
	http.Error(writer, errMessage, statusCode)
}
