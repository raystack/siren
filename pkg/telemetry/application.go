package telemetry

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opentelemetry.io/otel/metric/global"
)

var Meter = global.MeterProvider().Meter("siren")

var (
	TagReceiverType  = tag.MustNewKey("receiver_type")
	TagRoutingMethod = tag.MustNewKey("routing_method")
	TagMessageStatus = tag.MustNewKey("status")

	MetricNotificationMessageQueueTime = stats.Int64("notification.message.queue.time", "time of message from enqueued to be picked up", stats.UnitMilliseconds)

	MetricNotificationMessageEnqueue   = stats.Int64("notification.message.enqueue", "enqueued notification messages", stats.UnitDimensionless)
	MetricNotificationMessagePending   = stats.Int64("notification.message.pending", "processed notification messages", stats.UnitDimensionless)
	MetricNotificationMessageFailed    = stats.Int64("notification.message.failed", "failed to publish notification messages", stats.UnitDimensionless)
	MetricNotificationMessagePublished = stats.Int64("notification.message.published", "published notification messages", stats.UnitDimensionless)

	MetricNotificationSubscriberNotFound = stats.Int64("notification.subscriber.notfound", "notification does not match any subscription", stats.UnitDimensionless)

	MetricReceiverPreHookDBFailed     = stats.Int64("receiver.prehookdb.failed", "failed prehook db receiver", stats.UnitDimensionless)
	MetricReceiverPostHookDBFailed    = stats.Int64("receiver.posthookdb.failed", "failed posthook db receiver", stats.UnitDimensionless)
	MetricReceiverPreHookQueueFailed  = stats.Int64("receiver.prehookqueue.failed", "failed prehook queue receiver", stats.UnitDimensionless)
	MetricReceiverPostHookQueueFailed = stats.Int64("receiver.posthookqueue.failed", "failed posthook queue receiver", stats.UnitDimensionless)
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
			Name:        MetricNotificationMessageEnqueue.Name(),
			Description: MetricNotificationMessageEnqueue.Description(),
			TagKeys:     []tag.Key{TagReceiverType, TagRoutingMethod, TagMessageStatus},
			Measure:     MetricNotificationMessageEnqueue,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        MetricNotificationMessagePending.Name(),
			Description: MetricNotificationMessagePending.Description(),
			TagKeys:     []tag.Key{TagReceiverType, TagMessageStatus},
			Measure:     MetricNotificationMessagePending,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        MetricNotificationMessageFailed.Name(),
			Description: MetricNotificationMessageFailed.Description(),
			TagKeys:     []tag.Key{TagReceiverType, TagMessageStatus},
			Measure:     MetricNotificationMessageFailed,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        MetricNotificationMessagePublished.Name(),
			Description: MetricNotificationMessagePublished.Description(),
			TagKeys:     []tag.Key{TagReceiverType, TagMessageStatus},
			Measure:     MetricNotificationMessagePublished,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        MetricNotificationSubscriberNotFound.Name(),
			Description: MetricNotificationSubscriberNotFound.Description(),
			Measure:     MetricNotificationSubscriberNotFound,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        MetricReceiverPreHookDBFailed.Name(),
			Description: MetricReceiverPreHookDBFailed.Description(),
			TagKeys:     []tag.Key{TagReceiverType},
			Measure:     MetricReceiverPreHookDBFailed,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        MetricReceiverPostHookDBFailed.Name(),
			Description: MetricReceiverPostHookDBFailed.Description(),
			TagKeys:     []tag.Key{TagReceiverType},
			Measure:     MetricReceiverPostHookDBFailed,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        MetricReceiverPreHookQueueFailed.Name(),
			Description: MetricReceiverPreHookQueueFailed.Description(),
			TagKeys:     []tag.Key{TagReceiverType, TagRoutingMethod},
			Measure:     MetricReceiverPreHookQueueFailed,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        MetricReceiverPostHookQueueFailed.Name(),
			Description: MetricReceiverPostHookQueueFailed.Description(),
			TagKeys:     []tag.Key{TagReceiverType},
			Measure:     MetricReceiverPostHookQueueFailed,
			Aggregation: view.Sum(),
		},
	)
}
