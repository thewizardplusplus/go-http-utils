package httputils

import (
	"net/http"

	"github.com/go-log/log"
	"github.com/gorilla/mux"
)

// CatchingMiddleware ...
//
// It's a middleware that captures the last error from the Write() method
// of the http.ResponseWriter interface and logs it via the provided log.Logger
// interface.
//
func CatchingMiddleware(logger log.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			catchingWriter := NewCatchingResponseWriter(writer)
			next.ServeHTTP(catchingWriter, request)

			if err := catchingWriter.LastError(); err != nil {
				logger.Logf("unable to write the HTTP response: %v", err)
			}
		})
	}
}
