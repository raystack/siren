package rule

import (
	"errors"
)

var (
	ErrDuplicate = errors.New("rule conflicted with existing")
)
