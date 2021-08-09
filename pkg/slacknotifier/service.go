package slacknotifier

import (
	"fmt"
	"github.com/odpf/siren/domain"
	"github.com/pkg/errors"
)

type Service struct {
	client              SlackNotifier
	codeExchangeService domain.CodeExchangeService
}

func (s Service) Notify(message *domain.SlackMessage) (*domain.SlackMessageSendResponse, error) {
	m := &SlackMessage{}
	m = m.fromDomain(message)
	token, err := s.codeExchangeService.GetToken(message.Entity)
	res := &domain.SlackMessageSendResponse{
		OK: false,
	}
	if err != nil {
		return res, errors.Wrap(err, fmt.Sprintf("could not get token for entity: %s", message.Entity))
	}
	err = s.client.Notify(m, token)
	if err != nil {
		return res, errors.Wrap(err, fmt.Sprintf("could not send notification"))
	}
	res.OK = true
	return res, nil
}

func NewService(codeExchangeService domain.CodeExchangeService) domain.SlackNotifierService {
	return &Service{
		client:              NewSlackNotifierClient(),
		codeExchangeService: codeExchangeService,
	}
}
