package queues

type Kind string

const (
	KindInMemory Kind = "inmemory"
	KindPostgres Kind = "postgres"
)

func (k Kind) String() string {
	return string(k)
}

type Config struct {
	Kind Kind `mapstructure:"kind" default:"inmemory"`
}

type FilterCleanup struct {
	MessagePendingTimeThreshold   string
	MessagePublishedTimeThreshold string
}
