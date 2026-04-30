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

type BindZfLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 绑定正方密码
func NewBindZfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindZfLogic {
	return &BindZfLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindZfLogic) BindZf(req *types.BindPasswordReq) (resp *types.EmptyResp, err error) {
	userID, err := currentUserID(l.ctx)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.Password) == "" {
		return nil, errorsx.ErrParameterInvalid
	}

	return bindUser(l.ctx, l.svcCtx, &pb.BindRequest{
		UserId:     userID,
		Type:       pb.BindType_BIND_TYPE_ZF,
		ZfPassword: req.Password,
	})
}
