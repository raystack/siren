package v1beta1_test

import (
	"context"
	"testing"

	"github.com/odpf/salt/log"
	"github.com/odpf/siren/internal/api"
	"github.com/odpf/siren/internal/api/mocks"
	"github.com/odpf/siren/internal/api/v1beta1"
	"github.com/odpf/siren/pkg/errors"
	sirenv1beta1 "github.com/odpf/siren/proto/odpf/siren/v1beta1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
)

func TestGRPCServer_NotifyReceiver(t *testing.T) {
	const idempotencyHeaderKey = "idempotency-key"
	testCases := []struct {
		name           string
		idempotencyKey string
		setup          func(*mocks.NotificationService)
		errString      string
	}{
		{
			name:           "should return invalid argument if notify receiver return invalid argument",
			idempotencyKey: "test",
			setup: func(ns *mocks.NotificationService) {
				ns.EXPECT().CheckAndInsertIdempotency(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(1, nil)
				ns.EXPECT().DispatchToReceiver(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("notification.Notification"), mock.AnythingOfType("uint64")).Return(errors.ErrInvalid)
			},
			errString: "rpc error: code = InvalidArgument desc = request is not valid",
		},
		{
			name:           "should return internal error if notify receiver return some error",
			idempotencyKey: "test",
			setup: func(ns *mocks.NotificationService) {
				ns.EXPECT().CheckAndInsertIdempotency(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(1, nil)
				ns.EXPECT().DispatchToReceiver(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("notification.Notification"), mock.AnythingOfType("uint64")).Return(errors.New("some error"))
			},
			errString: "rpc error: code = Internal desc = some unexpected error occurred",
		},
		{
			name:           "should return success if request is idempotent",
			idempotencyKey: "test",
			setup: func(ns *mocks.NotificationService) {
				ns.EXPECT().CheckAndInsertIdempotency(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(0, errors.ErrConflict)
			},
		},
		{
			name:           "should return error if idempotency checking return error",
			idempotencyKey: "test",
			setup: func(ns *mocks.NotificationService) {
				ns.EXPECT().CheckAndInsertIdempotency(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(0, errors.New("some error"))
			},
			errString: "rpc error: code = Internal desc = some unexpected error occurred",
		},
		{
			name:           "should return error if error updating idempotency as success",
			idempotencyKey: "test",
			setup: func(ns *mocks.NotificationService) {
				ns.EXPECT().CheckAndInsertIdempotency(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(1, nil)
				ns.EXPECT().MarkIdempotencyAsSuccess(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("uint64")).Return(errors.New("some error"))
				ns.EXPECT().DispatchToReceiver(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("notification.Notification"), mock.AnythingOfType("uint64")).Return(nil)
			},
			errString: "rpc error: code = Internal desc = some unexpected error occurred",
		},
		{
			name:           "should return OK response if notify receiver succeed",
			idempotencyKey: "test",
			setup: func(ns *mocks.NotificationService) {
				ns.EXPECT().CheckAndInsertIdempotency(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(1, nil)
				ns.EXPECT().MarkIdempotencyAsSuccess(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("uint64")).Return(nil)
				ns.EXPECT().DispatchToReceiver(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("notification.Notification"), mock.AnythingOfType("uint64")).Return(nil)
			},
		},
		{
			name: "should return OK response if notify receiver succeed without idempotency",
			setup: func(ns *mocks.NotificationService) {
				ns.EXPECT().DispatchToReceiver(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("notification.Notification"), mock.AnythingOfType("uint64")).Return(nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				mockNotificationService = new(mocks.NotificationService)
			)

			if tc.setup != nil {
				tc.setup(mockNotificationService)
			}

			dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{
				IdempotencyKey: idempotencyHeaderKey,
			}, &api.Deps{NotificationService: mockNotificationService})
			ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
				idempotencyHeaderKey: tc.idempotencyKey,
			}))
			_, err := dummyGRPCServer.NotifyReceiver(ctx, &sirenv1beta1.NotifyReceiverRequest{})

			if (err != nil) && tc.errString != err.Error() {
				t.Errorf("NotifyReceiver() error = %v, wantErr %v", err, tc.errString)
			}

			mockNotificationService.AssertExpectations(t)
		})
	}
}
