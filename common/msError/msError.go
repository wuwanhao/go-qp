package msError

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Error struct {
	Code int
	Err  error
}

func NewError(code int, err error) *Error {
	return &Error{
		Code: code,
		Err:  err,
	}
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func GrpcError(err *Error) error {
	return status.Error(codes.Code(err.Code), err.Err.Error())
}

func ParseGrpcError(err Error) (int, string) {
	fromError, _ := status.FromError(&err)
	return int(fromError.Code()), fromError.Message()
}

func ToError(err error) *Error {
	fromError, _ := status.FromError(err)
	return NewError(int(fromError.Code()), errors.New(fromError.Message()))
}
