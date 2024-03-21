package v1beta1

import (
	"context"
	"fmt"

	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/core/template"
	"github.com/goto/siren/internal/api"
	"github.com/goto/siren/pkg/errors"
	sirenv1beta1 "github.com/goto/siren/proto/gotocompany/siren/v1beta1"
)

const notificationAPIScope = "notification_api"

func (s *GRPCServer) PostNotification(ctx context.Context, req *sirenv1beta1.PostNotificationRequest) (*sirenv1beta1.PostNotificationResponse, error) {
	idempotencyScope := api.GetHeaderString(ctx, s.headers.IdempotencyScope)
	if idempotencyScope == "" {
		idempotencyScope = notificationAPIScope
	}

	idempotencyKey := api.GetHeaderString(ctx, s.headers.IdempotencyKey)
	if idempotencyKey != "" {
		if notificationID, err := s.notificationService.CheckIdempotency(ctx, idempotencyScope, idempotencyKey); notificationID != "" {
			return &sirenv1beta1.PostNotificationResponse{
				NotificationId: notificationID,
			}, nil
		} else if errors.Is(err, errors.ErrNotFound) {
			s.logger.Debug("no idempotency found with detail", "scope", idempotencyScope, "key", idempotencyKey)
		} else {
			return nil, s.generateRPCErr(fmt.Errorf("error when checking idempotency: %w", err))
		}
	}

	var receiverSelectors = []map[string]string{}
	for _, pbSelector := range req.GetReceivers() {
		var mss = make(map[string]string)
		for k, v := range pbSelector.AsMap() {
			vString, ok := v.(string)
			if !ok {
				err := errors.ErrInvalid.WithMsgf("invalid receiver selectors, value must be string but found %v", v)
				return nil, s.generateRPCErr(err)
			}
			mss[k] = vString
		}
		receiverSelectors = append(receiverSelectors, mss)
	}

	// TODO once custom template is supported, this needs to be set
	var notificationTemplate = template.ReservedName_SystemDefault
	if req.GetTemplate() != "" {
		notificationTemplate = req.GetTemplate()
	}

	n := notification.Notification{
		Type:              notification.TypeEvent,
		Data:              req.GetData().AsMap(),
		Labels:            req.GetLabels(),
		Template:          notificationTemplate,
		ReceiverSelectors: receiverSelectors,
	}

	notificationID, err := s.notificationService.Dispatch(ctx, n)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	if idempotencyKey != "" {
		if err := s.notificationService.InsertIdempotency(ctx, idempotencyScope, idempotencyKey, notificationID); err != nil {
			return nil, s.generateRPCErr(err)
		}
	}

	return &sirenv1beta1.PostNotificationResponse{
		NotificationId: notificationID,
	}, nil
}
