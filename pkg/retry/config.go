package retry

import "time"

type Config struct {
	// duration to wait on the retries.
	WaitDuration time.Duration `mapstructure:"wait_duration" yaml:"wait_duration" default:"20ms"`
	// enable exponential backoff on the retry and jitter.
	EnableBackoff bool `mapstructure:"enable_backoff" yaml:"enable_backoff" default:"false"`
	// number of times that will be retried in case of error
	// before returning the error itself.
	MaxTries int `mapstructure:"max_retry" yaml:"max_tries" default:"3"`
	// short circuit the retrier if false
	Enable bool `mapstructure:"enable" yaml:"enable" default:"true"`
}
