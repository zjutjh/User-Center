package credential

import (
	"strings"

	"github.com/zjutjh/User-Center/common/errorsx"
)

type Service struct {
	codec *Codec
}

func NewService(secretKey string) (*Service, error) {
	codec, err := NewCodec(secretKey)
	if err != nil {
		return nil, err
	}

	return &Service{codec: codec}, nil
}

func (s *Service) DecryptPassword(ciphertext string) (string, error) {
	return s.codec.Decrypt(ciphertext)
}

func (s *Service) PrepareYxyBind(deviceIDValue, yxyUIDValue string) (string, string, error) {
	deviceID := strings.TrimSpace(deviceIDValue)
	yxyUID := strings.TrimSpace(yxyUIDValue)
	if deviceID == "" || yxyUID == "" {
		return "", "", errorsx.ErrParameterInvalid
	}

	return deviceID, yxyUID, nil
}

func (s *Service) PrepareZfBind(zfPasswordValue string) (string, error) {
	zfPassword := strings.TrimSpace(zfPasswordValue)
	if zfPassword == "" {
		return "", errorsx.ErrParameterInvalid
	}

	ciphertext, err := s.codec.Encrypt(zfPassword)
	if err != nil {
		return "", err
	}

	return ciphertext, nil
}

func (s *Service) PrepareOauthBind(oauthPasswordValue string) (string, error) {
	oauthPassword := strings.TrimSpace(oauthPasswordValue)
	if oauthPassword == "" {
		return "", errorsx.ErrParameterInvalid
	}

	ciphertext, err := s.codec.Encrypt(oauthPassword)
	if err != nil {
		return "", err
	}

	return ciphertext, nil
}
