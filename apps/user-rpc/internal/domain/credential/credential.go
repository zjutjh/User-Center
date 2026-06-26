package credential

import (
	"strings"

	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
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

func (s *Service) BuildBindUpdates(in *pb.BindRequest) (map[string]any, error) {
	switch in.GetType() {
	case pb.BindType_BIND_TYPE_YXY:
		return s.buildYXYBindUpdates(in)
	case pb.BindType_BIND_TYPE_ZF:
		return s.buildZFBindUpdates(in)
	case pb.BindType_BIND_TYPE_OAUTH:
		return s.buildOAuthBindUpdates(in)
	default:
		return nil, errorsx.ErrParameterInvalid
	}
}

func (s *Service) DecryptPassword(ciphertext string) (string, error) {
	return s.codec.Decrypt(ciphertext)
}

func (s *Service) buildYXYBindUpdates(in *pb.BindRequest) (map[string]any, error) {
	deviceID := strings.TrimSpace(in.GetDeviceId())
	yxyUID := strings.TrimSpace(in.GetYxyUid())
	if deviceID == "" || yxyUID == "" {
		return nil, errorsx.ErrParameterInvalid
	}

	return map[string]any{
		"device_id": deviceID,
		"yxy_uid":   yxyUID,
	}, nil
}

func (s *Service) buildZFBindUpdates(in *pb.BindRequest) (map[string]any, error) {
	zfPassword := strings.TrimSpace(in.GetZfPassword())
	if zfPassword == "" {
		return nil, errorsx.ErrParameterInvalid
	}

	ciphertext, err := s.codec.Encrypt(zfPassword)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"zf_password": ciphertext,
	}, nil
}

func (s *Service) buildOAuthBindUpdates(in *pb.BindRequest) (map[string]any, error) {
	oauthPassword := strings.TrimSpace(in.GetOauthPassword())
	if oauthPassword == "" {
		return nil, errorsx.ErrParameterInvalid
	}

	ciphertext, err := s.codec.Encrypt(oauthPassword)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"oauth_password": ciphertext,
	}, nil
}
