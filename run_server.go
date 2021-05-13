package httputils

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/go-log/log"
)

//go:generate mockery --name=Server --inpackage --case=underscore --testonly

// Server ...
//
// It represents the interface of the http.Server structure used
// by the RunServer() function.
//
type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// RunServer ...
//
// It runs the provided Server interface via its ListenAndServe() method,
// then waits for specified signals. If there're no signals specified,
// the function will wait for any possible signals.
//
// After receiving the signal, the function will stop the server
// via its Shutdown() method. The provided context will be used for calling
// this method, and just for that.
//
// Attention! The function will return only after completing of ListenAndServe()
// and Shutdown() methods of the server.
//
// Errors that occurred when calling ListenAndServe() and Shutdown() methods
// of the server will be processed by the provided log.Logger interface. (Use
// the github.com/go-log/log package to wrap the standard logger.) The fact
// that an error occurred will be reflected in the boolean flag returned
// by the function.
//
func RunServer(
	shutdownCtx context.Context,
	server Server,
	logger log.Logger,
	interruptSignals ...os.Signal,
) (ok bool) {
	// it's used to wait for the shutdown goroutine to complete
	var waiter sync.WaitGroup
	waiter.Add(1)

	// it's used to complete the shutdown goroutine if this function ends earlier
	// (before receiving any signal); it'll happen when an error occurs
	// when calling the ListenAndServe() method of the provided server
	interruptCtx, interruptCtxCancel := context.WithCancel(context.Background())
	defer interruptCtxCancel()

	go func() {
		defer waiter.Done()

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, interruptSignals...)

		select {
		case <-interrupt:
		case <-interruptCtx.Done():
			return
		}

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			// error with closing listeners
			logger.Logf("unable to shutdown the HTTP server: %v", err)
		}

		// update the result
		ok = err == nil
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// error with starting or closing listeners
		logger.Logf("unable to run the HTTP server: %v", err)
		return false
	}

	waiter.Wait()
	return ok
}
