package interceptors

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zjutjh/User-Center/common/errorsx"
	"google.golang.org/grpc"
)

func UnaryRequestLogInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		cost := time.Since(start)
		if err != nil {
			logc.Errorf(ctx, "gRPC 请求失败，方法=%s，耗时=%s，错误=%v", info.FullMethod, cost, err)
		} else {
			logc.Infof(ctx, "gRPC 请求完成，方法=%s，耗时=%s", info.FullMethod, cost)
		}
		return resp, err
	}
}

func UnaryErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return nil, errorsx.ToGRPC(err)
		}
		return resp, nil
	}
}
