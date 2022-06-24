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
	List(context.Context) ([]receiver.Receiver, error)
	Create(context.Context, *receiver.Receiver) (uint64, error)
	Get(context.Context, uint64) (*receiver.Receiver, error)
	Update(context.Context, *receiver.Receiver) (uint64, error)
	Delete(context.Context, uint64) error
	Notify(context.Context, uint64, receiver.NotificationMessage) error
}

func (s *GRPCServer) ListReceivers(ctx context.Context, _ *sirenv1beta1.ListReceiversRequest) (*sirenv1beta1.ListReceiversResponse, error) {
	receivers, err := s.receiverService.List(ctx)
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

func (s *GRPCServer) CreateReceiver(ctx context.Context, req *sirenv1beta1.CreateReceiverRequest) (*sirenv1beta1.CreateReceiverResponse, error) {
	configurations := req.GetConfigurations().AsMap()

	rcv := &receiver.Receiver{
		Name:           req.GetName(),
		Type:           req.GetType(),
		Labels:         req.GetLabels(),
		Configurations: configurations,
	}

	id, err := s.receiverService.Create(ctx, rcv)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.CreateReceiverResponse{
		Id: id,
	}, nil
}

func (s *GRPCServer) GetReceiver(ctx context.Context, req *sirenv1beta1.GetReceiverRequest) (*sirenv1beta1.GetReceiverResponse, error) {
	rcv, err := s.receiverService.Get(ctx, req.GetId())
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

func (s *GRPCServer) UpdateReceiver(ctx context.Context, req *sirenv1beta1.UpdateReceiverRequest) (*sirenv1beta1.UpdateReceiverResponse, error) {
	configurations := req.GetConfigurations().AsMap()

	rcv := &receiver.Receiver{
		ID:             req.GetId(),
		Name:           req.GetName(),
		Type:           req.GetType(),
		Labels:         req.GetLabels(),
		Configurations: configurations,
	}
	id, err := s.receiverService.Update(ctx, rcv)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.UpdateReceiverResponse{
		Id: id,
	}, nil
}

func (s *GRPCServer) DeleteReceiver(ctx context.Context, req *sirenv1beta1.DeleteReceiverRequest) (*sirenv1beta1.DeleteReceiverResponse, error) {
	err := s.receiverService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.DeleteReceiverResponse{}, nil
}

func (s *GRPCServer) NotifyReceiver(ctx context.Context, req *sirenv1beta1.NotifyReceiverRequest) (*sirenv1beta1.NotifyReceiverResponse, error) {
	if err := s.receiverService.Notify(ctx, req.GetId(), req.GetPayload().AsMap()); err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.NotifyReceiverResponse{}, nil
}
