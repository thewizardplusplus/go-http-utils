package httputils

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// ParsePathParameter ...
//
// It extracts the parameter with the specified name from the path part
// of the request URL and then scans it into the data.
//
// The extracting does not actually work with the request URL directly. Instead,
// this function relies on the use of the github.com/gorilla/mux package
// for routing.
//
// The scanning works via the fmt.Sscan() function and has the corresponding
// restrictions.
//
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
