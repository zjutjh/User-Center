package apiExpection

import "net/http"

type Error struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
}

var (
	ServerError     = NewError(http.StatusInternalServerError, 500, "系统异常，请稍后重试!")
	ParamError      = NewError(http.StatusInternalServerError, 501, "参数错误")
	HttpTimeout     = NewError(http.StatusInternalServerError, 505, "系统异常，请稍后重试!")
	UserAlreadyExit = NewError(http.StatusInternalServerError, 508, "用户名已经存在")
	RequestError    = NewError(http.StatusInternalServerError, 506, "系统异常，请稍后重试!")
	NotFound        = NewError(http.StatusNotFound, 404, http.StatusText(http.StatusNotFound))
	Unknown         = NewError(http.StatusInternalServerError, 300500, "系统异常，请稍后重试!")
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
