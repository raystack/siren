package telemetry

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opentelemetry.io/otel/metric/global"
)

var Meter = global.MeterProvider().Meter("siren")

var (
	TagReceiverType = tag.MustNewKey("receiver_type")

	MetricPublishedNotifications = stats.Int64("published_notifications", "published notifications", stats.UnitDimensionless)
	MetricEnqueueMQ              = stats.Int64("enqueue_message_queue", "enqueued messages to message queue", stats.UnitDimensionless)
	MetricDequeueMQ              = stats.Int64("dequeue_message_queue", "dequeued messages from message queue", stats.UnitDimensionless)
	MetricEnqueueDLQ             = stats.Int64("enqueue_dlq", "enqueued messages to dlq", stats.UnitDimensionless)
	MetricDequeueDLQ             = stats.Int64("dequeue_dlq", "dequeued messages from dlq", stats.UnitDimensionless)
)

func setupApplicationViews() error {
	return view.Register(
		&view.View{
			Name:        MetricPublishedNotifications.Name(),
			Description: MetricPublishedNotifications.Description(),
			TagKeys:     []tag.Key{TagReceiverType},
			Measure:     MetricPublishedNotifications,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        MetricEnqueueMQ.Name(),
			Description: MetricEnqueueMQ.Description(),
			TagKeys:     []tag.Key{TagReceiverType},
			Measure:     MetricEnqueueMQ,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        MetricDequeueMQ.Name(),
			Description: MetricDequeueMQ.Description(),
			TagKeys:     []tag.Key{TagReceiverType},
			Measure:     MetricDequeueMQ,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        MetricEnqueueDLQ.Name(),
			Description: MetricEnqueueDLQ.Description(),
			TagKeys:     []tag.Key{TagReceiverType},
			Measure:     MetricEnqueueDLQ,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        MetricDequeueDLQ.Name(),
			Description: MetricDequeueMQ.Description(),
			TagKeys:     []tag.Key{TagReceiverType},
			Measure:     MetricDequeueMQ,
			Aggregation: view.Sum(),
		},
	)
}
