package v1

import (
	"context"
	"encoding/json"
	sirenv1 "github.com/odpf/siren/api/proto/odpf/siren/v1"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/helper"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	Slack     string = "slack"
	Pagerduty string = "pagerduty"
	Http      string = "http"
)

func (s *GRPCServer) ListReceivers(_ context.Context, _ *emptypb.Empty) (*sirenv1.ListReceiversResponse, error) {
	receivers, err := s.container.ReceiverService.ListReceivers()
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	res := &sirenv1.ListReceiversResponse{
		Receivers: make([]*sirenv1.Receiver, 0),
	}
	for _, receiver := range receivers {
		configurations, err := structpb.NewStruct(receiver.Configurations)
		if err != nil {
			return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
		}

		item := &sirenv1.Receiver{
			Id:             receiver.Id,
			Name:           receiver.Name,
			Type:           receiver.Type,
			Configurations: configurations,
			Labels:         receiver.Labels,
			CreatedAt:      timestamppb.New(receiver.CreatedAt),
			UpdatedAt:      timestamppb.New(receiver.UpdatedAt),
		}
		res.Receivers = append(res.Receivers, item)
	}
	return res, nil
}

func (s *GRPCServer) CreateReceiver(_ context.Context, req *sirenv1.CreateReceiverRequest) (*sirenv1.Receiver, error) {
	configurations := req.GetConfigurations().AsMap()

	switch receiverType := req.GetType(); receiverType {
	case Slack:
		err := validateSlackConfigurations(configurations)
		if err != nil {
			return nil, err
		}
	case Pagerduty:
		err := validatePagerdutyConfigurations(configurations)
		if err != nil {
			return nil, err
		}
	case Http:
		err := validateHttpConfigurations(configurations)
		if err != nil {
			return nil, err
		}
	default:
		return nil, status.Errorf(codes.InvalidArgument, "receiver not supported")
	}

	receiver, err := s.container.ReceiverService.CreateReceiver(&domain.Receiver{
		Name:           req.GetName(),
		Type:           req.GetType(),
		Labels:         req.GetLabels(),
		Configurations: configurations,
	})
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	c, err := structpb.NewStruct(receiver.Configurations)
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	return &sirenv1.Receiver{
		Id:             receiver.Id,
		Name:           receiver.Name,
		Type:           receiver.Type,
		Labels:         receiver.Labels,
		Configurations: c,
		CreatedAt:      timestamppb.New(receiver.CreatedAt),
		UpdatedAt:      timestamppb.New(receiver.UpdatedAt),
	}, nil
}

func (s *GRPCServer) GetReceiver(_ context.Context, req *sirenv1.GetReceiverRequest) (*sirenv1.Receiver, error) {
	receiver, err := s.container.ReceiverService.GetReceiver(req.GetId())
	if receiver == nil {
		return nil, status.Errorf(codes.NotFound, "receiver not found")
	}
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	data, err := structpb.NewStruct(receiver.Data)
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	configuration, err := structpb.NewStruct(receiver.Configurations)
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	return &sirenv1.Receiver{
		Id:             receiver.Id,
		Name:           receiver.Name,
		Type:           receiver.Type,
		Labels:         receiver.Labels,
		Configurations: configuration,
		Data:           data,
		CreatedAt:      timestamppb.New(receiver.CreatedAt),
		UpdatedAt:      timestamppb.New(receiver.UpdatedAt),
	}, nil
}

func (s *GRPCServer) UpdateReceiver(_ context.Context, req *sirenv1.UpdateReceiverRequest) (*sirenv1.Receiver, error) {
	configurations := req.GetConfigurations().AsMap()

	switch receiverType := req.GetType(); receiverType {
	case Slack:
		err := validateSlackConfigurations(configurations)
		if err != nil {
			return nil, err
		}
	case Pagerduty:
		err := validatePagerdutyConfigurations(configurations)
		if err != nil {
			return nil, err
		}
	case Http:
		err := validateHttpConfigurations(configurations)
		if err != nil {
			return nil, err
		}
	default:
		return nil, status.Errorf(codes.InvalidArgument, "receiver not supported")
	}

	receiver, err := s.container.ReceiverService.UpdateReceiver(&domain.Receiver{
		Id:             req.GetId(),
		Name:           req.GetName(),
		Type:           req.GetType(),
		Labels:         req.GetLabels(),
		Configurations: configurations,
	})
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	configuration, err := structpb.NewStruct(receiver.Configurations)
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	return &sirenv1.Receiver{
		Id:             receiver.Id,
		Name:           receiver.Name,
		Type:           receiver.Type,
		Labels:         receiver.Labels,
		Configurations: configuration,
		CreatedAt:      timestamppb.New(receiver.CreatedAt),
		UpdatedAt:      timestamppb.New(receiver.UpdatedAt),
	}, nil
}

func (s *GRPCServer) DeleteReceiver(_ context.Context, req *sirenv1.DeleteReceiverRequest) (*emptypb.Empty, error) {
	err := s.container.ReceiverService.DeleteReceiver(uint64(req.GetId()))
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	return &emptypb.Empty{}, nil
}

func (s *GRPCServer) SendReceiverNotification(_ context.Context, req *sirenv1.SendReceiverNotificationRequest) (*sirenv1.SendReceiverNotificationResponse, error) {
	var res *sirenv1.SendReceiverNotificationResponse
	receiver, err := s.container.ReceiverService.GetReceiver(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	switch receiver.Type {
	case Slack:
		slackPayload := req.GetSlack()

		b, err := json.Marshal(slackPayload.GetBlocks())
		if err != nil {
			s.logger.Error("handler", zap.Error(err))
			return nil, status.Errorf(codes.InvalidArgument, "Invalid block")
		}

		blocks := slack.Blocks{}
		err = json.Unmarshal(b, &blocks)
		if err != nil {
			s.logger.Error("handler", zap.Error(err))
			return nil, status.Errorf(codes.InvalidArgument, "unable to parse block")
		}

		payload := &domain.SlackMessage{
			ReceiverName: slackPayload.GetReceiverName(),
			ReceiverType: slackPayload.GetReceiverType(),
			Token:        receiver.Configurations["token"].(string),
			Message:      slackPayload.GetMessage(),
			Blocks:       blocks,
		}
		result, err := s.container.NotifierServices.Slack.Notify(payload)
		if err != nil {
			return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
		}
		res = &sirenv1.SendReceiverNotificationResponse{
			Ok: result.OK,
		}
	default:
		return nil, status.Errorf(codes.NotFound, "Send notification not registered for this receiver")
	}
	return res, nil
}

func validateSlackConfigurations(configurations map[string]interface{}) error {
	_, err := helper.GetMapString(configurations, "configurations", "client_id")
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	_, err = helper.GetMapString(configurations, "configurations", "client_secret")
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	_, err = helper.GetMapString(configurations, "configurations", "auth_code")
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}
	return nil
}

func validatePagerdutyConfigurations(configurations map[string]interface{}) error {
	_, err := helper.GetMapString(configurations, "configurations", "service_key")
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}
	return nil
}

func validateHttpConfigurations(configurations map[string]interface{}) error {
	_, err := helper.GetMapString(configurations, "configurations", "url")
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}
	return nil
}
