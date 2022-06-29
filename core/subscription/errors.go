package subscription

import (
	"errors"
	"fmt"
)

var (
	ErrDuplicate = errors.New("urn already exist")
)

type NotFoundError struct {
	ID uint64
}

func (err NotFoundError) Error() string {
	if err.ID != 0 {
		return fmt.Sprintf("subscription with id %d not found", err.ID)
	}

	return "subscription not found"
}
