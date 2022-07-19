package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/odpf/siren/pkg/errors"
	goslack "github.com/slack-go/slack"
)

type Client struct {
	httpClient *http.Client
	data       *clientData
}

type clientData struct {
	authCode      string
	clientID      string
	clientSecret  string
	token         string
	creds         Credential
	goslackClient GoSlackCaller
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{}
	for _, opt := range opts {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}

	return c
}

// createClient create slack client with 3 options
// goslack Client, token / client secret
// the order that took precedence
// goslackClient - token - client secret
// e.g. if user passes goslackClient, it will ignore the others
func (c *Client) createGoSlackClient(ctx context.Context, opts ...ClientCallOption) (GoSlackCaller, error) {
	c.data = &clientData{}
	for _, opt := range opts {
		opt(c.data)
	}

	if c.data.goslackClient != nil {
		return c.data.goslackClient, nil
	}

	if c.data.token != "" {
		c.data.goslackClient = goslack.New(c.data.token)
		return c.data.goslackClient, nil
	}

	if c.data.authCode == "" || c.data.clientID == "" || c.data.clientSecret == "" {
		return nil, errors.New("no client id/secret credential provided")
	}

	creds, err := c.ExchangeAuth(ctx, c.data.authCode, c.data.clientID, c.data.clientSecret)
	if err != nil {
		return nil, err
	}
	c.data.creds = creds

	return c.data.goslackClient, nil
}

func (c *Client) ExchangeAuth(ctx context.Context, authCode, clientID, clientSecret string) (Credential, error) {
	data := url.Values{}
	data.Set("code", authCode)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	response := codeExchangeHTTPResponse{}
	req, err := http.NewRequest(http.MethodPost, oAuthServerEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return Credential{}, fmt.Errorf("failed to create request body: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Credential{}, fmt.Errorf("failure in http call: %w", err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Credential{}, fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return Credential{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	if !response.Ok {
		return Credential{}, errors.New("slack oauth call failed")
	}
	return Credential{
		AccessToken: response.AccessToken,
		TeamName:    response.Team.Name,
	}, nil
}

func (c *Client) GetWorkspaceChannels(ctx context.Context, opts ...ClientCallOption) ([]Channel, error) {
	gsc, err := c.createGoSlackClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("goslack client creation failure: %w", err)
	}

	joinedChannelList, err := c.getJoinedChannelsList(ctx, gsc)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch joined channel list: %w", err)
	}

	result := make([]Channel, 0)
	for _, c := range joinedChannelList {
		result = append(result, Channel{
			ID:   c.ID,
			Name: c.Name,
		})
	}
	return result, nil
}

func (c *Client) Notify(ctx context.Context, message *Message, opts ...ClientCallOption) error {
	gsc, err := c.createGoSlackClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("goslack client creation failure: %w", err)
	}

	var channelID string
	switch message.ReceiverType {
	case TypeReceiverChannel:
		joinedChannelList, err := c.getJoinedChannelsList(ctx, gsc)
		if err != nil {
			return fmt.Errorf("failed to fetch joined channel list: %w", err)
		}
		channelID = searchChannelId(joinedChannelList, message.ReceiverName)
		if channelID == "" {
			return fmt.Errorf("app is not part of the channel %q", message.ReceiverName)
		}
	case TypeReceiverUser:
		user, err := gsc.GetUserByEmailContext(ctx, message.ReceiverName)
		if err != nil {
			if err.Error() == "users_not_found" {
				return fmt.Errorf("failed to get id for %q", message.ReceiverName)
			}
			return err
		}
		channelID = user.ID
	default:
		return fmt.Errorf("unknown receiver type %q", message.ReceiverType)
	}
	_, _, _, err = gsc.SendMessageContext(ctx, channelID, goslack.MsgOptionText(message.Message, false), goslack.MsgOptionBlocks(message.Blocks.BlockSet...))
	if err != nil {
		return fmt.Errorf("failed to send message to %q", message.ReceiverName)
	}
	return nil
}

func (c *Client) getJoinedChannelsList(ctx context.Context, gsc GoSlackCaller) ([]goslack.Channel, error) {
	channelList := make([]goslack.Channel, 0)
	curr := ""
	for {
		channels, nextCursor, err := gsc.GetConversationsForUserContext(ctx, &goslack.GetConversationsForUserParameters{
			Types:  []string{"public_channel", "private_channel"},
			Cursor: curr,
			Limit:  1000})
		if err != nil {
			return channelList, err
		}
		channelList = append(channelList, channels...)
		curr = nextCursor
		if curr == "" {
			break
		}
	}
	return channelList, nil
}

func searchChannelId(channels []goslack.Channel, channelName string) string {
	for _, c := range channels {
		if c.Name == channelName {
			return c.ID
		}
	}
	return ""
}
