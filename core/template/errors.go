package template

import (
	"errors"
	"fmt"
)

var (
	ErrDuplicate = errors.New("name already exist")
)

type NotFoundError struct {
	Name string
}

func (err NotFoundError) Error() string {
	if err.Name != "" {
		return fmt.Sprintf("template with name %q not found", err.Name)
	}

	return "template not found"
}
