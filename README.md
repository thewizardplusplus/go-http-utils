# go-http-utils

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-http-utils?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-http-utils)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-http-utils)](https://goreportcard.com/report/github.com/thewizardplusplus/go-http-utils)
[![Build Status](https://travis-ci.org/thewizardplusplus/go-http-utils.svg?branch=master)](https://travis-ci.org/thewizardplusplus/go-http-utils)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-http-utils/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/go-http-utils)

The library that provides HTTP utility functions.

## Features

- simplified interface of the `http.Client` structure for mocking purposes;
- wrapper for the `http.ResponseWriter` interface for catching writing errors;
- middlewares:
  - middleware for catching writing errors;
  - middleware that fallback of requests to static assets to the index.html file (useful in a SPA);
- functions:
  - analogs:
    - analog of the `http.Redirect()` function with catching writing errors;
    - analog of the `http.FileServer()` function with applied above-mentioned middlewares;
    - analog of the `http.Error()` function with the additional improvements:
      - additional logging of the error;
      - accepting of an error object instead of an error string;
  - function to extract the parameter with the specified name from the path part of the request URL and then scan it into the data;
  - JSON:
    - function to read bytes from the reader and then unmarshal them into the data;
    - function to marshal the data and then write it in the writer:
      - additional setting of the corresponding content type;
      - additional setting of the specified status code;
  - function to start a server with support for graceful shutdown by a signal.

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

`httputils.StaticAssetHandler()`:

```go
package main

import (
	stdlog "log"
	"net/http"
	"os"

	"github.com/go-log/log/print"
	httputils "github.com/thewizardplusplus/go-http-utils"
)

func main() {
	// use the standard logger for error handling
	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags)
	staticAssetHandler := httputils.StaticAssetHandler(
		http.Dir("/var/www/example.com"),
		// wrap the standard logger via the github.com/go-log/log package
		print.New(logger),
	)

	http.Handle("/", staticAssetHandler)
	stdlog.Fatal(http.ListenAndServe(":8080", nil))
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

## Bibliography

1. [Check Errors When Calling `http.ResponseWriter.Write()`](https://stackoverflow.com/a/43976633)
2. [Mocking HTTP Requests in Golang](https://www.thegreatcodeadventure.com/mocking-http-requests-in-golang/)
3. [Proxying API Requests in Development](https://create-react-app.dev/docs/proxying-api-requests-in-development/)
4. https://golang.org/pkg/net/http/#Server.Shutdown

## License

The MIT License (MIT)

Copyright &copy; 2020-2021 thewizardplusplus
