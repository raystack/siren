package notification_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/odpf/siren/core/notification"
)

func TestMessage_Initialize(t *testing.T) {
	var (
		testID             = "some-id"
		testTimeNow        = time.Now()
		testExpiryDuration = 5 * time.Minute
	)
	testCases := []struct {
		name                string
		n                   notification.Notification
		receiverType        string
		notificationConfigs map[string]interface{}
		want                *notification.Message
		wantErr             bool
	}{
		{
			name: "all notification labels and variables should be merged to message detail and variable takes precedence if key conflict",
			n: notification.Notification{
				Labels: map[string]string{
					"labelkey1": "value1",
					"samekey":   "label_value",
				},
				Variables: map[string]interface{}{
					"varkey1": "value1",
					"samekey": "var_value",
				},
			},
			want: &notification.Message{
				ID:     testID,
				Status: notification.MessageStatusEnqueued,
				Detail: map[string]interface{}{
					"labelkey1": "value1",
					"varkey1":   "value1",
					"samekey":   "var_value",
				},
				MaxTries:  notification.DefaultMaxTries,
				CreatedAt: testTimeNow,
				UpdatedAt: testTimeNow,
				ExpiredAt: testTimeNow.Add(testExpiryDuration),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &notification.Message{}
			m.Initialize(tc.n, tc.receiverType, tc.notificationConfigs,
				notification.InitWithID(testID),
				notification.InitWithCreateTime(testTimeNow),
				notification.InitWithExpiryDuration(testExpiryDuration),
			)

			if diff := cmp.Diff(m, tc.want, cmpopts.IgnoreUnexported(notification.Message{})); diff != "" {
				t.Errorf("Notification.ToMessage() diff = %v", diff)
			}
		})
	}
}

func TestMessage_Mark(t *testing.T) {
	var (
		createTime     = time.Now()
		expiryDuration = 1 * time.Minute
		expiredAt      = createTime.Add(expiryDuration)
	)
	m := &notification.Message{
		ID:     "some-id",
		Status: notification.MessageStatusEnqueued,
		Detail: map[string]interface{}{
			"labelkey1": "value1",
			"varkey1":   "value1",
		},
		MaxTries:  notification.DefaultMaxTries,
		CreatedAt: createTime,
		UpdatedAt: createTime,
		ExpiredAt: expiredAt,
	}

	t.Run("mark failed should updates message to the failed state", func(t *testing.T) {
		var (
			testTimeNow     = time.Now()
			err             = errors.New("some error")
			expectedMessage = m
		)

		expectedMessage.TryCount = m.TryCount + 1
		expectedMessage.LastError = err.Error()
		expectedMessage.Status = notification.MessageStatusFailed
		expectedMessage.UpdatedAt = testTimeNow

		m.MarkFailed(testTimeNow, true, err)

		if diff := cmp.Diff(m, expectedMessage, cmpopts.IgnoreUnexported(notification.Message{})); diff != "" {
			t.Errorf("result not match, diff = %v", diff)
		}
	})
	t.Run("mark pending should updates message to the pending state", func(t *testing.T) {
		var (
			testTimeNow     = time.Now()
			expectedMessage = m
		)

		expectedMessage.TryCount = m.TryCount + 1
		expectedMessage.Status = notification.MessageStatusPending
		expectedMessage.UpdatedAt = testTimeNow

		m.MarkPending(testTimeNow)

		if diff := cmp.Diff(m, expectedMessage, cmpopts.IgnoreUnexported(notification.Message{})); diff != "" {
			t.Errorf("result not match, diff = %v", diff)
		}
	})
	t.Run("mark published should updates message to the published state", func(t *testing.T) {
		var (
			testTimeNow     = time.Now()
			expectedMessage = m
		)

		expectedMessage.TryCount = m.TryCount + 1
		expectedMessage.Status = notification.MessageStatusPublished
		expectedMessage.UpdatedAt = testTimeNow

		m.MarkPublished(testTimeNow)

		if diff := cmp.Diff(m, expectedMessage, cmpopts.IgnoreUnexported(notification.Message{})); diff != "" {
			t.Errorf("result not match, diff = %v", diff)
		}
	})
}
