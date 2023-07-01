package receiver

import (
	"fmt"

	"github.com/raystack/siren/pkg/errors"
)

var (
	ErrNotImplemented = errors.New("operation not supported")
)

type NotFoundError struct {
	ID uint64
}

func (err NotFoundError) Error() string {
	if err.ID != 0 {
		return fmt.Sprintf("receiver with id %d not found", err.ID)
	}

	return "receiver not found"
}
