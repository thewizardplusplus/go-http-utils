// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package httputils

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// MockHandler is an autogenerated mock type for the Handler type
type MockHandler struct {
	mock.Mock
}

// ServeHTTP provides a mock function with given fields: _a0, _a1
func (_m *MockHandler) ServeHTTP(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}
