package plugins

import "errors"

var (
	ErrNotImplemented                 = errors.New("operation not supported")
	ErrProviderSyncMethodNotSupported = errors.New("provider sync method not supported")
)
