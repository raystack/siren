package notification

import (
	"time"

	"github.com/goto/siren/plugins/queues"
)

type Config struct {
	Queue          queues.Config `mapstructure:"queue" yaml:"queue"`
	MessageHandler HandlerConfig `mapstructure:"message_handler" yaml:"message_handler"`
	DLQHandler     HandlerConfig `mapstructure:"dlq_handler" yaml:"dlq_handler"`
}

type HandlerConfig struct {
	Enabled       bool          `mapstructure:"enabled" yaml:"enabled" default:"true"`
	PollDuration  time.Duration `mapstructure:"poll_duration" yaml:"poll_duration" default:"5s"`
	ReceiverTypes []string      `mapstructure:"receiver_types" yaml:"receiver_types"`
	BatchSize     int           `mapstructure:"batch_size" yaml:"batch_size" default:"1"`
}
