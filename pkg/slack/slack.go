package slack

import (
	goslack "github.com/slack-go/slack"
)

const (
	oAuthServerEndpoint = "https://slack.com/api/oauth.v2.access"

	TypeReceiverChannel = "channel"
	TypeReceiverUser    = "user"
)

//go:generate mockery --name=GoSlackCaller -r --case underscore --with-expecter --structname GoSlackCaller --filename goslack_caller.go --output=./mocks
type GoSlackCaller interface {
	GetConversationsForUser(params *goslack.GetConversationsForUserParameters) (channels []goslack.Channel, nextCursor string, err error)
	GetUserByEmail(email string) (*goslack.User, error)
	SendMessage(channel string, options ...goslack.MsgOption) (string, string, string, error)
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
