package httputils

import (
	"net/http"
	"os"

	"github.com/go-log/log"
)

//go:generate mockery --name=Logger --inpackage --case=underscore --testonly

// Logger ...
//
// It's used only for mock generating.
type Logger interface {
	log.Logger
}

//go:generate mockery --name=Handler --inpackage --case=underscore --testonly

// Handler ...
//
// It's used only for mock generating.
type Handler interface {
	http.Handler
}

//go:generate mockery --name=ResponseWriter --inpackage --case=underscore --testonly

// ResponseWriter ...
//
// It's used only for mock generating.
type ResponseWriter interface {
	http.ResponseWriter
}

//go:generate mockery --name=FileInfo --inpackage --case=underscore --testonly

// FileInfo ...
//
// It's used only for mock generating.
type FileInfo interface {
	os.FileInfo
}

//go:generate mockery --name=File --inpackage --case=underscore --testonly

// File ...
//
// It's used only for mock generating.
type File interface {
	http.File
}

//go:generate mockery --name=FileSystem --inpackage --case=underscore --testonly

// FileSystem ...
//
// It's used only for mock generating.
type FileSystem interface {
	http.FileSystem
}
