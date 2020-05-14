package httputils

import (
	"net/http"

	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/mux"
)

// SPAFallbackMiddleware ...
//
// A SPA often manages its routing itself, so all relevant requests must
// be handled by an index.html file. A backend, in turn, must be responsible
// for an API and distribution of static files.
//
// To separate these types of requests, the following heuristic is used.
// If the request uses the GET method and contains the text/html value
// in the Accept header, the index.html file is returned (more precisely, a path
// part of a request URL is replaced to "/"). All other requests are processed
// as usual. The fact is that any routing requests sent from a modern browser
// meet described requirements.
//
// This solution is based on the proxy implementation in the development server
// of the Create React App project. See:
// https://create-react-app.dev/docs/proxying-api-requests-in-development/
//
func SPAFallbackMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			if isStaticAssetRequest(request) {
				request.URL.Path = "/"
			}
			next.ServeHTTP(writer, request)
		})
	}
}

func isStaticAssetRequest(request *http.Request) bool {
	if request.Method != http.MethodGet {
		return false
	}

	for _, spec := range header.ParseAccept(request.Header, "Accept") {
		if spec.Value == "text/html" {
			return true
		}
	}

	return false
}
