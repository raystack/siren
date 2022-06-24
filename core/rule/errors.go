package rule

import (
	"errors"
)

var (
	ErrDuplicate = errors.New("rule conflicted with existing")
	ErrRelation  = errors.New("provider's namespace does not exist")
)
