package slack

import (
	"context"

	goslack "github.com/slack-go/slack"
)

const (
	oAuthServerEndpoint = "https://slack.com/api/oauth.v2.access"

	TypeReceiverChannel = "channel"
	TypeReceiverUser    = "user"
)

//go:generate mockery --name=GoSlackCaller -r --case underscore --with-expecter --structname GoSlackCaller --filename goslack_caller.go --output=./mocks
type GoSlackCaller interface {
	GetConversationsForUserContext(ctx context.Context, params *goslack.GetConversationsForUserParameters) (channels []goslack.Channel, nextCursor string, err error)
	GetUserByEmailContext(ctx context.Context, email string) (*goslack.User, error)
	SendMessageContext(ctx context.Context, channel string, options ...goslack.MsgOption) (string, string, string, error)
}

type codeExchangeHTTPResponse struct {
	AccessToken string `json:"access_token"`
	Team        struct {
		Name string `json:"name"`
	} `json:"team"`
	Ok bool `json:"ok"`
}

type Credential struct {
	AccessToken string
	TeamName    string
}

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
