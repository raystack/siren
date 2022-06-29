package namespace

import (
	"errors"
	"fmt"
)

var (
	ErrDuplicate = errors.New("urn and provider pair already exist")
)

type NotFoundError struct {
	ID uint64
}

func (err NotFoundError) Error() string {
	if err.ID != 0 {
		return fmt.Sprintf("namespace with id %d not found", err.ID)
	}

	return "namespace not found"
}
