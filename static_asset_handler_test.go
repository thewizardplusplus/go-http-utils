package httputils

import (
	"net/http"
	"testing"

	"github.com/go-log/log"
	"github.com/stretchr/testify/mock"
)

func TestStaticAssetHandler(test *testing.T) {
	type args struct {
		fileSystem http.FileSystem
		logger     log.Logger
	}
	type handlerArgs struct {
		writer  http.ResponseWriter
		request *http.Request
	}

	for _, data := range []struct {
		name        string
		args        args
		handlerArgs handlerArgs
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			handler := StaticAssetHandler(data.args.fileSystem, data.args.logger)
			handler.ServeHTTP(data.handlerArgs.writer, data.handlerArgs.request)

			mock.AssertExpectationsForObjects(
				test,
				data.args.fileSystem,
				data.args.logger,
				data.handlerArgs.writer,
			)
		})
	}
}
