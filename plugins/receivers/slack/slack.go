package slack

import (
	"context"

	"github.com/goto/siren/pkg/secret"
)

type Encryptor interface {
	Encrypt(str secret.MaskableString) (secret.MaskableString, error)
	Decrypt(str secret.MaskableString) (secret.MaskableString, error)
}

type SlackCaller interface {
	ExchangeAuth(ctx context.Context, authCode, clientID, clientSecret string) (Credential, error)
	GetWorkspaceChannels(ctx context.Context, token secret.MaskableString) ([]Channel, error)
	Notify(ctx context.Context, conf NotificationConfig, message Message) error
}
