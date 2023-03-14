package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/goto/siren/pkg/retry"
	"github.com/stretchr/testify/assert"
)

type testError struct {
	numExecution int
	isRetryable  bool
}

func (te *testError) Execute(ctx context.Context) error {
	te.numExecution = te.numExecution + 1
	if te.isRetryable {
		return retry.RetryableError{Err: errors.New("retryable error")}
	}
	return errors.New("some error")
}

func TestRetrier_Error(t *testing.T) {
	var te *testError
	testCases := []struct {
		name                 string
		f                    func() func(ctx context.Context) error
		cfg                  retry.Config
		expectedNumExecution int
	}{
		{
			name: "execution should be retried with the same number of times if retryable",
			f: func() func(ctx context.Context) error {
				te = &testError{numExecution: 0, isRetryable: true}
				return te.Execute
			},
			cfg: retry.Config{
				MaxTries: 4,
				Enable:   true,
			},
			expectedNumExecution: 5,
		},
		{
			name: "execution should not be retried if error not retryable",
			f: func() func(ctx context.Context) error {
				te = &testError{numExecution: 0, isRetryable: false}
				return te.Execute
			},
			cfg: retry.Config{
				MaxTries: 4,
				Enable:   true,
			},
			expectedNumExecution: 1,
		},
		{
			name: "execution should not be retried if retrier is disabled",
			f: func() func(ctx context.Context) error {
				te = &testError{numExecution: 0, isRetryable: false}
				return te.Execute
			},
			cfg: retry.Config{
				MaxTries: 4,
			},
			expectedNumExecution: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rtr := retry.New(tc.cfg)
			_ = rtr.Run(context.TODO(), tc.f())

			assert.Equal(t, tc.expectedNumExecution, te.numExecution)
		})
	}
}

type testTime struct {
	prevExecution time.Time
	waitTimes     []time.Duration
}

func (tr *testTime) Execute(ctx context.Context) error {
	now := time.Now()

	if !tr.prevExecution.IsZero() {
		durationSince := now.Sub(tr.prevExecution)
		tr.waitTimes = append(tr.waitTimes, durationSince.Round(time.Millisecond))
	}
	tr.prevExecution = now

	return errors.New("some error")
}

func TestRetrier_RetryWithoutBackoff(t *testing.T) {
	testCases := []struct {
		name              string
		cfg               retry.Config
		expectedWaitTimes []time.Duration
	}{
		{
			name: "a retry executions without backoff should be at constant rate. (10ms, 4 retries)",
			cfg: retry.Config{
				MaxTries:      4,
				EnableBackoff: true,
				WaitDuration:  10 * time.Millisecond,
				Enable:        true,
			},
			expectedWaitTimes: []time.Duration{
				10 * time.Millisecond,
				10 * time.Millisecond,
				10 * time.Millisecond,
				10 * time.Millisecond,
			},
		},
		{
			name: "a retry executions without backoff should be at constant rate (30ms, 2 retries)",
			cfg: retry.Config{
				MaxTries:      2,
				EnableBackoff: false,
				WaitDuration:  30 * time.Millisecond,
				Enable:        true,
			},
			expectedWaitTimes: []time.Duration{
				30 * time.Millisecond,
				30 * time.Millisecond,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tr := &testTime{}
			rtr := retry.New(tc.cfg)
			_ = rtr.Run(context.TODO(), tr.Execute)

			assert.InEpsilonSlice(t, tc.expectedWaitTimes, tr.waitTimes, 0.1)
		})
	}
}
