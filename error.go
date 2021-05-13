package httputils

import (
	"net/http"

	"github.com/go-log/log"
)

// LoggingError ...
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
