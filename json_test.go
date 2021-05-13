package httputils

import (
	"io"
	"net/http"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReadJSON(test *testing.T) {
	type args struct {
		reader io.Reader
		data   interface{}
	}
	type testData struct {
		FieldOne int
		FieldTwo string
	}

	for _, data := range []struct {
		name     string
		args     args
		wantData interface{}
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "success with the null",
			args: args{
				reader: func() io.Reader {
					reader := new(MockReader)
					reader.
						On("Read", mock.AnythingOfType("[]uint8")).
						Return(func(buffer []byte) int { return copy(buffer, "null") }, io.EOF)

					return reader
				}(),
				data: new(testData),
			},
			wantData: new(testData),
			wantErr:  assert.NoError,
		},
		{
			name: "success with an empty data",
			args: args{
				reader: func() io.Reader {
					reader := new(MockReader)
					reader.
						On("Read", mock.AnythingOfType("[]uint8")).
						Return(func(buffer []byte) int { return copy(buffer, "{}") }, io.EOF)

					return reader
				}(),
				data: new(testData),
			},
			wantData: new(testData),
			wantErr:  assert.NoError,
		},
		{
			name: "success with a non-empty data",
			args: args{
				reader: func() io.Reader {
					reader := new(MockReader)
					reader.
						On("Read", mock.AnythingOfType("[]uint8")).
						Return(
							func(buffer []byte) int {
								return copy(buffer, `{"FieldOne": 23, "FieldTwo": "test"}`)
							},
							io.EOF,
						)

					return reader
				}(),
				data: new(testData),
			},
			wantData: &testData{FieldOne: 23, FieldTwo: "test"},
			wantErr:  assert.NoError,
		},
		{
			name: "error with a nil pointer",
			args: args{
				reader: new(MockReader),
				data:   (*testData)(nil),
			},
			wantData: (*testData)(nil),
			wantErr:  assert.Error,
		},
		{
			name: "error with a data type",
			args: args{
				reader: new(MockReader),
				data:   "",
			},
			wantData: "",
			wantErr:  assert.Error,
		},
		{
			name: "error with the data reading",
			args: args{
				reader: func() io.Reader {
					reader := new(MockReader)
					reader.
						On("Read", mock.AnythingOfType("[]uint8")).
						Return(0, iotest.ErrTimeout)

					return reader
				}(),
				data: new(testData),
			},
			wantData: new(testData),
			wantErr:  assert.Error,
		},
		{
			name: "error with the data unmarshalling",
			args: args{
				reader: func() io.Reader {
					reader := new(MockReader)
					reader.
						On("Read", mock.AnythingOfType("[]uint8")).
						Return(
							func(buffer []byte) int { return copy(buffer, "incorrect") },
							io.EOF,
						)

					return reader
				}(),
				data: new(testData),
			},
			wantData: new(testData),
			wantErr:  assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotErr := ReadJSON(data.args.reader, data.args.data)

			mock.AssertExpectationsForObjects(test, data.args.reader)
			assert.Equal(test, data.wantData, data.args.data)
			data.wantErr(test, gotErr)
		})
	}
}

func TestWriteJSON(test *testing.T) {
	type args struct {
		writer     http.ResponseWriter
		statusCode int
		data       interface{}
	}
	type testData struct {
		FieldOne int
		FieldTwo string
	}

	for _, data := range []struct {
		name       string
		args       args
		wantHeader http.Header
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				writer: func() http.ResponseWriter {
					const bytes = `{"FieldOne":23,"FieldTwo":"test"}`

					writer := new(MockResponseWriter)
					writer.On("Header").Return(http.Header{})
					writer.On("WriteHeader", http.StatusOK).Return()
					writer.On("Write", []byte(bytes)).Return(len(bytes), nil)

					return writer
				}(),
				statusCode: http.StatusOK,
				data:       testData{FieldOne: 23, FieldTwo: "test"},
			},
			wantHeader: http.Header{"Content-Type": {"application/json"}},
			wantErr:    assert.NoError,
		},
		{
			name: "error with the marshalling",
			args: args{
				writer:     new(MockResponseWriter),
				statusCode: http.StatusOK,
				data:       func() {},
			},
			wantHeader: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error with the writing",
			args: args{
				writer: func() http.ResponseWriter {
					const bytes = `{"FieldOne":23,"FieldTwo":"test"}`

					writer := new(MockResponseWriter)
					writer.On("Header").Return(http.Header{})
					writer.On("WriteHeader", http.StatusOK).Return()
					writer.On("Write", []byte(bytes)).Return(0, iotest.ErrTimeout)

					return writer
				}(),
				statusCode: http.StatusOK,
				data:       testData{FieldOne: 23, FieldTwo: "test"},
			},
			wantHeader: http.Header{"Content-Type": {"application/json"}},
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotErr := WriteJSON(data.args.writer, data.args.statusCode, data.args.data)

			mock.AssertExpectationsForObjects(test, data.args.writer)
			if data.wantHeader != nil {
				assert.Equal(test, data.wantHeader, data.args.writer.Header())
			}
			data.wantErr(test, gotErr)
		})
	}
}
