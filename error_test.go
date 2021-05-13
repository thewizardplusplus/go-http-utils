package httputils

import (
	"net/http"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoggingError(test *testing.T) {
	logMessage := iotest.ErrTimeout.Error()
	logger := new(MockLogger)
	logger.On("Log", logMessage).Return()

	body := logMessage + "\n"
	writer := new(MockResponseWriter)
	writer.On("Header").Return(http.Header{})
	writer.On("WriteHeader", http.StatusInternalServerError).Return()
	writer.On("Write", []byte(body)).Return(len(body), nil)

	LoggingError(
		logger,
		writer,
		iotest.ErrTimeout,
		http.StatusInternalServerError,
	)

	wantHeader := http.Header{
		"Content-Type":           {"text/plain; charset=utf-8"},
		"X-Content-Type-Options": {"nosniff"},
	}
	mock.AssertExpectationsForObjects(test, logger, writer)
	assert.Equal(test, wantHeader, writer.Header())
}
