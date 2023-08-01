package notification

import (
	"context"
	"strconv"

	"github.com/goto/siren/core/log"
	"github.com/goto/siren/pkg/errors"
)

type DispatchReceiverService struct {
	receiverService ReceiverService
	notifierPlugins map[string]Notifier
}

func NewDispatchReceiverService(receiverService ReceiverService, notifierPlugins map[string]Notifier) *DispatchReceiverService {
	return &DispatchReceiverService{
		receiverService: receiverService,
		notifierPlugins: notifierPlugins,
	}
}

func (s *DispatchReceiverService) getNotifierPlugin(receiverType string) (Notifier, error) {
	notifierPlugin, exist := s.notifierPlugins[receiverType]
	if !exist {
		return nil, errors.ErrInvalid.WithMsgf("unsupported receiver type: %q", receiverType)
	}
	return notifierPlugin, nil
}

func (s *DispatchReceiverService) PrepareMessage(ctx context.Context, n Notification) ([]Message, []log.Notification, bool, error) {

	var notificationLogs []log.Notification

	receiverID, err := strconv.ParseUint(n.Labels[ReceiverIDLabelKey], 0, 64)
	if err != nil {
		// should not goes here as this already have been checked
		return nil, nil, false, err
	}

	rcv, err := s.receiverService.Get(ctx, receiverID)
	if err != nil {
		return nil, nil, false, err
	}

	notifierPlugin, err := s.getNotifierPlugin(rcv.Type)
	if err != nil {
		return nil, nil, false, errors.ErrInvalid.WithMsgf("invalid receiver type: %s", err.Error())
	}

	message, err := InitMessage(
		ctx,
		notifierPlugin,
		n,
		rcv.Type,
		rcv.Configurations,
		InitWithExpiryDuration(n.ValidDuration),
	)
	if err != nil {
		return nil, nil, false, err
	}

	messages := []Message{message}
	notificationLogs = append(notificationLogs, log.Notification{
		NamespaceID:    n.NamespaceID,
		NotificationID: n.ID,
		ReceiverID:     rcv.ID,
		AlertIDs:       n.AlertIDs,
	})

	return messages, notificationLogs, false, nil
}
