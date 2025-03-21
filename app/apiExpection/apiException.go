package apiExpection

import "net/http"

type Error struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
}

var (
	ServerError       = NewError(http.StatusInternalServerError, 500, "系统异常，请稍后重试!")
	ParamError        = NewError(http.StatusInternalServerError, 501, "参数错误")
	HttpTimeout       = NewError(http.StatusInternalServerError, 505, "系统异常，请稍后重试!")
	UserAlreadyExit   = NewError(http.StatusInternalServerError, 508, "用户名已经存在")
	RequestError      = NewError(http.StatusInternalServerError, 506, "系统异常，请稍后重试!")
	ClosedError       = NewError(http.StatusInternalServerError, 507, "统一身份认证夜间不对外开放")
	NotFound          = NewError(http.StatusNotFound, 404, http.StatusText(http.StatusNotFound))
	Unknown           = NewError(http.StatusInternalServerError, 300500, "系统异常，请稍后重试!")
	EmailNotFound     = NewError(http.StatusOK, 404, "找不到对应邮箱")
	CodeError         = NewError(http.StatusOK, 500, "验证码错误")
	WrongPassword     = NewError(http.StatusOK, 409, "统一系统密码错误")
	NotActivatedError = NewError(http.StatusOK, 411, "统一系统账号未激活")
	WrongAccount      = NewError(http.StatusOK, 412, "统一系统账号错误")
)

func OtherError(message string) *Error {
	return NewError(http.StatusForbidden, 100403, message)
}

func (e *Error) Error() string {
	return e.Msg
}

func NewError(statusCode, Code int, msg string) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       Code,
		Msg:        msg,
	}
}
