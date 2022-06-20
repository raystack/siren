package v1beta1

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/pkg/slack"
	"github.com/odpf/siren/utils"
	sirenv1beta1 "go.buf.build/odpf/gw/odpf/proton/odpf/siren/v1beta1"
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

//go:generate mockery --name=ReceiverService -r --case underscore --with-expecter --structname ReceiverService --filename receiver_service.go --output=./mocks
type ReceiverService interface {
	ListReceivers() ([]*receiver.Receiver, error)
	CreateReceiver(*receiver.Receiver) error
	GetReceiver(uint64) (*receiver.Receiver, error)
	UpdateReceiver(*receiver.Receiver) error
	DeleteReceiver(uint64) error
	NotifyReceiver(rcv *receiver.Receiver, payloadMessage string, payloadReceiverName string, payloadReceiverType string, payloadBlock []byte) error
	Migrate() error
}

type NotifierServices struct { //TODO to be refactored, temporary only
	Slack SlackNotifierService
}

//go:generate mockery --name=SlackNotifierService -r --case underscore --with-expecter --structname SlackNotifierService --filename slack_notifier_service.go --output=./mocks
type SlackNotifierService interface { //TODO to be refactored, temporary only
	Notify(*slack.Message, ...slack.ClientCallOption) error
}

func (s *GRPCServer) ListReceivers(_ context.Context, _ *emptypb.Empty) (*sirenv1beta1.ListReceiversResponse, error) {
	receivers, err := s.receiverService.ListReceivers()
	if err != nil {
		return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
	}

	res := &sirenv1beta1.ListReceiversResponse{
		Receivers: make([]*sirenv1beta1.Receiver, 0),
	}
	for _, receiver := range receivers {
		configurations, err := structpb.NewStruct(receiver.Configurations)
		if err != nil {
			return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
		}

		item := &sirenv1beta1.Receiver{
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

func (s *GRPCServer) CreateReceiver(_ context.Context, req *sirenv1beta1.CreateReceiverRequest) (*sirenv1beta1.Receiver, error) {
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

	receiver := &receiver.Receiver{
		Name:           req.GetName(),
		Type:           req.GetType(),
		Labels:         req.GetLabels(),
		Configurations: configurations,
	}
	if err := s.receiverService.CreateReceiver(receiver); err != nil {
		return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
	}

	c, err := structpb.NewStruct(receiver.Configurations)
	if err != nil {
		return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
	}

	return &sirenv1beta1.Receiver{
		Id:             receiver.Id,
		Name:           receiver.Name,
		Type:           receiver.Type,
		Labels:         receiver.Labels,
		Configurations: c,
		CreatedAt:      timestamppb.New(receiver.CreatedAt),
		UpdatedAt:      timestamppb.New(receiver.UpdatedAt),
	}, nil
}

func (s *GRPCServer) GetReceiver(_ context.Context, req *sirenv1beta1.GetReceiverRequest) (*sirenv1beta1.Receiver, error) {
	receiver, err := s.receiverService.GetReceiver(req.GetId())
	if receiver == nil {
		return nil, status.Errorf(codes.NotFound, "receiver not found")
	}
	if err != nil {
		return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
	}

	data, err := structpb.NewStruct(receiver.Data)
	if err != nil {
		return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
	}

	configuration, err := structpb.NewStruct(receiver.Configurations)
	if err != nil {
		return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
	}

	return &sirenv1beta1.Receiver{
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

func (s *GRPCServer) UpdateReceiver(_ context.Context, req *sirenv1beta1.UpdateReceiverRequest) (*sirenv1beta1.Receiver, error) {
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

	receiver := &receiver.Receiver{
		Id:             req.GetId(),
		Name:           req.GetName(),
		Type:           req.GetType(),
		Labels:         req.GetLabels(),
		Configurations: configurations,
	}
	if err := s.receiverService.UpdateReceiver(receiver); err != nil {
		return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
	}

	configuration, err := structpb.NewStruct(receiver.Configurations)
	if err != nil {
		return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
	}

	return &sirenv1beta1.Receiver{
		Id:             receiver.Id,
		Name:           receiver.Name,
		Type:           receiver.Type,
		Labels:         receiver.Labels,
		Configurations: configuration,
		CreatedAt:      timestamppb.New(receiver.CreatedAt),
		UpdatedAt:      timestamppb.New(receiver.UpdatedAt),
	}, nil
}

func (s *GRPCServer) DeleteReceiver(_ context.Context, req *sirenv1beta1.DeleteReceiverRequest) (*emptypb.Empty, error) {
	err := s.receiverService.DeleteReceiver(uint64(req.GetId()))
	if err != nil {
		return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
	}

	return &emptypb.Empty{}, nil
}

func (s *GRPCServer) SendReceiverNotification(_ context.Context, req *sirenv1beta1.SendReceiverNotificationRequest) (*sirenv1beta1.SendReceiverNotificationResponse, error) {
	var res *sirenv1beta1.SendReceiverNotificationResponse
	rcv, err := s.receiverService.GetReceiver(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	switch rcv.Type {
	case Slack:
		slackPayload := req.GetSlack()

		b, err := json.Marshal(slackPayload.GetBlocks())
		if err != nil {
			s.logger.Error("failed to encode the payload JSON", "error", err)
			return nil, status.Errorf(codes.InvalidArgument, "Invalid block")
		}
		if err := s.receiverService.NotifyReceiver(rcv, slackPayload.GetMessage(), slackPayload.GetReceiverName(), slackPayload.GetReceiverType(), b); err != nil {
			if errors.Is(err, receiver.ErrInvalid) {
				return nil, status.Errorf(codes.InvalidArgument, err.Error())
			}
			return nil, utils.GRPCLogError(s.logger, codes.Internal, err)
		}
		res = &sirenv1beta1.SendReceiverNotificationResponse{
			Ok: true,
		}
	default:
		return nil, status.Errorf(codes.NotFound, "Send notification not registered for this receiver")
	}
	return res, nil
}

func validateSlackConfigurations(configurations map[string]interface{}) error {
	_, err := utils.GetMapString(configurations, "configurations", "client_id")
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	_, err = utils.GetMapString(configurations, "configurations", "client_secret")
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	_, err = utils.GetMapString(configurations, "configurations", "auth_code")
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}
	return nil
}

func validatePagerdutyConfigurations(configurations map[string]interface{}) error {
	_, err := utils.GetMapString(configurations, "configurations", "service_key")
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}
	return nil
}

func validateHttpConfigurations(configurations map[string]interface{}) error {
	_, err := utils.GetMapString(configurations, "configurations", "url")
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}
	return nil
}
