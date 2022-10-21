package notification

import (
	"time"

	"github.com/odpf/siren/plugins/queues"
)

type Config struct {
	Queue          queues.Config `mapstructure:"queue"`
	MessageHandler HandlerConfig `mapstructure:"message_handler"`
	DLQHandler     HandlerConfig `mapstructure:"dlq_handler"`
}

type HandlerConfig struct {
	Enabled       bool          `mapstructure:"enabled"`
	PollDuration  time.Duration `mapstructure:"poll_duration"`
	ReceiverTypes []string      `mapstructure:"receiver_types"`
	BatchSize     int           `mapstructure:"batch_size"`
}
