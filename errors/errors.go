package errors

import (
	"encoding/json"
)

type IgnorableError struct {
	Err       string `json:"error"`
	Ignorable bool   `json:"ignorable"`
}

func (e *IgnorableError) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// UnwrapIgnorableError tries to parse a JSON string into an error. If that
// fails, it will set the given string as the error detail.
func UnwrapIgnorableError(err string) (bool, string) {
	if err == "" {
		return false, err
	}

	igerr := new(IgnorableError)
	uerr := json.Unmarshal([]byte(err), igerr)
	if uerr != nil {
		return false, err
	}

	return igerr.Ignorable, igerr.Err
}

// IgnoreError generates a -1 error.
func WrapIgnorableError(err error) error {
	if err == nil {
		return nil
	}

	return &IgnorableError{
		Err:       err.Error(),
		Ignorable: true,
	}
}