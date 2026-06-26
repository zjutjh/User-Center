package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zjutjh/User-Center/apps/user-rpc/usercenterservice"
)

func TestUserRespFromRPCBuildsBindFlags(t *testing.T) {
	resp := userFromRPC(&usercenterservice.GetUserInfoResponse{
		UserId:        7,
		StudentId:     "20240001",
		UserType:      "student",
		Email:         "mango@example.com",
		PhoneNum:      "13800000000",
		YxyUid:        "yxy-uid",
		OauthPassword: "bound",
		CreateTime:    "2026-06-27T00:00:00Z",
	})

	require.Equal(t, int64(7), resp.User.Id)
	require.Equal(t, "20240001", resp.User.Username)
	require.True(t, resp.User.Bind.Yxy)
	require.True(t, resp.User.Bind.Oauth)
	require.False(t, resp.User.Bind.Zf)
}
