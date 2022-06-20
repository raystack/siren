package receiver

import "github.com/odpf/siren/pkg/slack"

//go:generate mockery --name=SlackClient -r --case underscore --with-expecter --structname SlackClient --filename slack_client.go --output=./mocks
type SlackClient interface {
	GetWorkspaceChannels(opts ...slack.ClientCallOption) ([]slack.Channel, error)
	Notify(message *slack.Message, opts ...slack.ClientCallOption) error
}
