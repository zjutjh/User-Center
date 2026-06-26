package credential

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

const ciphertextVersionPrefix = "v1:"

type Codec struct {
	gcm cipher.AEAD
}

func NewCodec(secretKey string) (*Codec, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &Codec{gcm: gcm}, nil
}

func (c *Codec) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := c.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return ciphertextVersionPrefix + base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (c *Codec) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	if !strings.HasPrefix(ciphertext, ciphertextVersionPrefix) {
		return ciphertext, nil
	}

	encoded := strings.TrimPrefix(ciphertext, ciphertextVersionPrefix)
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	nonceSize := c.gcm.NonceSize()
	if len(raw) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce := raw[:nonceSize]
	encrypted := raw[nonceSize:]

	plaintext, err := c.gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
