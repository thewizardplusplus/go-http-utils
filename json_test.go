package httputils

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReadJSON(test *testing.T) {
	type args struct {
		reader io.Reader
		data   interface{}
	}

	for _, data := range []struct {
		name     string
		args     args
		wantData interface{}
		wantErr  assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			gotErr := ReadJSON(data.args.reader, data.args.data)

			mock.AssertExpectationsForObjects(test, data.args.reader)
			assert.Equal(test, data.wantData, data.args.data)
			data.wantErr(test, gotErr)
		})
	}
}
