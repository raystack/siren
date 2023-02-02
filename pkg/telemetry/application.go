package telemetry

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

const (
	HookConditionPreHookDB     = "prehookdb"
	HookConditionPostHookDB    = "posthookdb"
	HookConditionPreHookQueue  = "prehookqueue"
	HookConditionPostHookQueue = "posthookqueue"
)

var (
	TagReceiverType  = tag.MustNewKey("receiver_type")
	TagRoutingMethod = tag.MustNewKey("routing_method")
	TagMessageStatus = tag.MustNewKey("status")
	TagHookCondition = tag.MustNewKey("hook_condition")

	MetricNotificationMessageQueueTime = stats.Int64("notification.message.queue.time", "time of message from enqueued to be picked up", stats.UnitMilliseconds)

	MetricNotificationMessageCounter = stats.Int64("notification.message", "notification messages counter", stats.UnitDimensionless)

	MetricNotificationSubscriberNotFound = stats.Int64("notification.subscriber.notfound", "notification does not match any subscription", stats.UnitDimensionless)

	MetricReceiverHookFailed = stats.Int64("receiver.hook.failed", "failed hook condition", stats.UnitDimensionless)
)

func setupApplicationViews() error {
	return view.Register(
		&view.View{
			Name:        MetricNotificationMessageQueueTime.Name(),
			Description: MetricNotificationMessageQueueTime.Description(),
			TagKeys:     []tag.Key{TagReceiverType},
			Measure:     MetricNotificationMessageQueueTime,
			Aggregation: view.Distribution(),
		},
		&view.View{
			Name:        MetricNotificationMessageCounter.Name(),
			Description: MetricNotificationMessageCounter.Description(),
			TagKeys:     []tag.Key{TagReceiverType, TagRoutingMethod, TagMessageStatus},
			Measure:     MetricNotificationMessageCounter,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        MetricNotificationSubscriberNotFound.Name(),
			Description: MetricNotificationSubscriberNotFound.Description(),
			Measure:     MetricNotificationSubscriberNotFound,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        MetricReceiverHookFailed.Name(),
			Description: MetricReceiverHookFailed.Description(),
			TagKeys:     []tag.Key{TagReceiverType, TagHookCondition},
			Measure:     MetricReceiverHookFailed,
			Aggregation: view.Count(),
		},
	)
}
