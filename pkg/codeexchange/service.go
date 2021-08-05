package codeexchange

import (
	"github.com/odpf/siren/domain"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CodeExchangeHTTPResponse struct {
	AccessToken string `json:"access_token"`
	Team        struct {
		Name string `json:"name"`
	} `json:"team"`
	Ok bool `json:"ok"`
}

// Service handles business logic
type Service struct {
	repository   ExchangeRepository
	exchanger    Exchanger
	clientID     string
	clientSecret string
}

// NewService returns repository struct
func NewService(db *gorm.DB, httpClient Doer, slackAppConfig domain.SlackApp, encryptionKey string) (domain.CodeExchangeService, error) {
	repository, err := NewRepository(db, encryptionKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create codeexchange repository")
	}

	return &Service{
		repository:   repository,
		clientID:     slackAppConfig.ClientID,
		clientSecret: slackAppConfig.ClientSecret,
		exchanger:    NewSlackClient(httpClient),
	}, nil
}

func (service Service) GetToken(workspace string) (string, error) {
	return service.repository.Get(workspace)
}

func (service Service) Exchange(payload domain.OAuthPayload) (*domain.OAuthExchangeResponse, error) {
	response, err := service.exchanger.Exchange(payload.Code, service.clientID, service.clientSecret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exchange code with slack OAuth server")
	}

	workspaceName := payload.Workspace
	if workspaceName == "" {
		// we can also use fetch domain name instead of Team.Name if needed
		// reference: https://api.slack.com/methods/team.info
		workspaceName = response.Team.Name
	}

	err = service.repository.Upsert(&AccessToken{
		AccessToken: response.AccessToken,
		Workspace:   workspaceName,
	})
	if err != nil {
		return nil, err
	}
	return &domain.OAuthExchangeResponse{OK: true}, nil
}

func (service Service) Migrate() error {
	return service.repository.Migrate()
}
