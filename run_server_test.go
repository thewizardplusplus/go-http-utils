package httputils

import (
	"context"
	stdlog "log"
	"net/http"
	"os"
	"testing"
	"testing/iotest"
	"time"

	"github.com/go-log/log"
	"github.com/go-log/log/print"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func ExampleRunServer() {
	server := &http.Server{Addr: ":8080"}
	// use the standard logger for error handling
	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags)
	if ok := RunServer(
		context.Background(),
		server,
		// wrap the standard logger via the github.com/go-log/log package
		print.New(logger),
		os.Interrupt,
	); !ok {
		// the error is already logged, so just end the program with the error status
		os.Exit(1)
	}
}

func TestRunServer(test *testing.T) {
	type args struct {
		shutdownCtx      context.Context
		server           Server
		logger           log.Logger
		interruptSignals []os.Signal
	}

	for _, data := range []struct {
		name   string
		args   args
		action func(test *testing.T)
		wantOk assert.BoolAssertionFunc
	}{
		{
			name: "success",
			args: args{
				shutdownCtx: context.Background(),
				server: func() Server {
					server := new(MockServer)
					server.On("ListenAndServe").Return(http.ErrServerClosed)
					server.
						On(
							"Shutdown",
							mock.MatchedBy(func(context.Context) bool { return true }),
						).
						Return(nil)

					return server
				}(),
				logger:           new(MockLogger),
				interruptSignals: []os.Signal{os.Interrupt},
			},
			action: func(test *testing.T) {
				time.Sleep(time.Second)

				currentProcess, err := os.FindProcess(os.Getpid())
				require.NoError(test, err)

				err = currentProcess.Signal(os.Interrupt)
				require.NoError(test, err)
			},
			wantOk: assert.True,
		},
		{
			name: "error on the ListenAndServe() call",
			args: args{
				shutdownCtx: context.Background(),
				server: func() Server {
					server := new(MockServer)
					server.On("ListenAndServe").Return(iotest.ErrTimeout)

					return server
				}(),
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
				interruptSignals: []os.Signal{os.Interrupt},
			},
			action: func(test *testing.T) {},
			wantOk: assert.False,
		},
		{
			name: "error on the Shutdown() call",
			args: args{
				shutdownCtx: context.Background(),
				server: func() Server {
					server := new(MockServer)
					server.On("ListenAndServe").Return(http.ErrServerClosed)
					server.
						On(
							"Shutdown",
							mock.MatchedBy(func(context.Context) bool { return true }),
						).
						Return(iotest.ErrTimeout)

					return server
				}(),
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
				interruptSignals: []os.Signal{os.Interrupt},
			},
			action: func(test *testing.T) {
				time.Sleep(time.Second)

				currentProcess, err := os.FindProcess(os.Getpid())
				require.NoError(test, err)

				err = currentProcess.Signal(os.Interrupt)
				require.NoError(test, err)
			},
			wantOk: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			go data.action(test)

			gotOk := RunServer(
				data.args.shutdownCtx,
				data.args.server,
				data.args.logger,
				data.args.interruptSignals...,
			)

			mock.AssertExpectationsForObjects(test, data.args.server, data.args.logger)
			data.wantOk(test, gotOk)
		})
	}
}
