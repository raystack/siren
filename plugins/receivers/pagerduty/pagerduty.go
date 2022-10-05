package pagerduty

import "context"

//go:generate mockery --name=PagerDutyCaller -r --case underscore --with-expecter --structname PagerDutyCaller --filename pagerduty_caller.go --output=./mocks
type PagerDutyCaller interface {
	NotifyV1(ctx context.Context, message MessageV1) error
}
