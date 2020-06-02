package httputils

import (
	"net/http"

	"github.com/go-log/log"
)

//go:generate mockery -name=Logger -inpkg -case=underscore -testonly

// Logger ...
//
// It's used only for mock generating.
type Logger interface {
	log.Logger
}

//go:generate mockery -name=Handler -inpkg -case=underscore -testonly

// Handler ...
//
// It's used only for mock generating.
type Handler interface {
	http.Handler
}

//go:generate mockery -name=ResponseWriter -inpkg -case=underscore -testonly

// ResponseWriter ...
//
// It's used only for mock generating.
type ResponseWriter interface {
	http.ResponseWriter
}

//go:generate mockery -name=File -inpkg -case=underscore -testonly

// File ...
//
// It's used only for mock generating.
type File interface {
	http.File
}

//go:generate mockery -name=FileSystem -inpkg -case=underscore -testonly

// FileSystem ...
//
// It's used only for mock generating.
type FileSystem interface {
	http.FileSystem
}
