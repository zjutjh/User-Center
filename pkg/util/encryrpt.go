package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func Encryrpt(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	pass := hex.EncodeToString(h.Sum(nil))
	return pass
}
