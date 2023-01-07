package v1beta1

import (
	"context"

	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/internal/api"
	"github.com/odpf/siren/pkg/errors"
	sirenv1beta1 "github.com/odpf/siren/proto/odpf/siren/v1beta1"
)

const notificationAPIScope = "notification_api"

func (s *GRPCServer) NotifyReceiver(ctx context.Context, req *sirenv1beta1.NotifyReceiverRequest) (*sirenv1beta1.NotifyReceiverResponse, error) {
	var (
		idempotentID uint64
		err          error
	)

	payloadMap := req.GetPayload().AsMap()

	idempotencyKey := api.GetHeaderString(ctx, s.headers.IdempotencyKey)
	if idempotencyKey != "" {
		idempotentID, err = s.notificationService.CheckAndInsertIdempotency(ctx, notificationAPIScope, idempotencyKey)
		if err != nil {
			// idempotent
			if errors.Is(err, errors.ErrConflict) {
				return &sirenv1beta1.NotifyReceiverResponse{}, nil
			}
			return nil, s.generateRPCErr(err)
		}
	}

	n, err := notification.BuildTypeReceiver(req.GetId(), payloadMap)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	if err := s.notificationService.Dispatch(ctx, n); err != nil {
		return nil, s.generateRPCErr(err)
	}

	if idempotencyKey != "" {
		if err := s.notificationService.MarkIdempotencyAsSuccess(ctx, idempotentID); err != nil {
			return nil, s.generateRPCErr(err)
		}
	}

	return &sirenv1beta1.NotifyReceiverResponse{}, nil
}
