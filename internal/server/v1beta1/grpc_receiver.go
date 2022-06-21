package v1beta1

import (
	"context"
	"errors"

	"github.com/odpf/siren/core/receiver"
	sirenv1beta1 "github.com/odpf/siren/internal/server/proto/odpf/siren/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate mockery --name=ReceiverService -r --case underscore --with-expecter --structname ReceiverService --filename receiver_service.go --output=./mocks
type ReceiverService interface {
	ListReceivers() ([]*receiver.Receiver, error)
	CreateReceiver(*receiver.Receiver) error
	GetReceiver(uint64) (*receiver.Receiver, error)
	UpdateReceiver(*receiver.Receiver) error
	DeleteReceiver(uint64) error
	NotifyReceiver(id uint64, payloadMessage receiver.NotificationMessage) error
}

func (s *GRPCServer) ListReceivers(_ context.Context, _ *emptypb.Empty) (*sirenv1beta1.ListReceiversResponse, error) {
	receivers, err := s.receiverService.ListReceivers()
	if err != nil {
		return nil, gRPCLogError(s.logger, codes.Internal, err)
	}

	res := &sirenv1beta1.ListReceiversResponse{
		Data: make([]*sirenv1beta1.Receiver, 0),
	}
	for _, rcv := range receivers {
		configurations, err := structpb.NewStruct(rcv.Configurations)
		if err != nil {
			return nil, gRPCLogError(s.logger, codes.Internal, err)
		}

		item := &sirenv1beta1.Receiver{
			Id:             rcv.ID,
			Name:           rcv.Name,
			Type:           rcv.Type,
			Configurations: configurations,
			Labels:         rcv.Labels,
			CreatedAt:      timestamppb.New(rcv.CreatedAt),
			UpdatedAt:      timestamppb.New(rcv.UpdatedAt),
		}
		res.Data = append(res.Data, item)
	}
	return res, nil
}

func (s *GRPCServer) CreateReceiver(_ context.Context, req *sirenv1beta1.CreateReceiverRequest) (*sirenv1beta1.CreateReceiverResponse, error) {
	configurations := req.GetConfigurations().AsMap()

	rcv := &receiver.Receiver{
		Name:           req.GetName(),
		Type:           req.GetType(),
		Labels:         req.GetLabels(),
		Configurations: configurations,
	}

	if err := s.receiverService.CreateReceiver(rcv); err != nil {
		if errors.Is(err, receiver.ErrInvalid) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, gRPCLogError(s.logger, codes.Internal, err)
	}

	return &sirenv1beta1.CreateReceiverResponse{
		Id: rcv.ID,
	}, nil
}

func (s *GRPCServer) GetReceiver(_ context.Context, req *sirenv1beta1.GetReceiverRequest) (*sirenv1beta1.GetReceiverResponse, error) {
	rcv, err := s.receiverService.GetReceiver(req.GetId())
	if rcv == nil {
		return nil, status.Errorf(codes.NotFound, "receiver not found")
	}
	if err != nil {
		return nil, gRPCLogError(s.logger, codes.Internal, err)
	}

	data, err := structpb.NewStruct(rcv.Data)
	if err != nil {
		return nil, gRPCLogError(s.logger, codes.Internal, err)
	}

	configuration, err := structpb.NewStruct(rcv.Configurations)
	if err != nil {
		return nil, gRPCLogError(s.logger, codes.Internal, err)
	}

	return &sirenv1beta1.GetReceiverResponse{
		Data: &sirenv1beta1.Receiver{
			Id:             rcv.ID,
			Name:           rcv.Name,
			Type:           rcv.Type,
			Labels:         rcv.Labels,
			Configurations: configuration,
			Data:           data,
			CreatedAt:      timestamppb.New(rcv.CreatedAt),
			UpdatedAt:      timestamppb.New(rcv.UpdatedAt),
		},
	}, nil
}

func (s *GRPCServer) UpdateReceiver(_ context.Context, req *sirenv1beta1.UpdateReceiverRequest) (*sirenv1beta1.UpdateReceiverResponse, error) {
	configurations := req.GetConfigurations().AsMap()

	rcv := &receiver.Receiver{
		ID:             req.GetId(),
		Name:           req.GetName(),
		Type:           req.GetType(),
		Labels:         req.GetLabels(),
		Configurations: configurations,
	}
	if err := s.receiverService.UpdateReceiver(rcv); err != nil {
		if errors.Is(err, receiver.ErrInvalid) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, gRPCLogError(s.logger, codes.Internal, err)
	}

	return &sirenv1beta1.UpdateReceiverResponse{
		Id: rcv.ID,
	}, nil
}

func (s *GRPCServer) DeleteReceiver(_ context.Context, req *sirenv1beta1.DeleteReceiverRequest) (*sirenv1beta1.DeleteReceiverResponse, error) {
	err := s.receiverService.DeleteReceiver(uint64(req.GetId()))
	if err != nil {
		return nil, gRPCLogError(s.logger, codes.Internal, err)
	}

	return &sirenv1beta1.DeleteReceiverResponse{}, nil
}

func (s *GRPCServer) NotifyReceiver(_ context.Context, req *sirenv1beta1.NotifyReceiverRequest) (*sirenv1beta1.NotifyReceiverResponse, error) {
	if err := s.receiverService.NotifyReceiver(req.GetId(), req.GetPayload().AsMap()); err != nil {
		if errors.Is(err, receiver.ErrInvalid) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, gRPCLogError(s.logger, codes.Internal, err)
	}

	return &sirenv1beta1.NotifyReceiverResponse{}, nil
}
