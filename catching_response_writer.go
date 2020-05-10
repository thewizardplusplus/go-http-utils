package httputils

import (
	"net/http"
)

// CatchingResponseWriter ...
//
// It wraps the http.ResponseWriter interface to save the last error
// that occurred while calling the http.ResponseWriter.Write() method.
//
// Errors of writing via the http.ResponseWriter interface are important
// to handle. See for details: https://stackoverflow.com/a/43976633
//
// Attention! The CatchingResponseWriter structure only supports methods
// that are declared directly in the http.ResponseWriter interface. Any optional
// interfaces (http.Flusher, http.Hijacker, etc.) aren't supported.
//
type CatchingResponseWriter struct {
	http.ResponseWriter

	lastError error
}

// NewCatchingResponseWriter ...
//
// It allocates and returns a new CatchingResponseWriter object wrapping
// the provided http.ResponseWriter interface.
//
func NewCatchingResponseWriter(
	writer http.ResponseWriter,
) *CatchingResponseWriter {
	return &CatchingResponseWriter{ResponseWriter: writer}
}

// LastError ...
//
// It returns the last saved error from the Write() method of the wrapped
// http.ResponseWriter interface.
//
func (writer CatchingResponseWriter) LastError() error {
	return writer.lastError
}

func (writer *CatchingResponseWriter) Write(p []byte) (n int, err error) {
	n, err = writer.ResponseWriter.Write(p)
	writer.lastError = err

	return n, err
}
