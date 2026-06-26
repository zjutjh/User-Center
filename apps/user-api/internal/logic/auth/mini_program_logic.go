// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"

	"github.com/zjutjh/User-Center/apps/user-api/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MiniProgramLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 小程序登录
func NewMiniProgramLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MiniProgramLogic {
	return &MiniProgramLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MiniProgramLogic) MiniProgram(req *types.MiniProgramLoginReq) (resp *types.UserResp, err error) {
	// todo: add your logic here and delete this line

	return
}
