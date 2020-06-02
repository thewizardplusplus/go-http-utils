package httputils

import (
	stdlog "log"
	"net/http"
	"os"
	"testing"

	"github.com/go-log/log"
	"github.com/go-log/log/print"
	"github.com/stretchr/testify/mock"
)

func ExampleStaticAssetHandler() {
	// use the standard logger for error handling
	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags)
	staticAssetHandler := StaticAssetHandler(
		http.Dir("/var/www/example.com"),
		// wrap the standard logger via the github.com/go-log/log package
		print.New(logger),
	)

	http.Handle("/", staticAssetHandler)
	stdlog.Fatal(http.ListenAndServe(":8080", nil))
}

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
