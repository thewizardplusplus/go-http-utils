# go-http-utils

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-http-utils?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-http-utils)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-http-utils)](https://goreportcard.com/report/github.com/thewizardplusplus/go-http-utils)
[![Build Status](https://travis-ci.org/thewizardplusplus/go-http-utils.svg?branch=master)](https://travis-ci.org/thewizardplusplus/go-http-utils)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-http-utils/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/go-http-utils)

HTTP utility functions.

## Installation

Prepare the directory:

```
$ mkdir --parents "$(go env GOPATH)/src/github.com/thewizardplusplus/"
$ cd "$(go env GOPATH)/src/github.com/thewizardplusplus/"
```

Clone this repository:

```
$ git clone https://github.com/thewizardplusplus/go-http-utils.git
$ cd go-http-utils
```

Install dependencies with the [dep](https://golang.github.io/dep/) tool:

```
$ dep ensure -vendor-only
```

## Examples

`httputils.CatchingResponseWriter`:

```go
package main

import (
	"log"
	"net/http"

	httputils "github.com/thewizardplusplus/go-http-utils"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		catchingWriter := httputils.NewCatchingResponseWriter(writer)

		// use the catchingWriter object as the usual http.ResponseWriter interface
		http.Redirect(
			catchingWriter,
			request,
			"http://example.com/",
			http.StatusMovedPermanently,
		)

		if err := catchingWriter.LastError(); err != nil {
			log.Printf("unable to write the HTTP response: %v", err)
		}
	})
}
```

`httputils.CatchingMiddleware()`:

```go
package main

import (
	"fmt"
	stdlog "log"
	"net/http"
	"os"

	"github.com/go-log/log/print"
	httputils "github.com/thewizardplusplus/go-http-utils"
)

func main() {
	// use the standard logger for error handling
	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags)
	catchingMiddleware := httputils.CatchingMiddleware(
		// wrap the standard logger via the github.com/go-log/log package
		print.New(logger),
	)

	var handler http.Handler
	handler = http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		// writing error will be handled by the catching middleware
		fmt.Fprintln(writer, "Hello, world!")
	})
	handler = catchingMiddleware(handler)

	http.Handle("/", handler)
	logger.Fatal(http.ListenAndServe(":8080", nil))
}
```

`httputils.SPAFallbackMiddleware()`:

```go
package main

import (
	"log"
	"net/http"

	httputils "github.com/thewizardplusplus/go-http-utils"
)

func main() {
	staticAssetHandler := http.FileServer(http.Dir("/var/www/example.com"))
	staticAssetHandler = httputils.SPAFallbackMiddleware()(staticAssetHandler)

	http.Handle("/", staticAssetHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

`httputils.RunServer()`:

```go
package main

import (
	"context"
	stdlog "log"
	"net/http"
	"os"

	"github.com/go-log/log/print"
	httputils "github.com/thewizardplusplus/go-http-utils"
)

func main() {
	server := &http.Server{Addr: ":8080"}
	// use the standard logger for error handling
	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags)
	if ok := httputils.RunServer(
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
```

## License

The MIT License (MIT)

Copyright &copy; 2020 thewizardplusplus
