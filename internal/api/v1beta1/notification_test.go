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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestGRPCServer_NotifyReceiver(t *testing.T) {
	var dummyReq = &sirenv1beta1.NotifyReceiverRequest{
		Id: 1,
		Payload: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"receiver_name": {
					Kind: &structpb.Value_StringValue{
						StringValue: "foo",
					},
				},
				"receiver_type": {
					Kind: &structpb.Value_StringValue{
						StringValue: "channel",
					},
				},
				"message": {
					Kind: &structpb.Value_StringValue{
						StringValue: "bar",
					},
				},
				"blocks": {
					Kind: &structpb.Value_ListValue{
						ListValue: &structpb.ListValue{
							Values: []*structpb.Value{
								{
									Kind: &structpb.Value_StructValue{
										StructValue: &structpb.Struct{
											Fields: map[string]*structpb.Value{
												"type": {
													Kind: &structpb.Value_StringValue{
														StringValue: "section",
													},
												},
												"text": {
													Kind: &structpb.Value_StructValue{
														StructValue: &structpb.Struct{
															Fields: map[string]*structpb.Value{
																"type": {
																	Kind: &structpb.Value_StringValue{
																		StringValue: "mrkdwn",
																	},
																},
																"text": {
																	Kind: &structpb.Value_StringValue{
																		StringValue: "Hello",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	t.Run("should return invalid argument if notify receiver return invalid argument", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		mockNotificationService := new(mocks.NotificationService)

		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), &api.Deps{ReceiverService: mockedReceiverService, NotificationService: mockNotificationService})

		mockNotificationService.EXPECT().Dispatch(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Notification")).Return(nil)
		mockedReceiverService.EXPECT().Notify(mock.AnythingOfType("*context.emptyCtx"),
			uint64(1),
			map[string]interface{}{
				"receiver_name": dummyReq.GetPayload().Fields["receiver_name"].GetStringValue(),
				"receiver_type": dummyReq.GetPayload().Fields["receiver_type"].GetStringValue(),
				"message":       dummyReq.GetPayload().Fields["message"].GetStringValue(),
				"blocks":        dummyReq.GetPayload().Fields["blocks"].GetListValue().AsSlice(),
			},
		).Return(errors.ErrInvalid)
		_, err := dummyGRPCServer.NotifyReceiver(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = request is not valid")
		mockedReceiverService.AssertExpectations(t)
	})

	t.Run("should return internal error if notify receiver return some error", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}

		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), &api.Deps{ReceiverService: mockedReceiverService})

		mockedReceiverService.EXPECT().Notify(
			mock.AnythingOfType("*context.emptyCtx"),
			uint64(1),
			map[string]interface{}{
				"receiver_name": dummyReq.GetPayload().Fields["receiver_name"].GetStringValue(),
				"receiver_type": dummyReq.GetPayload().Fields["receiver_type"].GetStringValue(),
				"message":       dummyReq.GetPayload().Fields["message"].GetStringValue(),
				"blocks":        dummyReq.GetPayload().Fields["blocks"].GetListValue().AsSlice(),
			},
		).Return(errors.New("some error"))
		_, err := dummyGRPCServer.NotifyReceiver(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
		mockedReceiverService.AssertExpectations(t)
	})

	t.Run("should return OK response if notify receiver succeed", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		mockNotificationService := new(mocks.NotificationService)

		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), &api.Deps{ReceiverService: mockedReceiverService, NotificationService: mockNotificationService})

		mockNotificationService.EXPECT().Dispatch(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Notification")).Return(nil)
		mockedReceiverService.EXPECT().Notify(
			mock.AnythingOfType("*context.emptyCtx"),
			uint64(1),
			map[string]interface{}{
				"receiver_name": dummyReq.GetPayload().Fields["receiver_name"].GetStringValue(),
				"receiver_type": dummyReq.GetPayload().Fields["receiver_type"].GetStringValue(),
				"message":       dummyReq.GetPayload().Fields["message"].GetStringValue(),
				"blocks":        dummyReq.GetPayload().Fields["blocks"].GetListValue().AsSlice(),
			},
		).Return(nil)
		_, err := dummyGRPCServer.NotifyReceiver(context.Background(), dummyReq)
		assert.Nil(t, err)
		mockedReceiverService.AssertExpectations(t)
	})

}
