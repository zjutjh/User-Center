package errorsx

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestGRPCRoundTrip(t *testing.T) {
	err := ToGRPC(ErrWrongAccountOrPassword)

	require.ErrorIs(t, FromGRPC(err), ErrWrongAccountOrPassword)
}

func TestGRPCRoundTripWithValidationError(t *testing.T) {
	err := ToGRPC(ErrPasswordLength)

	require.ErrorIs(t, FromGRPC(err), ErrPasswordLength)
}

func TestFromGRPCMapsUnavailableToThirdService(t *testing.T) {
	err := status.Error(codes.Unavailable, "rpc unavailable")

	require.ErrorIs(t, FromGRPC(err), ErrThirdService)
}

func TestGRPCRoundTripOverNetwork(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer lis.Close()

	s := grpc.NewServer()
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: "test.TestService",
		HandlerType: (*interface{})(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Fail",
				Handler: func(_ interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
					in := new(emptypb.Empty)
					require.NoError(t, dec(in))

					handler := func(context.Context, any) (any, error) {
						return nil, ToGRPC(ErrPasswordLength)
					}
					if interceptor == nil {
						return handler(ctx, in)
					}

					return interceptor(ctx, in, &grpc.UnaryServerInfo{
						FullMethod: "/test.TestService/Fail",
					}, handler)
				},
			},
		},
	}, nil)

	go func() {
		_ = s.Serve(lis)
	}()
	defer s.GracefulStop()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	callErr := conn.Invoke(context.Background(), "/test.TestService/Fail", &emptypb.Empty{}, &emptypb.Empty{})
	require.ErrorIs(t, FromGRPC(callErr), ErrPasswordLength)
}
