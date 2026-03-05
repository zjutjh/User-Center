package oauth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zjutjh/User-Center/common/errorsx"
)

func TestLogin(t *testing.T) {
	username := os.Getenv("UC_OAUTH_USERNAME")
	password := os.Getenv("UC_OAUTH_PASSWORD")
	if username == "" || password == "" {
		t.Skip("set UC_OAUTH_USERNAME and UC_OAUTH_PASSWORD to run integration test")
	}

	for range 3 {
		cookies, err := New(t.Context()).Login(username, password)
		require.NoError(t, err)
		require.Len(t, cookies, 1)
		require.Equal(t, "iPlanetDirectoryPro", cookies[0].Name)
		require.NotEmpty(t, cookies[0].Value)
	}

	for range 3 {
		cookies, err := New(t.Context()).Login(username, password+"1")
		require.ErrorIs(t, err, errorsx.ErrWrongAccountOrPassword)
		require.Nil(t, cookies)
	}
}

func TestSSOLogin(t *testing.T) {
	username := os.Getenv("UC_OAUTH_USERNAME")
	password := os.Getenv("UC_OAUTH_PASSWORD")
	if username == "" || password == "" {
		t.Skip("set UC_OAUTH_USERNAME and UC_OAUTH_PASSWORD to run integration test")
	}

	loginCookies, err := New(t.Context()).Login(username, password)
	require.NoError(t, err)
	require.NotEmpty(t, loginCookies)

	serviceURLs := []string{
		"http://www.gdjw.zjut.edu.cn/sso/zfiotlogin",
		"http://www.me.zjut.edu.cn/personal-center",
	}

	for _, serviceURL := range serviceURLs {
		ssoCookies, ssoErr := New(t.Context()).SSOLogin(loginCookies, serviceURL)
		require.NoError(t, ssoErr)
		require.NotEmpty(t, ssoCookies)
	}
}
