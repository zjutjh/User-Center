package account

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/model"
	"github.com/zjutjh/User-Center/common/errorsx"
)

func NormalizeStudentID(studentID string) string {
	return strings.ToUpper(strings.TrimSpace(studentID))
}

func NormalizeCardID(cardID string) string {
	return strings.ToUpper(strings.TrimSpace(cardID))
}

func HashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func VerifyPassword(hashedPassword, password string) bool {
	return hashedPassword == HashPassword(password)
}

func ValidatePasswordLength(password string) error {
	if len(password) < 6 || len(password) > 20 {
		return errorsx.ErrPasswordLength
	}
	return nil
}

func VerifyStudentIdentity(student *model.Student, cardID string) error {
	if student == nil {
		return errorsx.ErrUserNotExist
	}
	if !strings.EqualFold(student.IDCard, NormalizeCardID(cardID)) {
		return errorsx.ErrParameterInvalid
	}
	return nil
}

func BoundMarker(value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	return "BOUND"
}
