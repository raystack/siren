package slack

import (
	"errors"

	goslack "github.com/slack-go/slack"
	"gopkg.in/go-playground/validator.v9"
)

const oAuthServerEndpoint = "https://slack.com/api/oauth.v2.access"

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

type Message struct {
	ReceiverName string         `json:"receiver_name" validate:"required"`
	ReceiverType string         `json:"receiver_type" validate:"required,oneof=user channel"`
	Token        string         `json:"token" validate:"required"`
	Message      string         `json:"message"`
	Blocks       goslack.Blocks `json:"blocks"`
}

func (sm *Message) Validate() error {
	v := validator.New()
	if sm.Message == "" && len(sm.Blocks.BlockSet) == 0 {
		return errors.New("non empty message or non zero length block is required")
	}
	return v.Struct(sm)
}
