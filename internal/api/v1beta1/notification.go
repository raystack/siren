package v1beta1

import (
	"context"

	"github.com/mitchellh/mapstructure"
	"github.com/odpf/siren/core/notification"
	sirenv1beta1 "github.com/odpf/siren/proto/odpf/siren/v1beta1"
)

// TODO will be refactored later
func (s *GRPCServer) NotifyReceiver(ctx context.Context, req *sirenv1beta1.NotifyReceiverRequest) (*sirenv1beta1.NotifyReceiverResponse, error) {
	payloadMap := req.GetPayload().AsMap()

	n := &notification.Notification{}
	err := mapstructure.Decode(payloadMap, n)
	if err != nil {
		s.logger.Warn("failed to parse payload to notification", "payload", payloadMap)
	}

	if err == nil {
		if err := s.notificationService.Dispatch(ctx, *n); err != nil {
			s.logger.Warn("failed to send to notification service", "api", "notification", "notification", n, "error", err)
		}
	}

	return &sirenv1beta1.NotifyReceiverResponse{}, nil
}
