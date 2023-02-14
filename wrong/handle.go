package wrong

import (
	"errors"
	"google.golang.org/grpc/status"
)

func ShouldOk(err error) (error, bool) {
	if err == nil {
		return nil, true
	}
	if errors.Is(err, Warning{}) {
		return status.Error(600, err.Error()), false
	}
	return err, false
}

func New(text string) error {
	return status.Error(600, text)
}
