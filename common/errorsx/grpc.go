package errorsx

import (
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const grpcErrorDomain = "user-center"

func ToGRPC(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := status.FromError(err); ok {
		return err
	}

	codeErr, ok := As(err)
	if !ok {
		codeErr = ErrUnknown
	}

	st := status.New(toGRPCCode(codeErr), codeErr.Message)
	withDetails, detailErr := st.WithDetails(&errdetails.ErrorInfo{
		Domain: grpcErrorDomain,
		Reason: strconv.FormatInt(codeErr.Code, 10),
		Metadata: map[string]string{
			"code": strconv.FormatInt(codeErr.Code, 10),
			"msg":  codeErr.Message,
		},
	})
	if detailErr != nil {
		return st.Err()
	}

	return withDetails.Err()
}

func FromGRPC(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return ErrThirdService
	}

	for _, detail := range st.Details() {
		info, ok := detail.(*errdetails.ErrorInfo)
		if !ok || info.Domain != grpcErrorDomain {
			continue
		}

		codeValue := info.Metadata["code"]
		if codeValue == "" {
			codeValue = info.Reason
		}

		code, parseErr := strconv.ParseInt(codeValue, 10, 64)
		if parseErr == nil {
			if mapped := FromCode(code); mapped != nil {
				return mapped
			}

			msg := info.Metadata["msg"]
			if msg == "" {
				msg = st.Message()
			}
			return New(code, msg)
		}
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return ErrParameterInvalid
	case codes.Unauthenticated:
		return ErrNotLoggedIn
	case codes.NotFound:
		return ErrUserNotExist
	case codes.Unavailable, codes.DeadlineExceeded:
		return ErrThirdService
	default:
		return ErrUnknown
	}
}

func toGRPCCode(err *CodeError) codes.Code {
	switch err.Code {
	case CodeParameterInvalid, CodePasswordLength:
		return codes.InvalidArgument
	case CodeNotLoggedIn, CodeInvalidCookie:
		return codes.Unauthenticated
	case CodeUserNotExist:
		return codes.NotFound
	default:
		return codes.Internal
	}
}
