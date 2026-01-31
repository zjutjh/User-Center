package util

import (
	"log/slog"

	"google.golang.org/protobuf/types/known/structpb"

	typev1 "github.com/zjutjh/User-Center/api/user/v1alpha1"
	"github.com/zjutjh/User-Center/pkg/exception"
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
		return nil, exception.ServerError
	}
	return &typev1.Response{
		Code:    200,
		Message: "success",
		Data:    value,
	}, nil
}
