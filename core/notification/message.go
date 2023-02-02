package notification

import (
	"time"

	"github.com/google/uuid"
)

// MessageStatus determines the state of the message
type MessageStatus string

const (
	DefaultMaxTries = 3

	// additional details
	DetailsKeyRoutingMethod = "routing_method"

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
	ID     string
	Status MessageStatus

	ReceiverType string
	Configs      map[string]interface{} // the datasource to build vendor-specific configs
	Details      map[string]interface{} // the datasource to build vendor-specific message
	LastError    string

	MaxTries  int
	TryCount  int
	Retryable bool

	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time

	expiryDuration time.Duration
}

// Initialize initializes the message with some default value
// or the customized value
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

	details := make(map[string]interface{})
	for k, v := range n.Labels {
		details[k] = v
	}
	for k, v := range n.Data {
		details[k] = v
	}

	m.Details = details

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

// AddDetail adds a custom kv string detail
func (m *Message) AddStringDetail(key, value string) {
	m.Details[key] = value
}
