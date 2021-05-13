package httputils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/gorilla/mux"
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
		{
			name: "success with an integer",
			args: args{
				request: mux.SetURLVars(
					httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
					map[string]string{"test": "23"},
				),
				name: "test",
				data: pointer.ToInt(0),
			},
			wantData: pointer.ToInt(23),
			wantErr:  assert.NoError,
		},
		{
			name: "success with a string",
			args: args{
				request: mux.SetURLVars(
					httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
					map[string]string{"test": "data"},
				),
				name: "test",
				data: pointer.ToString(""),
			},
			wantData: pointer.ToString("data"),
			wantErr:  assert.NoError,
		},
		{
			name: "error with a nil pointer",
			args: args{
				request: mux.SetURLVars(
					httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
					map[string]string{"test": "23"},
				),
				name: "test",
				data: (*int)(nil),
			},
			wantData: (*int)(nil),
			wantErr:  assert.Error,
		},
		{
			name: "error with a data type",
			args: args{
				request: mux.SetURLVars(
					httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
					map[string]string{"test": "23"},
				),
				name: "test",
				data: "",
			},
			wantData: "",
			wantErr:  assert.Error,
		},
		{
			name: "error with the parameter existence",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
				name:    "test",
				data:    pointer.ToInt(0),
			},
			wantData: pointer.ToInt(0),
			wantErr:  assert.Error,
		},
		{
			name: "error with the data scanning",
			args: args{
				request: mux.SetURLVars(
					httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
					map[string]string{"test": "incorrect"},
				),
				name: "test",
				data: pointer.ToInt(0),
			},
			wantData: pointer.ToInt(0),
			wantErr:  assert.Error,
		},
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
