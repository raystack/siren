package notification

import (
	"time"

	"github.com/google/uuid"
)

type MessageStatus string

const (
	DefaultMaxTries = 3

	MessageStatusEnqueued  MessageStatus = "enqueued"
	MessageStatusFailed    MessageStatus = "failed"
	MessageStatusPending   MessageStatus = "pending"
	MessageStatusPublished MessageStatus = "published"
)

type MessageOption func(*Message)

func InitWithCreateTime(timeNow time.Time) MessageOption {
	return func(m *Message) {
		m.CreatedAt = timeNow
		m.UpdatedAt = timeNow
	}
}

func InitWithID(id string) MessageOption {
	return func(m *Message) {
		m.ID = id
	}
}

func InitWithExpiryDuration(dur time.Duration) MessageOption {
	return func(m *Message) {
		m.expiryDuration = dur
	}
}

func InitWithMaxTries(mt int) MessageOption {
	return func(m *Message) {
		m.MaxTries = mt
	}
}

type Message struct {
	ID     string
	Status MessageStatus

	ReceiverType string
	Configs      map[string]interface{} // the datasource to build vendor-specific configs
	Detail       map[string]interface{} // the datasource to build vendor-specific message
	LastError    string

	MaxTries int
	TryCount int

	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time

	expiryDuration time.Duration
}

func (m *Message) Initialize(
	n Notification,
	receiverType string,
	notificationConfigs map[string]interface{},
	opts ...MessageOption,
) {
	var timeNow = time.Now()

	m.ID = uuid.NewString()
	m.Status = MessageStatusEnqueued

	m.ReceiverType = receiverType
	m.Configs = notificationConfigs
	detail := make(map[string]interface{})
	for k, v := range n.Labels {
		detail[k] = v
	}
	for k, v := range n.Variables {
		detail[k] = v
	}
	m.Detail = detail

	m.MaxTries = DefaultMaxTries

	m.CreatedAt = timeNow
	m.UpdatedAt = timeNow

	for _, opt := range opts {
		opt(m)
	}

	if m.expiryDuration != 0 {
		m.ExpiredAt = m.CreatedAt.Add(m.expiryDuration)
	}
}

func (m *Message) MarkFailed(updatedAt time.Time, err error) {
	m.TryCount = m.TryCount + 1
	m.LastError = err.Error()
	m.Status = MessageStatusFailed
	m.UpdatedAt = updatedAt
}

func (m *Message) MarkPending(updatedAt time.Time) {
	m.Status = MessageStatusPending
	m.UpdatedAt = updatedAt
}

func (m *Message) MarkPublished(updatedAt time.Time) {
	m.TryCount = m.TryCount + 1
	m.Status = MessageStatusPublished
	m.UpdatedAt = updatedAt
}
