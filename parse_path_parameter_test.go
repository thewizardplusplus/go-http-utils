package httputils

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePathParameter(test *testing.T) {
	type args struct {
		request *http.Request
		name    string
		data    interface{}
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
			gotErr := ParsePathParameter(
				data.args.request,
				data.args.name,
				data.args.data,
			)

			assert.Equal(test, data.wantData, data.args.data)
			data.wantErr(test, gotErr)
		})
	}
}
