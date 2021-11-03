package slacknotifier

import (
	"fmt"
	"github.com/odpf/siren/domain"
	"github.com/pkg/errors"
)

type Service struct {
	client              SlackNotifier
}

func (s Service) Notify(message *domain.SlackMessage) (*domain.SlackMessageSendResponse, error) {
	payload := &SlackMessage{}
	payload = payload.fromDomain(message)
	res := &domain.SlackMessageSendResponse{
		OK: false,
	}
	err := s.client.Notify(payload, message.Token)
	if err != nil {
		return res, errors.Wrap(err, fmt.Sprintf("could not send notification"))
	}
	res.OK = true
	return res, nil
}

func NewService() domain.SlackNotifierService {
	return &Service{
		client: NewSlackNotifierClient(),
	}
}
