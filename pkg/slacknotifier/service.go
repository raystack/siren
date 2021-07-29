package slacknotifier

import (
	"fmt"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/codeexchange"
	"github.com/pkg/errors"
	"gorm.io/gorm"
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
		return res, err
	}
	res.OK = true
	return res, nil
}

func NewService(db *gorm.DB, encryptionKey string) (domain.SlackNotifierService, error) {
	svc, err := codeexchange.NewService(db, nil, domain.SlackApp{}, encryptionKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init slack notifier service")
	}
	return &Service{client: NewSlackNotifierClient(), codeExchangeService: svc}, nil
}
