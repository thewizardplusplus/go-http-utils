package httputils

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/pkg/errors"
)

// ReadJSON ...
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
