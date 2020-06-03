package httputils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCatchingRedirect(test *testing.T) {
	type args struct {
		writer     http.ResponseWriter
		request    *http.Request
		url        string
		statusCode int
	}

	for _, data := range []struct {
		name       string
		args       args
		wantHeader http.Header
		wantErr    error
	}{
		{
			name: "success without the Content-Type header",
			args: args{
				writer: func() http.ResponseWriter {
					body := fmt.Sprintf(
						"<a href=\"http://example.com/two\">%s</a>.\n\n",
						http.StatusText(http.StatusMovedPermanently),
					)

					writer := new(MockResponseWriter)
					writer.On("Header").Return(http.Header{})
					writer.On("WriteHeader", http.StatusMovedPermanently).Return()
					writer.On("Write", []byte(body)).Return(len(body), nil)

					return writer
				}(),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/one",
					nil,
				),
				url:        "http://example.com/two",
				statusCode: http.StatusMovedPermanently,
			},
			wantHeader: http.Header{
				"Content-Type": {"text/html; charset=utf-8"},
				"Location":     {"http://example.com/two"},
			},
			wantErr: nil,
		},
		{
			name: "success with the Content-Type header",
			args: args{
				writer: func() http.ResponseWriter {
					writer := new(MockResponseWriter)
					writer.
						On("Header").
						Return(http.Header{"Content-Type": {"application/octet-stream"}})
					writer.On("WriteHeader", http.StatusMovedPermanently).Return()

					return writer
				}(),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/one",
					nil,
				),
				url:        "http://example.com/two",
				statusCode: http.StatusMovedPermanently,
			},
			wantHeader: http.Header{
				"Content-Type": {"application/octet-stream"},
				"Location":     {"http://example.com/two"},
			},
			wantErr: nil,
		},
		{
			name: "error",
			args: args{
				writer: func() http.ResponseWriter {
					body := fmt.Sprintf(
						"<a href=\"http://example.com/two\">%s</a>.\n\n",
						http.StatusText(http.StatusMovedPermanently),
					)

					writer := new(MockResponseWriter)
					writer.On("Header").Return(http.Header{})
					writer.On("WriteHeader", http.StatusMovedPermanently).Return()
					writer.On("Write", []byte(body)).Return(0, iotest.ErrTimeout)

					return writer
				}(),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/one",
					nil,
				),
				url:        "http://example.com/two",
				statusCode: http.StatusMovedPermanently,
			},
			wantHeader: http.Header{
				"Content-Type": {"text/html; charset=utf-8"},
				"Location":     {"http://example.com/two"},
			},
			wantErr: iotest.ErrTimeout,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotErr := CatchingRedirect(
				data.args.writer,
				data.args.request,
				data.args.url,
				data.args.statusCode,
			)

			mock.AssertExpectationsForObjects(test, data.args.writer)
			assert.Equal(test, data.wantHeader, data.args.writer.Header())
			assert.Equal(test, data.wantErr, gotErr)
		})
	}
}
