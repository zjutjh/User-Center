package domain

import (
	"context"

	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/types"
	"github.com/zjutjh/User-Center/apps/user-rpc/usercenterservice"
	"github.com/zjutjh/User-Center/common/errorsx"
)

type UserCenter struct {
	rpc usercenterservice.UserCenterService
}

func NewUserCenter(conf zrpc.RpcClientConf) *UserCenter {
	client := zrpc.MustNewClient(conf)
	return newUserCenterWithRPC(usercenterservice.NewUserCenterService(client))
}

func newUserCenterWithRPC(rpc usercenterservice.UserCenterService) *UserCenter {
	return &UserCenter{rpc: rpc}
}

func (d *UserCenter) Register(ctx context.Context, studentID, password, cardID, email string) error {
	_, err := d.rpc.Register(ctx, &usercenterservice.RegisterRequest{
		StudentId: studentID,
		Password:  password,
		CardId:    cardID,
		Email:     email,
	})
	return errorsx.FromGRPC(err)
}

func (d *UserCenter) Login(ctx context.Context, studentID, password string) (int64, string, error) {
	resp, err := d.rpc.Login(ctx, &usercenterservice.LoginRequest{
		StudentId: studentID,
		Password:  password,
	})
	if err != nil {
		return 0, "", errorsx.FromGRPC(err)
	}
	return resp.UserId, resp.Session, nil
}

func (d *UserCenter) BindOauth(ctx context.Context, userID int64, oauthPassword string) error {
	_, err := d.rpc.BindOauth(ctx, &usercenterservice.BindOauthRequest{
		UserId:        userID,
		OauthPassword: oauthPassword,
	})
	return errorsx.FromGRPC(err)
}

func (d *UserCenter) BindZf(ctx context.Context, userID int64, zfPassword string) error {
	_, err := d.rpc.BindZf(ctx, &usercenterservice.BindZfRequest{
		UserId:     userID,
		ZfPassword: zfPassword,
	})
	return errorsx.FromGRPC(err)
}

func (d *UserCenter) BindYxy(ctx context.Context, userID int64, deviceID, yxyUID string) error {
	_, err := d.rpc.BindYxy(ctx, &usercenterservice.BindYxyRequest{
		UserId:   userID,
		DeviceId: deviceID,
		YxyUid:   yxyUID,
	})
	return errorsx.FromGRPC(err)
}

func (d *UserCenter) GetUser(ctx context.Context, userID int64) (*types.UserResp, error) {
	resp, err := d.rpc.GetUserInfo(ctx, &usercenterservice.GetUserInfoRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, errorsx.FromGRPC(err)
	}
	return userFromRPC(resp), nil
}

func (d *UserCenter) ResetPassword(ctx context.Context, userID int64, studentID, cardID, password string) error {
	_, err := d.rpc.ResetPassword(ctx, &usercenterservice.ResetPasswordRequest{
		UserId:    userID,
		StudentId: studentID,
		CardId:    cardID,
		Password:  password,
	})
	return errorsx.FromGRPC(err)
}

func (d *UserCenter) DeleteAccount(ctx context.Context, userID int64, studentID, cardID string) error {
	_, err := d.rpc.DeleteAccount(ctx, &usercenterservice.DeleteAccountRequest{
		UserId:    userID,
		StudentId: studentID,
		CardId:    cardID,
	})
	return errorsx.FromGRPC(err)
}

func userFromRPC(info *usercenterservice.GetUserInfoResponse) *types.UserResp {
	return &types.UserResp{
		User: types.UserInfo{
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
		},
	}
}
