package warning

import (
	"errors"
)

type Warning struct {
	Detail string
	Code   uint32
}

func (w Warning) Error() string {
	return w.Detail
}

func (w Warning) Message() string {
	return w.Detail
}
func New(text string) error {
	return &Warning{Detail: text, Code: 600}
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
