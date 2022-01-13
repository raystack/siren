package v1

import (
	"context"
	sirenv1beta1 "github.com/odpf/siren/api/proto/odpf/siren/v1beta1"
	"github.com/odpf/siren/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

func (s *GRPCServer) ListSubscriptions(_ context.Context, _ *emptypb.Empty) (*sirenv1beta1.ListSubscriptionsResponse, error) {
	subscriptions, err := s.container.SubscriptionService.ListSubscriptions()
	if err != nil {
		s.log.Error("failed to list subscriptions", "error", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res := &sirenv1beta1.ListSubscriptionsResponse{
		Subscriptions: make([]*sirenv1beta1.Subscription, 0),
	}
	for _, subscription := range subscriptions {
		item := &sirenv1beta1.Subscription{
			Id:        subscription.Id,
			Urn:       subscription.Urn,
			Namespace: subscription.Namespace,
			Match:     subscription.Match,
			Receivers: getReceiverMetadataListFromDomainObject(subscription.Receivers),
			CreatedAt: timestamppb.New(subscription.CreatedAt),
			UpdatedAt: timestamppb.New(subscription.UpdatedAt),
		}
		res.Subscriptions = append(res.Subscriptions, item)
	}
	return res, nil
}

func (s *GRPCServer) CreateSubscription(_ context.Context, req *sirenv1beta1.CreateSubscriptionRequest) (*sirenv1beta1.Subscription, error) {
	subscription, err := s.container.SubscriptionService.CreateSubscription(&domain.Subscription{
		Namespace: req.GetNamespace(),
		Urn:       req.GetUrn(),
		Receivers: getReceiverMetadataListInDomainObject(req.GetReceivers()),
		Match:     req.GetMatch(),
	})
	if err != nil {
		s.log.Error("failed to create subscription", "error", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	receivers := make([]*sirenv1beta1.ReceiverMetadata, 0)
	for _, receiverMetadataItem := range subscription.Receivers {
		item := getReceiverMetadataFromDomainObject(&receiverMetadataItem)
		receivers = append(receivers, &item)
	}
	return &sirenv1beta1.Subscription{
		Id:        subscription.Id,
		Urn:       subscription.Urn,
		Namespace: subscription.Namespace,
		Match:     subscription.Match,
		Receivers: receivers,
		CreatedAt: timestamppb.New(subscription.CreatedAt),
		UpdatedAt: timestamppb.New(subscription.UpdatedAt),
	}, nil
}

func (s *GRPCServer) GetSubscription(_ context.Context, req *sirenv1beta1.GetSubscriptionRequest) (*sirenv1beta1.Subscription, error) {
	subscription, err := s.container.SubscriptionService.GetSubscription(req.GetId())
	if err != nil {
		s.log.Error("failed to fetch subscription", "error", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if subscription == nil {
		return nil, status.Errorf(codes.NotFound, "subscription not found")
	}

	receivers := make([]*sirenv1beta1.ReceiverMetadata, 0)
	for _, receiverMetadataItem := range subscription.Receivers {
		item := getReceiverMetadataFromDomainObject(&receiverMetadataItem)
		receivers = append(receivers, &item)
	}

	return &sirenv1beta1.Subscription{
		Id:        subscription.Id,
		Urn:       subscription.Urn,
		Namespace: subscription.Namespace,
		Match:     subscription.Match,
		Receivers: receivers,
		CreatedAt: timestamppb.New(subscription.CreatedAt),
		UpdatedAt: timestamppb.New(subscription.UpdatedAt),
	}, nil
}

func (s *GRPCServer) UpdateSubscription(_ context.Context, req *sirenv1beta1.UpdateSubscriptionRequest) (*sirenv1beta1.Subscription, error) {
	subscription, err := s.container.SubscriptionService.UpdateSubscription(&domain.Subscription{
		Id:        req.GetId(),
		Namespace: req.GetNamespace(),
		Urn:       req.GetUrn(),
		Receivers: getReceiverMetadataListInDomainObject(req.GetReceivers()),
		Match:     req.GetMatch(),
	})
	if err != nil {
		if strings.Contains(err.Error(), `violates unique constraint "urn_provider_id_unique"`) {
			return nil, status.Errorf(codes.InvalidArgument, "urn and provider pair already exist")
		}
		s.log.Error("failed to update subscription", "error", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	receivers := make([]*sirenv1beta1.ReceiverMetadata, 0)
	for _, receiverMetadataItem := range subscription.Receivers {
		item := getReceiverMetadataFromDomainObject(&receiverMetadataItem)
		receivers = append(receivers, &item)
	}

	return &sirenv1beta1.Subscription{
		Id:        subscription.Id,
		Urn:       subscription.Urn,
		Namespace: subscription.Namespace,
		Match:     subscription.Match,
		Receivers: receivers,
		CreatedAt: timestamppb.New(subscription.CreatedAt),
		UpdatedAt: timestamppb.New(subscription.UpdatedAt),
	}, nil
}

func (s *GRPCServer) DeleteSubscription(_ context.Context, req *sirenv1beta1.DeleteSubscriptionRequest) (*emptypb.Empty, error) {
	err := s.container.SubscriptionService.DeleteSubscription(req.GetId())
	if err != nil {
		s.log.Error("failed to delete subscription", "error", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func getReceiverMetadataFromDomainObject(item *domain.ReceiverMetadata) sirenv1beta1.ReceiverMetadata {
	return sirenv1beta1.ReceiverMetadata{
		Id:            item.Id,
		Configuration: item.Configuration,
	}
}

func getReceiverMetadataInDomainObject(item *sirenv1beta1.ReceiverMetadata) domain.ReceiverMetadata {
	return domain.ReceiverMetadata{
		Id:            item.Id,
		Configuration: item.Configuration,
	}
}

func getReceiverMetadataListInDomainObject(domainReceivers []*sirenv1beta1.ReceiverMetadata) []domain.ReceiverMetadata {
	receivers := make([]domain.ReceiverMetadata, 0)
	for _, receiverMetadataItem := range domainReceivers {
		receivers = append(receivers, getReceiverMetadataInDomainObject(receiverMetadataItem))
	}
	return receivers
}

func getReceiverMetadataListFromDomainObject(domainReceivers []domain.ReceiverMetadata) []*sirenv1beta1.ReceiverMetadata {
	receivers := make([]*sirenv1beta1.ReceiverMetadata, 0)
	for _, receiverMetadataItem := range domainReceivers {
		item := getReceiverMetadataFromDomainObject(&receiverMetadataItem)
		receivers = append(receivers, &item)
	}
	return receivers
}
