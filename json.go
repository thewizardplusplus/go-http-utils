package httputils

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/pkg/errors"
)

// ReadJSON ...
//
// It reads bytes from the reader and then unmarshals them into the data.
// The data should be a non-nil pointer.
//
// If the reader is the io.ReadCloser interface, this function
// does not close it.
//
// If the data is not a non-nil pointer, this function will return an error,
// and it will happen before the reader is read.
//
func ReadJSON(reader io.Reader, data interface{}) error {
	dataReflection := reflect.ValueOf(data)
	if dataReflection.Kind() != reflect.Ptr || dataReflection.IsNil() {
		return errors.New("the data is incorrect: it should be a non-nil pointer")
	}

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.Wrap(err, "unable to read the data")
	}

	if err := json.Unmarshal(bytes, data); err != nil {
		return errors.Wrap(err, "unable to unmarshal the data")
	}

	return nil
}
