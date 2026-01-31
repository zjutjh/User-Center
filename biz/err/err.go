package bizerr

import (
	"fmt"

	typev1 "github.com/zjutjh/User-Center/api/user/v1alpha1"
)

var (
	UnknownError           = New(typev1.BizCode_UnknownError, "系统异常，请稍后重试!")
	ParamError             = New(typev1.BizCode_ParamError, "参数错误")
	PasswordLengthError    = New(typev1.BizCode_PasswordLengthError, "密码长度不符合要求")
	WrongAccountOrPassword = New(typev1.BizCode_WrongAccountOrPassword, "账号或密码错误")
	NotActivated           = New(typev1.BizCode_NotActivated, "账号未激活")
	UserNotExist           = New(typev1.BizCode_UserNotExist, "用户不存在")
	ClosedError            = New(typev1.BizCode_ClosedError, "统一身份认证夜间不对外开放")
	UserExisted            = New(typev1.BizCode_UserExisted, "用户已经存在")
)

type BizStatusErrorIface interface {
	BizCode() int32
	BizMessage() string
	Error() string
}

type BizStatusError struct {
	Code    typev1.BizCode
	Message string
}

func New(code typev1.BizCode, msg string) *BizStatusError {
	return &BizStatusError{
		Code:    code,
		Message: msg,
	}
}

func (b *BizStatusError) BizCode() typev1.BizCode { return b.Code }
func (b *BizStatusError) BizMessage() string      { return b.Message }
func (b *BizStatusError) Error() string           { return fmt.Sprintf("BizError[%d]: %s", b.Code, b.Message) }
