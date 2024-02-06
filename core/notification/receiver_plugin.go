package notification

import (
	"context"

	"github.com/goto/siren/plugins/queues"
)

type Notifier interface {
	PreHookQueueTransformConfigs(ctx context.Context, notificationConfigMap map[string]any) (map[string]any, error)
	PostHookQueueTransformConfigs(ctx context.Context, notificationConfigMap map[string]any) (map[string]any, error)
	GetSystemDefaultTemplate() string
	Send(ctx context.Context, message Message) (bool, error)
}

type Queuer interface {
	Enqueue(ctx context.Context, ms ...Message) error
	Dequeue(ctx context.Context, receiverTypes []string, batchSize int, handlerFn func(context.Context, []Message) error) error
	SuccessCallback(ctx context.Context, ms Message) error
	ErrorCallback(ctx context.Context, ms Message) error
	Cleanup(ctx context.Context, filter queues.FilterCleanup) error
	Stop(ctx context.Context) error
}
