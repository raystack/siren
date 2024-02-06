package pagerduty

import "context"

type PagerDutyCaller interface {
	NotifyV1(ctx context.Context, message MessageV1) error
}
