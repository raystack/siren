package retry

type RetryableError struct {
	Err error
}

func (rt RetryableError) Error() string {
	return rt.Err.Error()
}
