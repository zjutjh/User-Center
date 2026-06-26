package credential

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
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

func TestServiceBuildBindUpdatesEncryptsOAuthPassword(t *testing.T) {
	service, err := NewService(testKey)
	require.NoError(t, err)

	plaintext := "secret-password"
	updates, err := service.BuildBindUpdates(&pb.BindRequest{
		Type:          pb.BindType_BIND_TYPE_OAUTH,
		OauthPassword: plaintext,
	})
	require.NoError(t, err)

	ciphertext, ok := updates["oauth_password"].(string)
	require.True(t, ok)
	require.NotEmpty(t, ciphertext)
	require.NotEqual(t, plaintext, ciphertext)

	decrypted, err := service.DecryptPassword(ciphertext)
	require.NoError(t, err)
	require.Equal(t, plaintext, decrypted)
}

func TestServiceBuildBindUpdatesYXYReturnsIdentifiers(t *testing.T) {
	service, err := NewService(testKey)
	require.NoError(t, err)

	updates, err := service.BuildBindUpdates(&pb.BindRequest{
		Type:     pb.BindType_BIND_TYPE_YXY,
		DeviceId: "device-id",
		YxyUid:   "yxy-uid",
	})
	require.NoError(t, err)

	require.Equal(t, map[string]any{
		"device_id": "device-id",
		"yxy_uid":   "yxy-uid",
	}, updates)
}
