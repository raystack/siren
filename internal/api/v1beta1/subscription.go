package v1beta1

import (
	"context"
	"fmt"

	"github.com/odpf/siren/core/subscription"
	sirenv1beta1 "github.com/odpf/siren/proto/odpf/siren/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GRPCServer) ListSubscriptions(ctx context.Context, _ *sirenv1beta1.ListSubscriptionsRequest) (*sirenv1beta1.ListSubscriptionsResponse, error) {
	subscriptions, err := s.subscriptionService.List(ctx, subscription.Filter{})
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

	err := s.subscriptionService.Create(ctx, sub)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.CreateSubscriptionResponse{
		Id: sub.ID,
	}, nil
}

func (s *GRPCServer) GetSubscription(ctx context.Context, req *sirenv1beta1.GetSubscriptionRequest) (*sirenv1beta1.GetSubscriptionResponse, error) {
	sub, err := s.subscriptionService.Get(ctx, req.GetId())
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

	err := s.subscriptionService.Update(ctx, sub)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.UpdateSubscriptionResponse{
		Id: sub.ID,
	}, nil
}

func (s *GRPCServer) DeleteSubscription(ctx context.Context, req *sirenv1beta1.DeleteSubscriptionRequest) (*sirenv1beta1.DeleteSubscriptionResponse, error) {
	err := s.subscriptionService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, s.generateRPCErr(err)
	}
	return &sirenv1beta1.DeleteSubscriptionResponse{}, nil
}

func getReceiverMetadataFromDomainObject(item *subscription.Receiver) sirenv1beta1.ReceiverMetadata {
	configMap := make(map[string]string)
	for k, v := range item.Configuration {
		configMap[k] = fmt.Sprintf("%v", v)
	}

	return sirenv1beta1.ReceiverMetadata{
		Id:            item.ID,
		Configuration: configMap,
	}
}

func getReceiverMetadataInDomainObject(item *sirenv1beta1.ReceiverMetadata) subscription.Receiver {
	configMapInterface := make(map[string]interface{})
	for k, v := range item.Configuration {
		configMapInterface[k] = v
	}

	return subscription.Receiver{
		ID:            item.Id,
		Configuration: configMapInterface,
	}
}

func getReceiverMetadataListInDomainObject(domainReceivers []*sirenv1beta1.ReceiverMetadata) []subscription.Receiver {
	receivers := make([]subscription.Receiver, 0)
	for _, receiverMetadataItem := range domainReceivers {
		receivers = append(receivers, getReceiverMetadataInDomainObject(receiverMetadataItem))
	}
	return receivers
}

func getReceiverMetadataListFromDomainObject(domainReceivers []subscription.Receiver) []*sirenv1beta1.ReceiverMetadata {
	receivers := make([]*sirenv1beta1.ReceiverMetadata, 0)
	for _, receiverMetadataItem := range domainReceivers {
		item := getReceiverMetadataFromDomainObject(&receiverMetadataItem)
		receivers = append(receivers, &item)
	}
	return receivers
}
