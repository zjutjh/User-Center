package validatorx

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type requiredRequest struct {
	Name string `json:"name" validate:"required"`
}

func TestHTTPValidatorRejectsRequiredBlankString(t *testing.T) {
	err := HTTPValidator{}.Validate(httptest.NewRequest("POST", "/", nil), &requiredRequest{
		Name: "  ",
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "is not set")
}

func TestHTTPValidatorAcceptsRequiredString(t *testing.T) {
	err := HTTPValidator{}.Validate(httptest.NewRequest("POST", "/", nil), &requiredRequest{
		Name: "mango",
	})

	require.NoError(t, err)
}
