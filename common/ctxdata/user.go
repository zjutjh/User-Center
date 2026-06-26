package ctxdata

import (
	"context"

	"github.com/zjutjh/User-Center/common/constants"
)

func WithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, constants.ContextKeyUserID, userID)
}

func UserID(ctx context.Context) (int64, bool) {
	value := ctx.Value(constants.ContextKeyUserID)
	userID, ok := value.(int64)
	return userID, ok
}
