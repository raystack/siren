package slack

import (
	"context"

	"github.com/raystack/siren/pkg/secret"
)

//go:generate mockery --name=Encryptor -r --case underscore --with-expecter --structname Encryptor --filename encryptor.go --output=./mocks
type Encryptor interface {
	Encrypt(str secret.MaskableString) (secret.MaskableString, error)
	Decrypt(str secret.MaskableString) (secret.MaskableString, error)
}

//go:generate mockery --name=SlackCaller -r --case underscore --with-expecter --structname SlackCaller --filename slack_caller.go --output=./mocks
type SlackCaller interface {
	ExchangeAuth(ctx context.Context, authCode, clientID, clientSecret string) (Credential, error)
	GetWorkspaceChannels(ctx context.Context, token secret.MaskableString) ([]Channel, error)
	Notify(ctx context.Context, conf NotificationConfig, message Message) error
}
