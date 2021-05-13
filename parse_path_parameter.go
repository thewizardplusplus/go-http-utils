package httputils

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// ParsePathParameter ...
func ParsePathParameter(
	request *http.Request,
	name string,
	data interface{},
) error {
	dataReflection := reflect.ValueOf(data)
	if dataReflection.Kind() != reflect.Ptr || dataReflection.IsNil() {
		return errors.New("the data is incorrect: it should be a non-nil pointer")
	}

	value, ok := mux.Vars(request)[name]
	if !ok {
		return errors.New("the parameter is missing")
	}

	if _, err := fmt.Sscan(value, data); err != nil {
		return errors.Wrap(err, "unable to scan the data")
	}

	return nil
}
