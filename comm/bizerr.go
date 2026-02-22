package comm

import (
	"fmt"

	typev1 "github.com/zjutjh/User-Center/idl/user/v1alpha1"
	"github.com/zjutjh/mygo/kit"
)

var (
	UnknownError           = New(typev1.BizCode_UnknownError, "系统异常，请稍后重试!")
	PasswordLengthError    = New(typev1.BizCode_PasswordLengthError, "密码长度不符合要求")
	WrongAccountOrPassword = New(typev1.BizCode_WrongAccountOrPassword, "账号或密码错误")
	NotActivated           = New(typev1.BizCode_NotActivated, "账号未激活")
	UserNotExist           = New(typev1.BizCode_UserNotExist, "用户不存在")
	ClosedError            = New(typev1.BizCode_ClosedError, "统一身份认证夜间不对外开放")
	UserExisted            = New(typev1.BizCode_UserExisted, "用户已经存在")
)

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
func (b *BizStatusError) ToKitCode() kit.Code {
	return kit.NewCode(int64(b.Code.Number()), b.Message)
}

func FromBizCode(code typev1.BizCode) kit.Code {
	switch code {
	case typev1.BizCode_OK:
		return CodeOK
	case typev1.BizCode_UnknownError:
		return UnknownError.ToKitCode()
	case typev1.BizCode_PasswordLengthError:
		return PasswordLengthError.ToKitCode()
	case typev1.BizCode_WrongAccountOrPassword:
		return WrongAccountOrPassword.ToKitCode()
	case typev1.BizCode_NotActivated:
		return NotActivated.ToKitCode()
	case typev1.BizCode_UserNotExist:
		return UserNotExist.ToKitCode()
	case typev1.BizCode_ClosedError:
		return ClosedError.ToKitCode()
	case typev1.BizCode_UserExisted:
		return UserExisted.ToKitCode()
	}
	return UnknownError.ToKitCode()
}
