package errors

import (
	"fmt"
	"net/http"

	"github.com/micro/go-micro/errors"
)

const (
	StatusIgnorableError = -1
)

// IgnoreError generates a -1 error.
func IgnorableError(id, format string, a ...interface{}) error {
	return &errors.Error{
		Id:     id,
		Code:   StatusIgnorableError,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(200),
	}
}