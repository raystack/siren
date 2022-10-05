package slack

import "context"

//go:generate mockery --name=Encryptor -r --case underscore --with-expecter --structname Encryptor --filename encryptor.go --output=./mocks
type Encryptor interface {
	Encrypt(str string) (string, error)
	Decrypt(str string) (string, error)
}

//go:generate mockery --name=SlackCaller -r --case underscore --with-expecter --structname SlackCaller --filename slack_caller.go --output=./mocks
type SlackCaller interface {
	ExchangeAuth(ctx context.Context, authCode, clientID, clientSecret string) (Credential, error)
	GetWorkspaceChannels(ctx context.Context, token string) ([]Channel, error)
	Notify(ctx context.Context, conf NotificationConfig, message Message) error
}
