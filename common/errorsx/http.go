package errorsx

import (
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"
)

// NormalizeHTTPError maps transport-layer request parsing errors to business errors.
func NormalizeHTTPError(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := As(err); ok {
		return err
	}
	if isHTTPParameterError(err) {
		return ErrParameterInvalid
	}

	return err
}

func isHTTPParameterError(err error) bool {
	var syntaxErr *json.SyntaxError
	var unmarshalTypeErr *json.UnmarshalTypeError
	var numErr *strconv.NumError

	switch {
	case errors.As(err, &syntaxErr):
		return true
	case errors.As(err, &unmarshalTypeErr):
		return true
	case errors.As(err, &numErr):
		return true
	case errors.Is(err, io.EOF):
		return true
	case errors.Is(err, io.ErrUnexpectedEOF):
		return true
	}

	msg := err.Error()
	for _, marker := range []string{
		"type mismatch",
		"is not set",
		"is not fully set",
		"mustn't be nil",
		"value out of range",
		"not string",
		"not string or json.Number",
		"not defined in options",
		"expect string for field",
	} {
		if strings.Contains(msg, marker) {
			return true
		}
	}

	return false
}
