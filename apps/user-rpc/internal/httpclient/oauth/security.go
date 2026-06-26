package oauth

import (
	"crypto/rsa"
	"encoding/hex"
	"math/big"
	"strconv"
)

type publicKeyData struct {
	Modulus  string `json:"modulus"`
	Exponent string `json:"exponent"`
}

func rsaEncrypt(publicKey *rsa.PublicKey, text []byte) []byte {
	chunkSize := 2 * (publicKey.N.BitLen()/16 - 1)
	textLen := len(text)
	result := make([][]byte, 0)

	for i := textLen; i > 0; i -= chunkSize {
		textChunk := new(big.Int)
		textChunk.SetBytes(text[max(i-chunkSize, 0):i])
		textChunk.Exp(textChunk, big.NewInt(int64(publicKey.E)), publicKey.N)
		result = append(result, textChunk.Bytes())
	}

	return result[0]
}

func getEncryptedPassword(data *publicKeyData, password string) (string, error) {
	modulus := new(big.Int)
	modulus.SetString(data.Modulus, 16)
	e, err := strconv.ParseInt(data.Exponent, 16, 32)
	if err != nil {
		return "", err
	}
	publicKey := &rsa.PublicKey{
		N: modulus,
		E: int(e),
	}
	return hex.EncodeToString(rsaEncrypt(publicKey, []byte(password))), nil
}
