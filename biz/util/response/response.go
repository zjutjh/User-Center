package response

import (
	"errors"
	"log/slog"

	"google.golang.org/protobuf/types/known/structpb"

	typev1 "github.com/zjutjh/User-Center/api/user/v1alpha1"
	"github.com/zjutjh/User-Center/biz/err"
)

func Success(data interface{}) (*typev1.Response, error) {
	value, err := structpb.NewValue(data)
	if err != nil {
		slog.Error("failed to convert data to protobuf Value: %v", err)
		return nil, bizerr.UnknownError
	}
	return &typev1.Response{
		Code:    typev1.BizCode_OK,
		Message: "success",
		Data:    value,
	}, nil
}

func OK() (*typev1.Response, error) {
	return Success(nil)
}

func Error(err error) (*typev1.Response, error) {
	return ErrorExtra(err, "")
}

func ErrorExtra(err error, extra string) (*typev1.Response, error) {
	var bizErr *bizerr.BizStatusError
	if errors.As(err, &bizErr) {
		return &typev1.Response{
			Code:    bizErr.BizCode(),
			Message: extra,
		}, nil
	}
	slog.Error("unknown error: %v", err)
	return &typev1.Response{
		Code:    typev1.BizCode_UnknownError,
		Message: extra,
	}, nil
}
