package subscription

import (
	"context"

	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/pkg/errors"
)

func (s *Service) SyncBatchToUpstream(
	ctx context.Context,
	ns *namespace.Namespace,
	pluginService ProviderPlugin) error {

	// fetch all subscriptions in this namespace.
	subscriptionsInNamespace, err := s.repository.List(ctx, Filter{
		NamespaceID: ns.ID,
	})
	if err != nil {
		return err
	}
	receiversMap, err := CreateReceiversMap(ctx, s.receiverService, subscriptionsInNamespace)
	if err != nil {
		return err
	}

	subscriptionsInNamespace, err = AssignReceivers(s.receiverService, receiversMap, subscriptionsInNamespace)
	if err != nil {
		return err
	}
	if err := pluginService.SyncSubscriptions(ctx, subscriptionsInNamespace, ns.URN); err != nil {
		return err
	}

	return nil
}

func CreateReceiversMap(ctx context.Context, receiverService ReceiverService, subscriptions []Subscription) (map[uint64]*receiver.Receiver, error) {
	receiversMap := map[uint64]*receiver.Receiver{}
	for _, subs := range subscriptions {
		for _, rcv := range subs.Receivers {
			if rcv.ID != 0 {
				receiversMap[rcv.ID] = nil
			}
		}
	}

	// empty receivers map
	if len(receiversMap) == 0 {
		return nil, errors.New("no receivers found in subscription")
	}

	listOfReceiverIDs := []uint64{}
	for k := range receiversMap {
		listOfReceiverIDs = append(listOfReceiverIDs, k)
	}

	filteredReceivers, err := receiverService.List(ctx, receiver.Filter{
		ReceiverIDs: listOfReceiverIDs,
	})
	if err != nil {
		return nil, err
	}

	for i, rcv := range filteredReceivers {
		receiversMap[rcv.ID] = &filteredReceivers[i]
	}

	nilReceivers := []uint64{}
	for id, rcv := range receiversMap {
		if rcv == nil {
			nilReceivers = append(nilReceivers, id)
			continue
		}
	}

	if len(nilReceivers) > 0 {
		return nil, errors.ErrInvalid.WithMsgf("receiver id %v don't exist", nilReceivers)
	}

	return receiversMap, nil
}

func AssignReceivers(receiverService ReceiverService, receiversMap map[uint64]*receiver.Receiver, subscriptions []Subscription) ([]Subscription, error) {
	for is := range subscriptions {
		for ir, subsRcv := range subscriptions[is].Receivers {
			if mappedRcv := receiversMap[subsRcv.ID]; mappedRcv == nil {
				return nil, errors.ErrInvalid.WithMsgf("receiver id %d not found", subsRcv.ID)
			}
			subsConfig, err := receiverService.EnrichSubscriptionConfig(subsRcv.Configuration, receiversMap[subsRcv.ID])
			if err != nil {
				return nil, errors.ErrInvalid.WithMsgf(err.Error())
			}
			subscriptions[is].Receivers[ir].ID = receiversMap[subsRcv.ID].ID
			subscriptions[is].Receivers[ir].Type = receiversMap[subsRcv.ID].Type
			subscriptions[is].Receivers[ir].Configuration = subsConfig
		}
	}

	return subscriptions, nil
}
