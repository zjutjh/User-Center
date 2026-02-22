package response

import (
	"errors"
	"reflect"

	"github.com/zjutjh/User-Center/comm"
	pb "github.com/zjutjh/User-Center/idl/user/v1alpha1"
	"github.com/zjutjh/mygo/nlog"
)

func New[T any](code pb.BizCode) (*T, error) {
	resp := new(T)
	reflect.ValueOf(resp).Elem().FieldByName("Code").Set(reflect.ValueOf(code))
	return resp, nil
}

func OK[T any]() (*T, error) {
	return New[T](pb.BizCode_OK)
}

func Error[T any](err error) (*T, error) {
	var bizErr *comm.BizStatusError
	if errors.As(err, &bizErr) {
		return New[T](bizErr.BizCode())
	}
	nlog.Pick().Errorf("未知错误: %v", err)
	return New[T](pb.BizCode_UnknownError)
}
