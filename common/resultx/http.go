package resultx

import (
	"context"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zjutjh/User-Center/common/errorsx"
	"github.com/zjutjh/User-Center/common/validatorx"
)

type Response struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func InstallHTTPHandlers() {
	httpx.SetValidator(validatorx.HTTPValidator{})
	httpx.SetOkHandler(func(_ context.Context, v any) any {
		return Success(v)
	})
	httpx.SetErrorHandlerCtx(func(_ context.Context, err error) (int, any) {
		err = errorsx.NormalizeHTTPError(err)
		code, msg := errorsx.CodeOf(err)
		return errorsx.HTTPStatus(err), Response{
			Code: code,
			Msg:  msg,
			Data: struct{}{},
		}
	})
}

func Success(v any) Response {
	if v == nil {
		v = struct{}{}
	}
	return Response{
		Code: errorsx.CodeOK,
		Msg:  "ok",
		Data: v,
	}
}
