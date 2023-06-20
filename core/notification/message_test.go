package notification_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/core/notification/mocks"
	"github.com/stretchr/testify/mock"
)

func TestMessage_InitMessage(t *testing.T) {
	var (
		testID             = "some-id"
		testTimeNow        = time.Now()
		testExpiryDuration = 5 * time.Minute
	)
	testCases := []struct {
		name                string
		setup               func(*mocks.Notifier)
		n                   notification.Notification
		receiverType        string
		notificationConfigs map[string]any
		want                notification.Message
		errString           string
	}{
		{
			name: "all notification labels and data should be merged to message detail and data takes precedence if key conflict",
			setup: func(n *mocks.Notifier) {
				n.EXPECT().PreHookQueueTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("map[string]interface {}")).Return(nil, nil)
			},
			n: notification.Notification{
				Type: notification.TypeSubscriber,
				Labels: map[string]string{
					"labelkey1": "value1",
					"samekey":   "label_value",
				},
				Data: map[string]any{
					"varkey1": "value1",
					"samekey": "var_value",
				},
			},
			want: notification.Message{
				ID:     testID,
				Status: notification.MessageStatusEnqueued,
				Details: map[string]any{
					"labelkey1":                             "value1",
					"varkey1":                               "value1",
					"samekey":                               "var_value",
					notification.DetailsKeyNotificationType: notification.TypeSubscriber,
				},
				CreatedAt: testTimeNow,
				UpdatedAt: testTimeNow,
				ExpiredAt: testTimeNow.Add(testExpiryDuration),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockNotifierPlugin := new(mocks.Notifier)

			if tc.setup != nil {
				tc.setup(mockNotifierPlugin)
			}

			m, err := notification.InitMessage(
				context.TODO(),
				mockNotifierPlugin,
				tc.n,
				tc.receiverType,
				tc.notificationConfigs,
				notification.InitWithID(testID),
				notification.InitWithCreateTime(testTimeNow),
				notification.InitWithExpiryDuration(testExpiryDuration),
			)
			if err != nil {
				if err.Error() != tc.errString {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.errString)
				}
			}

			if diff := cmp.Diff(m, tc.want,
				cmpopts.IgnoreUnexported(notification.Message{}),
				cmpopts.IgnoreFields(notification.Message{}, "MaxTries")); diff != "" {
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
		Details: map[string]any{
			"labelkey1": "value1",
			"varkey1":   "value1",
		},
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

		if diff := cmp.Diff(m, expectedMessage,
			cmpopts.IgnoreUnexported(notification.Message{}),
			cmpopts.IgnoreFields(notification.Message{}, "MaxTries")); diff != "" {
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

		if diff := cmp.Diff(m, expectedMessage,
			cmpopts.IgnoreUnexported(notification.Message{}),
			cmpopts.IgnoreFields(notification.Message{}, "MaxTries")); diff != "" {
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

		if diff := cmp.Diff(m, expectedMessage,
			cmpopts.IgnoreUnexported(notification.Message{}),
			cmpopts.IgnoreFields(notification.Message{}, "MaxTries")); diff != "" {
			t.Errorf("result not match, diff = %v", diff)
		}
	})
}
