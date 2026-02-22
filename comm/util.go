package comm

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256Hash(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	pass := hex.EncodeToString(h.Sum(nil))
	return pass
}
