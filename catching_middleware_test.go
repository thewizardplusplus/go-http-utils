package httputils

import (
	"fmt"
	stdlog "log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"testing/iotest"

	"github.com/go-log/log"
	"github.com/go-log/log/print"
	"github.com/stretchr/testify/mock"
)

func ExampleCatchingMiddleware() {
	// use the standard logger for error handling
	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags)
	catchingMiddleware := CatchingMiddleware(
		// wrap the standard logger via the github.com/go-log/log package
		print.New(logger),
	)

	var handler http.Handler // nolint: staticcheck
	handler = http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		// writing error will be handled by the catching middleware
		fmt.Fprintln(writer, "Hello, world!") // nolint: errcheck
	})
	handler = catchingMiddleware(handler)

	http.Handle("/", handler)
	logger.Fatal(http.ListenAndServe(":8080", nil))
}

func TestCatchingMiddleware(test *testing.T) {
	type args struct {
		logger log.Logger
	}
	type middlewareArgs struct {
		next http.Handler
	}
	type handlerArgs struct {
		writer  http.ResponseWriter
		request *http.Request
	}

	for _, data := range []struct {
		name           string
		args           args
		middlewareArgs middlewareArgs
		handlerArgs    handlerArgs
	}{
		{
			name: "success",
			args: args{
				logger: new(MockLogger),
			},
			middlewareArgs: middlewareArgs{
				next: func() http.Handler {
					handler := new(MockHandler)
					handler.
						On(
							"ServeHTTP",
							mock.MatchedBy(func(writer http.ResponseWriter) bool {
								writer.Write([]byte("test")) // nolint: errcheck
								return true
							}),
							httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
						).
						Return()

					return handler
				}(),
			},
			handlerArgs: handlerArgs{
				writer: func() http.ResponseWriter {
					writer := new(MockResponseWriter)
					writer.On("Write", []byte("test")).Return(4, nil)

					return writer
				}(),
				request: httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
			},
		},
		{
			name: "error",
			args: args{
				logger: func() log.Logger {
					logger := new(MockLogger)
					logger.
						On(
							"Logf",
							mock.MatchedBy(func(string) bool { return true }),
							iotest.ErrTimeout,
						).
						Return()

					return logger
				}(),
			},
			middlewareArgs: middlewareArgs{
				next: func() http.Handler {
					handler := new(MockHandler)
					handler.
						On(
							"ServeHTTP",
							mock.MatchedBy(func(writer http.ResponseWriter) bool {
								writer.Write([]byte("test")) // nolint: errcheck
								return true
							}),
							httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
						).
						Return()

					return handler
				}(),
			},
			handlerArgs: handlerArgs{
				writer: func() http.ResponseWriter {
					writer := new(MockResponseWriter)
					writer.On("Write", []byte("test")).Return(2, iotest.ErrTimeout)

					return writer
				}(),
				request: httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			middleware := CatchingMiddleware(data.args.logger)
			handler := middleware(data.middlewareArgs.next)
			handler.ServeHTTP(data.handlerArgs.writer, data.handlerArgs.request)

			mock.AssertExpectationsForObjects(
				test,
				data.args.logger,
				data.middlewareArgs.next,
				data.handlerArgs.writer,
			)
		})
	}
}
