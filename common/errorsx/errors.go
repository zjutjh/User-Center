package errorsx

import (
	"errors"
	"fmt"
)

const (
	CodeOK               int64 = 0
	CodeUnknown          int64 = 10000
	CodeThirdService     int64 = 10001
	CodeNotLoggedIn      int64 = 20000
	CodeParameterInvalid int64 = 20003
	CodePasswordLength   int64 = 40001
	CodeWrongAccount     int64 = 40002
	CodeNotActivated     int64 = 40003
	CodeUserNotExist     int64 = 40004
	CodeOAuthClosed      int64 = 40005
	CodeUserExisted      int64 = 40006
	CodePasswordEdited   int64 = 40007
	CodeInvalidCookie    int64 = 40008
)

type CodeError struct {
	Code    int64
	Message string
}

func New(code int64, msg string) *CodeError {
	return &CodeError{
		Code:    code,
		Message: msg,
	}
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("code=%d msg=%s", e.Code, e.Message)
}

var (
	ErrUnknown                = New(CodeUnknown, "系统异常，请稍后重试")
	ErrThirdService           = New(CodeThirdService, "下游服务调用失败")
	ErrNotLoggedIn            = New(CodeNotLoggedIn, "用户未登录或登录已过期")
	ErrParameterInvalid       = New(CodeParameterInvalid, "参数非法")
	ErrPasswordLength         = New(CodePasswordLength, "密码长度必须在6~20位之间")
	ErrWrongAccountOrPassword = New(CodeWrongAccount, "账号或密码错误")
	ErrNotActivated           = New(CodeNotActivated, "账号未激活")
	ErrUserNotExist           = New(CodeUserNotExist, "用户不存在")
	ErrOAuthClosed            = New(CodeOAuthClosed, "统一身份认证夜间不对外开放")
	ErrUserExisted            = New(CodeUserExisted, "用户已经存在")
	ErrPasswordNeedEdited     = New(CodePasswordEdited, "密码需要修改")
	ErrInvalidCookie          = New(CodeInvalidCookie, "无效的 Cookie")
)

func As(err error) (*CodeError, bool) {
	var target *CodeError
	if errors.As(err, &target) {
		return target, true
	}
	return nil, false
}

func CodeOf(err error) (int64, string) {
	if err == nil {
		return CodeOK, "ok"
	}
	if target, ok := As(err); ok {
		return target.Code, target.Message
	}
	return ErrUnknown.Code, ErrUnknown.Message
}

func HTTPStatus(error) int {
	return 200
}

func FromCode(code int64) *CodeError {
	switch code {
	case CodeOK:
		return nil
	case CodeUnknown:
		return ErrUnknown
	case CodeThirdService:
		return ErrThirdService
	case CodeNotLoggedIn:
		return ErrNotLoggedIn
	case CodeParameterInvalid:
		return ErrParameterInvalid
	case CodePasswordLength:
		return ErrPasswordLength
	case CodeWrongAccount:
		return ErrWrongAccountOrPassword
	case CodeNotActivated:
		return ErrNotActivated
	case CodeUserNotExist:
		return ErrUserNotExist
	case CodeOAuthClosed:
		return ErrOAuthClosed
	case CodeUserExisted:
		return ErrUserExisted
	case CodePasswordEdited:
		return ErrPasswordNeedEdited
	case CodeInvalidCookie:
		return ErrInvalidCookie
	default:
		return ErrUnknown
	}
}
