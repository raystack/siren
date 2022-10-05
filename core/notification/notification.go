package notification

import (
	"context"
	"time"

	"github.com/odpf/siren/core/subscription"
)

//go:generate mockery --name=Notifier -r --case underscore --with-expecter --structname Notifier --filename notifier.go --output=./mocks
type Notifier interface {
	ValidateConfigMap(notificationConfigMap map[string]interface{}) error
	Publish(ctx context.Context, message Message) (bool, error)
}

//go:generate mockery --name=SubscriptionService -r --case underscore --with-expecter --structname SubscriptionService --filename subscription_service.go --output=./mocks
type SubscriptionService interface {
	MatchByLabels(ctx context.Context, labels map[string]string) ([]subscription.Subscription, error)
}

//go:generate mockery --name=Queuer -r --case underscore --with-expecter --structname Queuer --filename queuer.go --output=./mocks
type Queuer interface {
	Enqueue(ctx context.Context, ms ...Message) error
	Dequeue(ctx context.Context, receiverTypes []string, batchSize int, handlerFn func(context.Context, []Message) error) error
	SuccessHandler(ctx context.Context, ms Message) error
	ErrorHandler(ctx context.Context, ms Message) error
}

// Notification is a model of notification
type Notification struct {
	ID                  string                 `json:"id"`
	Variables           map[string]interface{} `json:"variables"`
	Labels              map[string]string      `json:"labels"`
	ValidDurationString string                 `json:"valid_duration"`
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
			return nil, err
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
