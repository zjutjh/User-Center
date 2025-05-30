package apiExpection

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	typev1 "github.com/zjutjh/User-Center-grpc/api/types/v1alpha1"
)

type Error struct {
	StatusCode codes.Code `json:"-"`
	Code       int        `json:"code"`
	Msg        string     `json:"msg"`
}

var (
	ServerError         = NewError(codes.Internal, 500, "系统异常，请稍后重试!")
	ParamError          = NewError(codes.InvalidArgument, 501, "参数错误")
	HttpTimeout         = NewError(codes.DeadlineExceeded, 505, "系统异常，请稍后重试!")
	RequestError        = NewError(codes.Canceled, 506, "系统异常，请稍后重试!")
	NotFound            = NewError(codes.NotFound, 404, http.StatusText(http.StatusNotFound))
	Unknown             = NewError(codes.Unknown, 500, "系统异常，请稍后重试!")
	WrongPassword       = NewError(codes.OK, 409, "统一系统密码错误")
	NotActivatedError   = NewError(codes.OK, 411, "统一系统账号未激活")
	WrongAccount        = NewError(codes.OK, 412, "统一系统账号错误")
	ClosedError         = NewError(codes.OK, 507, "统一身份认证夜间不对外开放")
	UserNotFound        = NewError(codes.OK, 400, "学号和身份证不存在或者不匹配")
	PasswordLengthError = NewError(codes.OK, 401, "密码长度不符合要求")
	AuthError           = NewError(codes.OK, 407, "密码错误")
	UserAlreadyExit     = NewError(codes.OK, 403, "用户已经存在")
	UserNotExit         = NewError(codes.OK, 404, "用户不存在")
)

func OtherError(message string) *Error {
	return NewError(http.StatusForbidden, 100403, message)
}

func (e *Error) Error() string {
	return e.Msg
}

func NewError(statusCode codes.Code, Code int, msg string) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       Code,
		Msg:        msg,
	}
}

func (e *Error) ToResponse() (*typev1.Response, error) {
	if e.StatusCode != codes.OK {
		err := status.Errorf(e.StatusCode, e.Msg)
		return nil, err
	}
	response := &typev1.Response{
		Code:    int32(e.Code),
		Message: e.Msg,
	}
	return response, nil
}
