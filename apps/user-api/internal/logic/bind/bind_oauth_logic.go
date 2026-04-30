// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package bind

import (
	"context"
	"strings"

	"github.com/zjutjh/User-Center/apps/user-api/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/types"
	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
	"github.com/zjutjh/User-Center/common/errorsx"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindOauthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 绑定统一认证密码
func NewBindOauthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindOauthLogic {
	return &BindOauthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindOauthLogic) BindOauth(req *types.BindPasswordReq) (resp *types.EmptyResp, err error) {
	userID, err := currentUserID(l.ctx)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.Password) == "" {
		return nil, errorsx.ErrParameterInvalid
	}

	return bindUser(l.ctx, l.svcCtx, &pb.BindRequest{
		UserId:        userID,
		Type:          pb.BindType_BIND_TYPE_OAUTH,
		OauthPassword: req.Password,
	})
}
