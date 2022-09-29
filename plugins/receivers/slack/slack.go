package slack

import "context"

//go:generate mockery --name=Encryptor -r --case underscore --with-expecter --structname Encryptor --filename encryptor.go --output=./mocks
type Encryptor interface {
	Encrypt(str string) (string, error)
	Decrypt(str string) (string, error)
}

//go:generate mockery --name=SlackClient -r --case underscore --with-expecter --structname SlackClient --filename slack_client.go --output=./mocks
type SlackClient interface {
	ExchangeAuth(ctx context.Context, authCode, clientID, clientSecret string) (Credential, error)
	GetWorkspaceChannels(ctx context.Context, opts ...ClientCallOption) ([]Channel, error)
	Notify(ctx context.Context, message *MessageGoSlack, opts ...ClientCallOption) error
}
