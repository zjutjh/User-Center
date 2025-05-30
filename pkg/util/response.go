package util

import (
	"log/slog"

	"google.golang.org/protobuf/types/known/structpb"

	typev1 "github.com/zjutjh/User-Center-grpc/api/types/v1alpha1"
	"github.com/zjutjh/User-Center-grpc/pkg/apiExpection"
)

func ResponseSuccess(data interface{}) (*typev1.Response, error) {
	if data == nil {
		return &typev1.Response{
			Code:    200,
			Message: "success",
		}, nil
	}
	value, err := structpb.NewValue(data)
	if err != nil {
		slog.Error("failed to convert data to protobuf Value: %v", err)
		return nil, apiExpection.ServerError
	}
	return &typev1.Response{
		Code:    200,
		Message: "success",
		Data:    value,
	}, nil
}
