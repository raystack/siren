package notification

import (
	"context"
	"time"

	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/plugins/queues"
)

//go:generate mockery --name=Notifier -r --case underscore --with-expecter --structname Notifier --filename notifier.go --output=./mocks
type Notifier interface {
	PreHookTransformConfigs(ctx context.Context, notificationConfigMap map[string]interface{}) (map[string]interface{}, error)
	PostHookTransformConfigs(ctx context.Context, notificationConfigMap map[string]interface{}) (map[string]interface{}, error)
	DefaultTemplateOfProvider(templateName string) string
	Publish(ctx context.Context, message Message) (bool, error)
}

//go:generate mockery --name=Queuer -r --case underscore --with-expecter --structname Queuer --filename queuer.go --output=./mocks
type Queuer interface {
	Enqueue(ctx context.Context, ms ...Message) error
	Dequeue(ctx context.Context, receiverTypes []string, batchSize int, handlerFn func(context.Context, []Message) error) error
	SuccessCallback(ctx context.Context, ms Message) error
	ErrorCallback(ctx context.Context, ms Message) error
	Cleanup(ctx context.Context, filter queues.FilterCleanup) error
	Stop(ctx context.Context) error
}

// Notification is a model of notification
type Notification struct {
	ID                  string
	Data                map[string]interface{} `json:"data"`
	Labels              map[string]string      `json:"labels"`
	ValidDurationString string                 `json:"valid_duration"`
	Template            string                 `json:"template"`
	CreatedAt           time.Time
}

// ToMessage transforms Notification model to one or several Messages
func (n Notification) ToMessage(receiverType string, notificationConfigMap map[string]interface{}) (*Message, error) {
	var (
		expiryDuration time.Duration
		err            error
	)

	if n.ValidDurationString != "" {
		expiryDuration, err = time.ParseDuration(n.ValidDurationString)
		if err != nil {
			return nil, errors.ErrInvalid.WithMsgf(err.Error())
		}
	}

	nm := &Message{}
	nm.Initialize(
		n,
		receiverType,
		notificationConfigMap,
		InitWithExpiryDuration(expiryDuration),
	)

	return nm, nil
}
