package receiver

import (
	"context"

	"github.com/odpf/siren/pkg/slack"
)

//go:generate mockery --name=SlackClient -r --case underscore --with-expecter --structname SlackClient --filename slack_client.go --output=./mocks
type SlackClient interface {
	GetWorkspaceChannels(ctx context.Context, opts ...slack.ClientCallOption) ([]slack.Channel, error)
	Notify(ctx context.Context, message *slack.Message, opts ...slack.ClientCallOption) error
}
