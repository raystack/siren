package cli

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/siren/pkg/errors"
)

var (
	ErrClientConfigNotFound = errors.New(heredoc.Doc(`
		Siren client config not found.
		Run "siren config init" to initialize a new client config or
		Run "siren help environment" for more information.
	`))
)
