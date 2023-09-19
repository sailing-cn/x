package warning

import (
	"errors"
	"google.golang.org/grpc/status"
)

type Warning struct {
	Detail string `json:"detail"`
	Code   uint32 `json:"code"`
}

func (w Warning) Error() string {
	return w.Detail
}

func (w Warning) Message() string {
	return w.Detail
}
func New(text string) *Warning {
	return &Warning{Detail: text, Code: 600}
}

// GRPCError GRPC错误
func GRPCError(text string) error {
	return status.Error(600, text)
}

//func GrpcWarning(msg string) error {
//	return status.Error(600, msg)
//}

func MustOk(err error) (error, bool) {
	if err == nil {
		return nil, true
	}
	if errors.Is(err, Warning{}) {
		return err, false
	}
	return err, false
}

func FromError(err error) (s *Warning, ok bool) {
	if err == nil {
		return nil, true
	}
	if se, ok := err.(interface {
		GRPCStatus() *Warning
	}); ok {
		return se.GRPCStatus(), true
	}
	return New(err.Error()), false
}
