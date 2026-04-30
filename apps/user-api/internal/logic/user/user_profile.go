package user

import (
	"context"

	"github.com/zjutjh/User-Center/apps/user-api/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/types"
	"github.com/zjutjh/User-Center/apps/user-rpc/usercenterservice"
	"github.com/zjutjh/User-Center/common/ctxdata"
	"github.com/zjutjh/User-Center/common/errorsx"
)

func currentUserID(ctx context.Context) (int64, error) {
	userID, ok := ctxdata.UserID(ctx)
	if !ok || userID == 0 {
		return 0, errorsx.ErrNotLoggedIn
	}
	return userID, nil
}

func getUserResp(ctx context.Context, svcCtx *svc.ServiceContext, userID int64) (*types.UserResp, error) {
	rpcResp, err := svcCtx.UserRpc.GetUserInfo(ctx, &usercenterservice.GetUserInfoRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, errorsx.FromGRPC(err)
	}

	return &types.UserResp{
		User: buildUserInfo(rpcResp),
	}, nil
}

func buildUserInfo(info *usercenterservice.GetUserInfoResponse) types.UserInfo {
	return types.UserInfo{
		Id:        info.UserId,
		Username:  info.StudentId,
		StudentId: info.StudentId,
		Bind: types.BindInfo{
			Zf:    info.ZfPassword != "",
			Yxy:   info.YxyUid != "",
			Oauth: info.OauthPassword != "",
		},
		UserType:   info.UserType,
		Email:      info.Email,
		PhoneNum:   info.PhoneNum,
		CreateTime: info.CreateTime,
	}
}
