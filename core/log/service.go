package log

import "context"

type Service struct {
	notificationLogRepo NotificationLogRepository
}

func NewService(nlr NotificationLogRepository) *Service {
	return &Service{nlr}
}

func (s *Service) LogNotifications(ctx context.Context, nlogs ...Notification) error {
	return s.notificationLogRepo.BulkCreate(ctx, nlogs)
}

func (s *Service) ListAlertIDsBySilenceID(ctx context.Context, silenceID string) ([]int64, error) {
	return s.notificationLogRepo.ListAlertIDsBySilenceID(ctx, silenceID)
}

func (s *Service) ListSubscriptionIDsBySilenceID(ctx context.Context, silenceID string) ([]int64, error) {
	return s.notificationLogRepo.ListSubscriptionIDsBySilenceID(ctx, silenceID)
}
