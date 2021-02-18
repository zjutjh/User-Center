package util

import (
	"github.com/asaskevich/govalidator"
)

func ValidateUsername(username string) bool {
	if username == "" || len(username) > 32 {
		return false
	}
	if !govalidator.IsPrintableASCII(username) {
		return false
	}
	return true
}

func ValidatePassword(password string) bool {
	// only accepts SHA-256 hashed password (length = 64)
	if len(password) != 64 {
		return false
	}
	if !govalidator.IsPrintableASCII(password) {
		return false
	}
	return true
}
