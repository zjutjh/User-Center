package oauth

import (
	"context"
	"net/http"

	oauthclient "github.com/zjutjh/User-Center/apps/user-rpc/internal/httpclient/oauth"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Login(ctx context.Context, username, password string) ([]*http.Cookie, error) {
	return oauthclient.New(ctx).Login(username, password)
}

func (s *Service) SSOLogin(ctx context.Context, cookies []*http.Cookie, serviceURL string) ([]*http.Cookie, error) {
	return oauthclient.New(ctx).SSOLogin(cookies, serviceURL)
}
