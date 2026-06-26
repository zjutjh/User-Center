package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/svc"
	"github.com/zjutjh/User-Center/common/errorsx"
)

func ensureBindUserExists(ctx context.Context, svcCtx *svc.ServiceContext, logger logx.Logger, userID int64) error {
	_, err := svcCtx.UserRepo.GetUserById(ctx, userID)
	if err != nil {
		if codeErr, ok := errorsx.As(err); ok {
			return codeErr
		}
		logger.Errorf("根据用户 ID 查询用户失败: %v", err)
		return errorsx.ErrUnknown
	}

	return nil
}
