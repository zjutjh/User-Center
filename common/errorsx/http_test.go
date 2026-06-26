package errorsx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeHTTPErrorKeepsCodeError(t *testing.T) {
	err := NormalizeHTTPError(ErrWrongAccountOrPassword)

	require.ErrorIs(t, err, ErrWrongAccountOrPassword)
}

func TestNormalizeHTTPErrorMapsJSONAndMappingErrors(t *testing.T) {
	testCases := []error{
		&json.SyntaxError{},
		&json.UnmarshalTypeError{},
		&strconv.NumError{},
		io.ErrUnexpectedEOF,
		io.EOF,
		errors.New("type mismatch"),
		errors.New(`"age" is not set`),
		errors.New(`"profile" is not fully set`),
		errors.New(`field "name" mustn't be nil`),
		errors.New(`the value in map is not string or json.Number, but bool`),
		fmt.Errorf("fullName: `user.age`, error: `%w`", errors.New("type mismatch")),
	}

	for _, testErr := range testCases {
		err := NormalizeHTTPError(testErr)
		require.ErrorIs(t, err, ErrParameterInvalid, "input=%v", testErr)
	}
}

func TestNormalizeHTTPErrorLeavesUnknownErrors(t *testing.T) {
	original := errors.New("database unavailable")

	err := NormalizeHTTPError(original)

	require.ErrorIs(t, err, original)
	require.NotErrorIs(t, err, ErrParameterInvalid)
}

func TestHTTPStatusAlwaysOK(t *testing.T) {
	testCases := []error{
		nil,
		ErrParameterInvalid,
		ErrNotLoggedIn,
		ErrUnknown,
		errors.New("plain error"),
	}

	for _, testErr := range testCases {
		require.Equal(t, 200, HTTPStatus(testErr))
	}
}
