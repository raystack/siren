package postgresq

type QueueOption func(*Queue)

func WithStrategy(s Strategy) QueueOption {
	return func(q *Queue) {
		q.strategy = s
	}
}
