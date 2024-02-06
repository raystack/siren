package notification

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/goto/siren/core/template"
	"github.com/goto/siren/pkg/errors"
	"gopkg.in/yaml.v3"
)

// MessageStatus determines the state of the message
type MessageStatus string

const (
	defaultMaxTries int = 3

	// additional details
	DetailsKeyNotificationType = "notification_type"

	MessageStatusEnqueued  MessageStatus = "enqueued"
	MessageStatusFailed    MessageStatus = "failed"
	MessageStatusPending   MessageStatus = "pending"
	MessageStatusPublished MessageStatus = "published"
)

func (ms MessageStatus) String() string {
	return string(ms)
}

// MessageOption provides ability to configure the message initialization
type MessageOption func(*Message)

// InitWithCreateTime initializes the message with custom create time
func InitWithCreateTime(timeNow time.Time) MessageOption {
	return func(m *Message) {
		m.CreatedAt = timeNow
		m.UpdatedAt = timeNow
	}
}

// InitWithID initializes the message with some ID
func InitWithID(id string) MessageOption {
	return func(m *Message) {
		m.ID = id
	}
}

// InitWithExpiryDuration initializes the message with the specified expiry duration
func InitWithExpiryDuration(dur time.Duration) MessageOption {
	return func(m *Message) {
		m.expiryDuration = dur
	}
}

// InitWithMaxTries initializes the message with custom max tries
func InitWithMaxTries(mt int) MessageOption {
	return func(m *Message) {
		m.MaxTries = mt
	}
}

// Message is the model to be sent for a specific subscription's receiver
type Message struct {
	ID           string
	Status       MessageStatus
	ReceiverType string
	Configs      map[string]any // the datasource to build vendor-specific configs
	Details      map[string]any // the datasource to build vendor-specific message
	MaxTries     int
	ExpiredAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time

	LastError string
	TryCount  int
	Retryable bool

	expiryDuration time.Duration
}

// Initialize initializes the message with some default value
// or the customized value
func InitMessage(
	ctx context.Context,
	notifierPlugin Notifier,
	n Notification,
	receiverType string,
	messageConfig map[string]any,
	opts ...MessageOption,
) (Message, error) {
	if notifierPlugin == nil {
		return Message{}, errors.New("notifierPlugin cannot be nil")
	}

	newConfigs, err := notifierPlugin.PreHookQueueTransformConfigs(ctx, messageConfig)
	if err != nil {
		return Message{}, err
	}

	var (
		timeNow = time.Now()
		details = make(map[string]any)
	)

	for k, v := range n.Labels {
		details[k] = v
	}
	for k, v := range n.Data {
		details[k] = v
	}

	m := &Message{
		ID:           uuid.NewString(),
		Status:       MessageStatusEnqueued,
		ReceiverType: receiverType,
		Configs:      newConfigs,
		Details:      details,
		MaxTries:     defaultMaxTries,
		CreatedAt:    timeNow,
		UpdatedAt:    timeNow,
	}

	for _, opt := range opts {
		opt(m)
	}

	if m.expiryDuration != 0 {
		m.ExpiredAt = m.CreatedAt.Add(m.expiryDuration)
	}

	//TODO fetch template if any, if not exist, check provider type, if exist use the default template, if not pass as-is
	// if there is template, render and replace detail with the new one
	if n.Template != "" {
		var templateBody string

		if template.IsReservedName(n.Template) {
			templateBody = notifierPlugin.GetSystemDefaultTemplate()
		}

		if templateBody != "" {
			renderedDetailString, err := template.RenderBody(templateBody, n)
			if err != nil {
				return Message{}, errors.ErrInvalid.WithMsgf("failed to render template receiver %s: %s", receiverType, err.Error())
			}

			var messageDetails map[string]any
			if err := yaml.Unmarshal([]byte(renderedDetailString), &messageDetails); err != nil {
				return Message{}, errors.ErrInvalid.WithMsgf("failed to unmarshal rendered template receiver %s: %s, rendered template: %v", receiverType, err.Error(), renderedDetailString)
			}
			m.Details = messageDetails
		}
	}

	m.Details[DetailsKeyNotificationType] = n.Type

	return *m, nil
}

// MarkFailed update message to the failed state
func (m *Message) MarkFailed(updatedAt time.Time, retryable bool, err error) {
	m.TryCount = m.TryCount + 1
	m.LastError = err.Error()
	m.Retryable = retryable
	m.Status = MessageStatusFailed
	m.UpdatedAt = updatedAt
}

// MarkPending update message to the pending state
func (m *Message) MarkPending(updatedAt time.Time) {
	m.Status = MessageStatusPending
	m.UpdatedAt = updatedAt
}

// MarkPublished update message to the published state
func (m *Message) MarkPublished(updatedAt time.Time) {
	m.TryCount = m.TryCount + 1
	m.Status = MessageStatusPublished
	m.UpdatedAt = updatedAt
}
