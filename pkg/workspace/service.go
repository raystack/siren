package workspace

import (
	"fmt"
	"github.com/odpf/siren/domain"
	"github.com/pkg/errors"
)

type Service struct {
	client              SlackRepository
	codeExchangeService domain.CodeExchangeService
}

func NewService(service domain.CodeExchangeService) (domain.WorkspaceService, error) {
	return &Service{client: NewRepository(), codeExchangeService: service}, nil
}

func (s Service) GetChannels(workspace string) ([]domain.Channel, error) {
	token, err := s.codeExchangeService.GetToken(workspace)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not get token for workspce: %s", workspace))
	}
	channels, err := s.client.GetWorkspaceChannel(token)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not get channels for workspace: %s", workspace))
	}

	result := make([]domain.Channel, 0)
	for _, c := range channels {
		result = append(result, c.toDomain())
	}

	return result, nil
}
