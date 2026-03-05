// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"
	"strings"

	"github.com/zjutjh/User-Center/apps/user-api/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/types"
	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
	"github.com/zjutjh/User-Center/apps/user-rpc/usercenterservice"
	"github.com/zjutjh/User-Center/common/ctxdata"
	"github.com/zjutjh/User-Center/common/errorsx"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 绑定外部系统
func NewBindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindLogic {
	return &BindLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindLogic) Bind(req *types.BindReq) (resp *types.BindResp, err error) {
	userID, ok := ctxdata.UserID(l.ctx)
	if !ok || userID == 0 {
		return nil, errorsx.ErrNotLoggedIn
	}

	bindType, err := parseBindType(req.Type)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.UserRpc.Bind(l.ctx, &usercenterservice.BindRequest{
		UserId:        userID,
		Type:          bindType,
		DeviceId:      req.DeviceId,
		YxyUid:        req.YxyUid,
		ZfPassword:    req.ZfPassword,
		OauthPassword: req.OauthPassword,
	})
	if err != nil {
		l.Errorf("调用用户 RPC 绑定接口失败: %v", err)
		return nil, errorsx.FromGRPC(err)
	}

	return &types.BindResp{}, nil
}

func parseBindType(bindType string) (pb.BindType, error) {
	switch strings.ToLower(strings.TrimSpace(bindType)) {
	case "yxy":
		return pb.BindType_BIND_TYPE_YXY, nil
	case "zf":
		return pb.BindType_BIND_TYPE_ZF, nil
	case "oauth":
		return pb.BindType_BIND_TYPE_OAUTH, nil
	default:
		return pb.BindType_BIND_TYPE_UNSPECIFIED, errorsx.ErrParameterInvalid
	}
}
