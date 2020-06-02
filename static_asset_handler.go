package httputils

import (
	"net/http"

	"github.com/go-log/log"
)

// StaticAssetHandler ...
func StaticAssetHandler(
	fileSystem http.FileSystem,
	logger log.Logger,
) http.Handler {
	handler := http.FileServer(fileSystem)
	handler = SPAFallbackMiddleware()(handler)
	handler = CatchingMiddleware(logger)(handler)

	return handler
}
