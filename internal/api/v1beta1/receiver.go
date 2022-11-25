package v1beta1

import (
	"context"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/pkg/secret"
	sirenv1beta1 "github.com/odpf/siren/proto/odpf/siren/v1beta1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GRPCServer) ListReceivers(ctx context.Context, _ *sirenv1beta1.ListReceiversRequest) (*sirenv1beta1.ListReceiversResponse, error) {
	receivers, err := s.receiverService.List(ctx, receiver.Filter{})
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	items := []*sirenv1beta1.Receiver{}
	for _, rcv := range receivers {
		configurations, err := structpb.NewStruct(sanitizeConfigMap(rcv.Configurations))
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
	if !receiver.IsTypeSupported(req.GetType()) {
		return nil, s.generateRPCErr(errors.ErrInvalid.WithMsgf("unsupported type %s", req.GetType()))
	}

	rcv := &receiver.Receiver{
		Name:           req.GetName(),
		Type:           req.GetType(),
		Labels:         req.GetLabels(),
		Configurations: req.GetConfigurations().AsMap(),
	}

	err := s.receiverService.Create(ctx, rcv)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.CreateReceiverResponse{
		Id: rcv.ID,
	}, nil
}

func (s *GRPCServer) GetReceiver(ctx context.Context, req *sirenv1beta1.GetReceiverRequest) (*sirenv1beta1.GetReceiverResponse, error) {
	rcv, err := s.receiverService.Get(ctx, req.GetId(), receiver.GetWithData(true))
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	data, err := structpb.NewStruct(rcv.Data)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	configuration, err := structpb.NewStruct(sanitizeConfigMap(rcv.Configurations))
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
	rcv := &receiver.Receiver{
		ID:             req.GetId(),
		Name:           req.GetName(),
		Labels:         req.GetLabels(),
		Configurations: req.GetConfigurations().AsMap(),
	}

	err := s.receiverService.Update(ctx, rcv)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.UpdateReceiverResponse{
		Id: rcv.ID,
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
	payloadMap := req.GetPayload().AsMap()

	n := notification.Notification{}
	if err := mapstructure.Decode(payloadMap, &n); err != nil {
		err = errors.ErrInvalid.WithMsgf("failed to parse payload to notification: %s", err.Error())
		return nil, s.generateRPCErr(err)
	}

	if err := s.notificationService.DispatchToReceiver(ctx, n, req.GetId()); err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.NotifyReceiverResponse{}, nil
}

// sanitizeConfigMap does all sanitization to present receiver configurations to the user
func sanitizeConfigMap(receiverConfigMap map[string]interface{}) map[string]interface{} {
	var newConfigMap = make(map[string]interface{})
	for k, v := range receiverConfigMap {
		// sanitize maskable string. convert `secret.MaskableString` to string to be compatible with structpb
		if reflect.TypeOf(v) == reflect.TypeOf(secret.MaskableString("")) {
			newConfigMap[k] = fmt.Sprintf("%v", v)
		} else {
			newConfigMap[k] = v
		}
	}
	return newConfigMap
}
