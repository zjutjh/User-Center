package bind

import (
	"context"

	"github.com/zjutjh/User-Center/apps/user-api/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/types"
	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
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

func bindUser(ctx context.Context, svcCtx *svc.ServiceContext, req *pb.BindRequest) (*types.EmptyResp, error) {
	if _, err := svcCtx.UserRpc.Bind(ctx, req); err != nil {
		return nil, errorsx.FromGRPC(err)
	}
	return &types.EmptyResp{}, nil
}
