package v1beta1

import (
	"context"

	"github.com/odpf/siren/core/receiver"
	sirenv1beta1 "github.com/odpf/siren/internal/server/proto/odpf/siren/v1beta1"
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

func (s *GRPCServer) ListReceivers(_ context.Context, _ *sirenv1beta1.ListReceiversRequest) (*sirenv1beta1.ListReceiversResponse, error) {
	receivers, err := s.receiverService.ListReceivers()
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	items := []*sirenv1beta1.Receiver{}
	for _, rcv := range receivers {
		configurations, err := structpb.NewStruct(rcv.Configurations)
		if err != nil {
			return nil, s.generateRPCErr(err)
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
		items = append(items, item)
	}
	return &sirenv1beta1.ListReceiversResponse{
		Receivers: items,
	}, nil
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
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.CreateReceiverResponse{
		Id: rcv.ID,
	}, nil
}

func (s *GRPCServer) GetReceiver(_ context.Context, req *sirenv1beta1.GetReceiverRequest) (*sirenv1beta1.GetReceiverResponse, error) {
	rcv, err := s.receiverService.GetReceiver(req.GetId())
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	data, err := structpb.NewStruct(rcv.Data)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	configuration, err := structpb.NewStruct(rcv.Configurations)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.GetReceiverResponse{
		Receiver: &sirenv1beta1.Receiver{
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
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.UpdateReceiverResponse{
		Id: rcv.ID,
	}, nil
}

func (s *GRPCServer) DeleteReceiver(_ context.Context, req *sirenv1beta1.DeleteReceiverRequest) (*sirenv1beta1.DeleteReceiverResponse, error) {
	err := s.receiverService.DeleteReceiver(uint64(req.GetId()))
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.DeleteReceiverResponse{}, nil
}

func (s *GRPCServer) NotifyReceiver(_ context.Context, req *sirenv1beta1.NotifyReceiverRequest) (*sirenv1beta1.NotifyReceiverResponse, error) {
	if err := s.receiverService.NotifyReceiver(req.GetId(), req.GetPayload().AsMap()); err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.NotifyReceiverResponse{}, nil
}
