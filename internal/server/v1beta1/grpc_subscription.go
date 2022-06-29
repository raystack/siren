package v1beta1

import (
	"context"

	"github.com/odpf/siren/core/subscription"
	sirenv1beta1 "github.com/odpf/siren/internal/server/proto/odpf/siren/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate mockery --name=SubscriptionService -r --case underscore --with-expecter --structname SubscriptionService --filename subscription_service.go --output=./mocks
type SubscriptionService interface {
	ListSubscriptions(context.Context) ([]*subscription.Subscription, error)
	CreateSubscription(context.Context, *subscription.Subscription) error
	GetSubscription(context.Context, uint64) (*subscription.Subscription, error)
	UpdateSubscription(context.Context, *subscription.Subscription) error
	DeleteSubscription(context.Context, uint64) error
}

func (s *GRPCServer) ListSubscriptions(ctx context.Context, _ *sirenv1beta1.ListSubscriptionsRequest) (*sirenv1beta1.ListSubscriptionsResponse, error) {
	subscriptions, err := s.subscriptionService.ListSubscriptions(ctx)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	items := []*sirenv1beta1.Subscription{}

	for _, sub := range subscriptions {
		item := &sirenv1beta1.Subscription{
			Id:        sub.ID,
			Urn:       sub.URN,
			Namespace: sub.Namespace,
			Match:     sub.Match,
			Receivers: getReceiverMetadataListFromDomainObject(sub.Receivers),
			CreatedAt: timestamppb.New(sub.CreatedAt),
			UpdatedAt: timestamppb.New(sub.UpdatedAt),
		}
		items = append(items, item)
	}
	return &sirenv1beta1.ListSubscriptionsResponse{
		Subscriptions: items,
	}, nil
}

func (s *GRPCServer) CreateSubscription(ctx context.Context, req *sirenv1beta1.CreateSubscriptionRequest) (*sirenv1beta1.CreateSubscriptionResponse, error) {
	sub := &subscription.Subscription{
		Namespace: req.GetNamespace(),
		URN:       req.GetUrn(),
		Receivers: getReceiverMetadataListInDomainObject(req.GetReceivers()),
		Match:     req.GetMatch(),
	}
	if err := s.subscriptionService.CreateSubscription(ctx, sub); err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.CreateSubscriptionResponse{
		Id: sub.ID,
	}, nil
}

func (s *GRPCServer) GetSubscription(ctx context.Context, req *sirenv1beta1.GetSubscriptionRequest) (*sirenv1beta1.GetSubscriptionResponse, error) {
	sub, err := s.subscriptionService.GetSubscription(ctx, req.GetId())
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	receivers := make([]*sirenv1beta1.ReceiverMetadata, 0)
	for _, receiverMetadataItem := range sub.Receivers {
		item := getReceiverMetadataFromDomainObject(&receiverMetadataItem)
		receivers = append(receivers, &item)
	}

	return &sirenv1beta1.GetSubscriptionResponse{
		Subscription: &sirenv1beta1.Subscription{
			Id:        sub.ID,
			Urn:       sub.URN,
			Namespace: sub.Namespace,
			Match:     sub.Match,
			Receivers: receivers,
			CreatedAt: timestamppb.New(sub.CreatedAt),
			UpdatedAt: timestamppb.New(sub.UpdatedAt),
		},
	}, nil
}

func (s *GRPCServer) UpdateSubscription(ctx context.Context, req *sirenv1beta1.UpdateSubscriptionRequest) (*sirenv1beta1.UpdateSubscriptionResponse, error) {
	sub := &subscription.Subscription{
		ID:        req.GetId(),
		Namespace: req.GetNamespace(),
		URN:       req.GetUrn(),
		Receivers: getReceiverMetadataListInDomainObject(req.GetReceivers()),
		Match:     req.GetMatch(),
	}
	if err := s.subscriptionService.UpdateSubscription(ctx, sub); err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.UpdateSubscriptionResponse{
		Id: sub.ID,
	}, nil
}

func (s *GRPCServer) DeleteSubscription(ctx context.Context, req *sirenv1beta1.DeleteSubscriptionRequest) (*sirenv1beta1.DeleteSubscriptionResponse, error) {
	err := s.subscriptionService.DeleteSubscription(ctx, req.GetId())
	if err != nil {
		return nil, s.generateRPCErr(err)
	}
	return &sirenv1beta1.DeleteSubscriptionResponse{}, nil
}

func getReceiverMetadataFromDomainObject(item *subscription.ReceiverMetadata) sirenv1beta1.ReceiverMetadata {
	return sirenv1beta1.ReceiverMetadata{
		Id:            item.ID,
		Configuration: item.Configuration,
	}
}

func getReceiverMetadataInDomainObject(item *sirenv1beta1.ReceiverMetadata) subscription.ReceiverMetadata {
	return subscription.ReceiverMetadata{
		ID:            item.Id,
		Configuration: item.Configuration,
	}
}

func getReceiverMetadataListInDomainObject(domainReceivers []*sirenv1beta1.ReceiverMetadata) []subscription.ReceiverMetadata {
	receivers := make([]subscription.ReceiverMetadata, 0)
	for _, receiverMetadataItem := range domainReceivers {
		receivers = append(receivers, getReceiverMetadataInDomainObject(receiverMetadataItem))
	}
	return receivers
}

func getReceiverMetadataListFromDomainObject(domainReceivers []subscription.ReceiverMetadata) []*sirenv1beta1.ReceiverMetadata {
	receivers := make([]*sirenv1beta1.ReceiverMetadata, 0)
	for _, receiverMetadataItem := range domainReceivers {
		item := getReceiverMetadataFromDomainObject(&receiverMetadataItem)
		receivers = append(receivers, &item)
	}
	return receivers
}
