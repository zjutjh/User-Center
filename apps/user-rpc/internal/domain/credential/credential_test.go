package credential

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const testKey = "1234567890abcdef"

func TestCodecEncryptDecrypt(t *testing.T) {
	codec, err := NewCodec(testKey)
	require.NoError(t, err)

	plaintext := "secret-password"
	ciphertext, err := codec.Encrypt(plaintext)
	require.NoError(t, err)

	require.NotEmpty(t, ciphertext)
	require.NotEqual(t, plaintext, ciphertext)

	decrypted, err := codec.Decrypt(ciphertext)
	require.NoError(t, err)
	require.Equal(t, plaintext, decrypted)
}

func TestCodecEncryptAddsVersionPrefix(t *testing.T) {
	codec, err := NewCodec(testKey)
	require.NoError(t, err)

	ciphertext, err := codec.Encrypt("secret-password")
	require.NoError(t, err)

	require.True(t, strings.HasPrefix(ciphertext, "v1:"))
}

func TestCodecDecryptLegacyPlaintext(t *testing.T) {
	codec, err := NewCodec(testKey)
	require.NoError(t, err)

	plaintext := "legacy-password"
	decrypted, err := codec.Decrypt(plaintext)
	require.NoError(t, err)
	require.Equal(t, plaintext, decrypted)
}

func TestServicePrepareOauthBindEncryptsPassword(t *testing.T) {
	service, err := NewService(testKey)
	require.NoError(t, err)

	plaintext := "secret-password"
	ciphertext, err := service.PrepareOauthBind(plaintext)
	require.NoError(t, err)

	require.NotEmpty(t, ciphertext)
	require.NotEqual(t, plaintext, ciphertext)

	decrypted, err := service.DecryptPassword(ciphertext)
	require.NoError(t, err)
	require.Equal(t, plaintext, decrypted)
}

func TestServicePrepareYxyBindReturnsTrimmedIdentifiers(t *testing.T) {
	service, err := NewService(testKey)
	require.NoError(t, err)

	deviceID, yxyUID, err := service.PrepareYxyBind("  device-id  ", "  yxy-uid  ")
	require.NoError(t, err)

	require.Equal(t, "device-id", deviceID)
	require.Equal(t, "yxy-uid", yxyUID)
}
