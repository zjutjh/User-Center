package account

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/model"
	"github.com/zjutjh/User-Center/common/errorsx"
)

func TestNormalizeStudentIDTrimsAndUppercases(t *testing.T) {
	require.Equal(t, "ABC123", NormalizeStudentID("  abc123  "))
}

func TestNormalizeCardIDTrimsAndUppercases(t *testing.T) {
	require.Equal(t, "ID123X", NormalizeCardID("  id123x  "))
}

func TestHashAndVerifyPassword(t *testing.T) {
	hashed := HashPassword("secret123")

	require.NotEmpty(t, hashed)
	require.NotEqual(t, "secret123", hashed)
	require.True(t, VerifyPassword(hashed, "secret123"))
	require.False(t, VerifyPassword(hashed, "wrong"))
}

func TestValidatePasswordLength(t *testing.T) {
	require.NoError(t, ValidatePasswordLength("123456"))
	require.NoError(t, ValidatePasswordLength("12345678901234567890"))
	require.ErrorIs(t, ValidatePasswordLength("12345"), errorsx.ErrPasswordLength)
	require.ErrorIs(t, ValidatePasswordLength("123456789012345678901"), errorsx.ErrPasswordLength)
}

func TestVerifyStudentIdentity(t *testing.T) {
	student := &model.Student{StudentID: "ABC123", IDCard: "ID123X"}

	require.NoError(t, VerifyStudentIdentity(student, "id123x"))
	require.ErrorIs(t, VerifyStudentIdentity(nil, "id123x"), errorsx.ErrUserNotExist)
	require.ErrorIs(t, VerifyStudentIdentity(student, "wrong"), errorsx.ErrParameterInvalid)
}

func TestBoundMarker(t *testing.T) {
	require.Equal(t, "", BoundMarker(""))
	require.Equal(t, "", BoundMarker("   "))
	require.Equal(t, "BOUND", BoundMarker("secret"))
}
