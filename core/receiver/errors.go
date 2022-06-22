package receiver

import (
	"fmt"
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
