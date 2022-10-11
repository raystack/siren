package httpreceiver

import "context"

//go:generate mockery --name=HTTPCaller -r --case underscore --with-expecter --structname HTTPCaller --filename http_caller.go --output=./mocks
type HTTPCaller interface {
	Notify(ctx context.Context, apiURL string, body []byte) error
}
