package httputils

import (
	"io"
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
