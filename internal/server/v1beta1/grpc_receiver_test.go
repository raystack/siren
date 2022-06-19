package v1beta1

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/receiver"
	sirenv1beta1 "github.com/odpf/siren/internal/server/proto/odpf/siren/v1beta1"
	"github.com/odpf/siren/internal/server/v1beta1/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestGRPCServer_ListReceiver(t *testing.T) {
	configurations := make(map[string]interface{})
	configurations["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"
	dummyResult := []*receiver.Receiver{
		{
			ID:             1,
			Name:           "foo",
			Type:           "bar",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	t.Run("should return list of all receiver", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			receiverService: mockedReceiverService,
			logger:          log.NewNoop(),
		}
		mockedReceiverService.EXPECT().ListReceivers().
			Return(dummyResult, nil).Once()

		res, err := dummyGRPCServer.ListReceivers(context.Background(), &emptypb.Empty{})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res.GetData()))
		assert.Equal(t, uint64(1), res.GetData()[0].GetId())
		assert.Equal(t, "foo", res.GetData()[0].GetName())
		assert.Equal(t, "bar", res.GetData()[0].GetType())
		assert.Equal(t, "bar", res.GetData()[0].GetConfigurations().AsMap()["foo"])
		assert.Equal(t, "bar", res.GetData()[0].GetLabels()["foo"])
	})

	t.Run("should return error Internal if getting providers failed", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			receiverService: mockedReceiverService,
			logger:          log.NewNoop(),
		}
		mockedReceiverService.EXPECT().ListReceivers().
			Return(nil, errors.New("random error"))

		res, err := dummyGRPCServer.ListReceivers(context.Background(), &emptypb.Empty{})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})

	t.Run("should return error Internal if NewStruct conversion failed", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			receiverService: mockedReceiverService,
			logger:          log.NewNoop(),
		}
		configurations["foo"] = string([]byte{0xff})
		dummyResult := []*receiver.Receiver{
			{
				ID:             1,
				Name:           "foo",
				Type:           "bar",
				Labels:         labels,
				Configurations: configurations,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
		}

		mockedReceiverService.EXPECT().ListReceivers().
			Return(dummyResult, nil)
		res, err := dummyGRPCServer.ListReceivers(context.Background(), &emptypb.Empty{})
		assert.Nil(t, res)
		assert.Equal(t, strings.Replace(err.Error(), "\u00a0", " ", -1),
			"rpc error: code = Internal desc = proto: invalid UTF-8 in string: \"\\xff\"")
	})
}

func TestGRPCServer_CreateReceiver(t *testing.T) {
	configurations := make(map[string]interface{})
	configurations["client_id"] = "foo"
	configurations["client_secret"] = "bar"
	configurations["auth_code"] = "foo"
	labels := make(map[string]string)
	labels["foo"] = "bar"
	generatedID := uint64(77)

	configurationsData, _ := structpb.NewStruct(configurations)
	dummyReq := &sirenv1beta1.CreateReceiverRequest{
		Name:           "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurationsData,
	}
	payload := &receiver.Receiver{
		Name:           "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurations,
	}

	t.Run("Should create a receiver object", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			receiverService: mockedReceiverService,
			logger:          log.NewNoop(),
		}
		mockedReceiverService.EXPECT().CreateReceiver(payload).Run(func(rcv *receiver.Receiver) {
			rcv.ID = generatedID
		}).Return(nil).Once()

		res, err := dummyGRPCServer.CreateReceiver(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, generatedID, res.GetId())
	})

	t.Run("should return error Invalid Argument if create receiver failed with err invalid", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			receiverService: mockedReceiverService,
			logger:          log.NewNoop(),
		}

		mockedReceiverService.EXPECT().CreateReceiver(payload).Return(fmt.Errorf("some error: %w", receiver.ErrInvalid)).Once()

		res, err := dummyGRPCServer.CreateReceiver(context.Background(), dummyReq)
		assert.EqualError(t, err,
			"rpc error: code = InvalidArgument desc = some error: bad_request")
		assert.Nil(t, res)
	})

	t.Run("should return error Internal if creating receiver failed", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			receiverService: mockedReceiverService,
			logger:          log.NewNoop(),
		}
		mockedReceiverService.EXPECT().CreateReceiver(payload).
			Return(errors.New("random error")).Once()

		res, err := dummyGRPCServer.CreateReceiver(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})

}

func TestGRPCServer_GetReceiver(t *testing.T) {
	configurations := make(map[string]interface{})
	configurations["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	receiverId := uint64(1)
	dummyReq := &sirenv1beta1.GetReceiverRequest{
		Id: 1,
	}
	payload := &receiver.Receiver{
		Name:           "foo",
		Type:           "bar",
		Labels:         labels,
		Configurations: configurations,
	}

	t.Run("should return a receiver", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			receiverService: mockedReceiverService,
			logger:          log.NewNoop(),
		}
		mockedReceiverService.EXPECT().GetReceiver(receiverId).
			Return(payload, nil).Once()

		res, err := dummyGRPCServer.GetReceiver(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, "foo", res.GetData().GetName())
		assert.Equal(t, "bar", res.GetData().GetType())
		assert.Equal(t, "bar", res.GetData().GetLabels()["foo"])
		assert.Equal(t, "bar", res.GetData().GetConfigurations().AsMap()["foo"])
	})

	t.Run("should return error code 5 if no receiver found", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			receiverService: mockedReceiverService,
			logger:          log.NewNoop(),
		}
		mockedReceiverService.EXPECT().GetReceiver(receiverId).
			Return(nil, nil).Once()

		res, err := dummyGRPCServer.GetReceiver(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = NotFound desc = receiver not found")
	})

	t.Run("should return error Internal if getting receiver failed", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			receiverService: mockedReceiverService,
			logger:          log.NewNoop(),
		}
		mockedReceiverService.EXPECT().GetReceiver(receiverId).
			Return(payload, errors.New("random error")).Once()

		res, err := dummyGRPCServer.GetReceiver(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})

	t.Run("should return error Internal if NewStruct conversion of configuration failed", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			receiverService: mockedReceiverService,
			logger:          log.NewNoop(),
		}

		configurations["foo"] = string([]byte{0xff})
		payload := &receiver.Receiver{
			Name:           "foo",
			Type:           "bar",
			Labels:         labels,
			Configurations: configurations,
		}

		mockedReceiverService.EXPECT().GetReceiver(receiverId).
			Return(payload, nil)
		res, err := dummyGRPCServer.GetReceiver(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.Equal(t, strings.Replace(err.Error(), "\u00a0", " ", -1),
			"rpc error: code = Internal desc = proto: invalid UTF-8 in string: \"\\xff\"")
	})

	t.Run("should return error Internal if data NewStruct conversion of data failed", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			receiverService: mockedReceiverService,
			logger:          log.NewNoop(),
		}
		data := make(map[string]interface{})
		data["channels"] = string([]byte{0xff})
		payload := &receiver.Receiver{
			Name:           "foo",
			Type:           "bar",
			Labels:         labels,
			Configurations: configurations,
			Data:           data,
		}

		mockedReceiverService.EXPECT().GetReceiver(receiverId).
			Return(payload, nil)
		res, err := dummyGRPCServer.GetReceiver(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.Equal(t, strings.Replace(err.Error(), "\u00a0", " ", -1),
			"rpc error: code = Internal desc = proto: invalid UTF-8 in string: \"\\xff\"")
	})
}

func TestGRPCServer_UpdateReceiver(t *testing.T) {
	configurations := make(map[string]interface{})
	configurations["client_id"] = "foo"
	configurations["client_secret"] = "bar"
	configurations["auth_code"] = "foo"

	labels := make(map[string]string)
	labels["foo"] = "bar"

	configurationsData, _ := structpb.NewStruct(configurations)
	dummyReq := &sirenv1beta1.UpdateReceiverRequest{
		Id:             uint64(22),
		Name:           "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurationsData,
	}
	payload := &receiver.Receiver{
		ID:             uint64(22),
		Name:           "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurations,
	}

	t.Run("should update receiver object", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			receiverService: mockedReceiverService,
			logger:          log.NewNoop(),
		}
		mockedReceiverService.EXPECT().UpdateReceiver(payload).
			Return(nil).Once()

		res, err := dummyGRPCServer.UpdateReceiver(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, payload.ID, res.GetId())
	})

	t.Run("should return error Invalid Argument if update receiver return invalid error", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			receiverService: mockedReceiverService,
			logger:          log.NewNoop(),
		}

		mockedReceiverService.EXPECT().UpdateReceiver(payload).
			Return(fmt.Errorf("some invalid error: %w", receiver.ErrInvalid)).Once()

		res, err := dummyGRPCServer.UpdateReceiver(context.Background(), dummyReq)
		assert.EqualError(t, err,
			"rpc error: code = InvalidArgument desc = some invalid error: bad_request")
		assert.Nil(t, res)
	})

	t.Run("should return error Internal if updating receiver failed", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			receiverService: mockedReceiverService,
			logger:          log.NewNoop(),
		}
		mockedReceiverService.EXPECT().UpdateReceiver(payload).
			Return(errors.New("random error"))

		res, err := dummyGRPCServer.UpdateReceiver(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})
}

func TestGRPCServer_DeleteReceiver(t *testing.T) {
	providerId := uint64(10)
	dummyReq := &sirenv1beta1.DeleteReceiverRequest{
		Id: uint64(10),
	}

	t.Run("should delete receiver object", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			receiverService: mockedReceiverService,
			logger:          log.NewNoop(),
		}
		mockedReceiverService.EXPECT().DeleteReceiver(providerId).
			Return(nil).Once()

		res, err := dummyGRPCServer.DeleteReceiver(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, "", res.String())
	})

	t.Run("should return error Internal if deleting receiver failed", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			receiverService: mockedReceiverService,
			logger:          log.NewNoop(),
		}
		mockedReceiverService.EXPECT().DeleteReceiver(providerId).
			Return(errors.New("random error")).Once()

		res, err := dummyGRPCServer.DeleteReceiver(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})
}

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

		dummyGRPCServer := GRPCServer{
			receiverService: mockedReceiverService,
			logger:          log.NewNoop(),
		}

		mockedReceiverService.EXPECT().NotifyReceiver(
			uint64(1),
			receiver.NotificationMessage{
				"receiver_name": dummyReq.GetPayload().Fields["receiver_name"].GetStringValue(),
				"receiver_type": dummyReq.GetPayload().Fields["receiver_type"].GetStringValue(),
				"message":       dummyReq.GetPayload().Fields["message"].GetStringValue(),
				"blocks":        dummyReq.GetPayload().Fields["blocks"].GetListValue().AsSlice(),
			},
		).Return(fmt.Errorf("some invalid argument error: %w", receiver.ErrInvalid))
		_, err := dummyGRPCServer.NotifyReceiver(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = some invalid argument error: bad_request")
		mockedReceiverService.AssertExpectations(t)
	})

	t.Run("should return internal error if notify receiver return some error", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}

		dummyGRPCServer := GRPCServer{
			receiverService: mockedReceiverService,
			logger:          log.NewNoop(),
		}

		mockedReceiverService.EXPECT().NotifyReceiver(
			uint64(1),
			receiver.NotificationMessage{
				"receiver_name": dummyReq.GetPayload().Fields["receiver_name"].GetStringValue(),
				"receiver_type": dummyReq.GetPayload().Fields["receiver_type"].GetStringValue(),
				"message":       dummyReq.GetPayload().Fields["message"].GetStringValue(),
				"blocks":        dummyReq.GetPayload().Fields["blocks"].GetListValue().AsSlice(),
			},
		).Return(errors.New("some error"))
		_, err := dummyGRPCServer.NotifyReceiver(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some error")
		mockedReceiverService.AssertExpectations(t)
	})

	t.Run("should return OK response if notify receiver succeed", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}

		dummyGRPCServer := GRPCServer{
			receiverService: mockedReceiverService,
			logger:          log.NewNoop(),
		}

		mockedReceiverService.EXPECT().NotifyReceiver(
			uint64(1),
			receiver.NotificationMessage{
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
