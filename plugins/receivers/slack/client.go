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
	"github.com/odpf/siren/pkg/httpclient"
	"github.com/odpf/siren/pkg/retry"
	goslack "github.com/slack-go/slack"
)

const (
	defaultSlackAPIHost = "https://slack.com/api"
	oAuthSlackPath      = "/oauth.v2.access"

	TypeChannelChannel = "channel"
	TypeChannelUser    = "user"
	DefaultChannelType = TypeChannelChannel
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

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Credential struct {
	AccessToken string
	TeamName    string
}

type Client struct {
	cfg        AppConfig
	httpClient *httpclient.Client
	// data       *clientData
	retrier retry.Runner
}

// type clientData struct {
// 	// customPath   string
// 	// authCode     string
// 	// clientID     string
// 	// clientSecret string
// 	// token        string
// 	// creds         Credential
// 	goslackClient GoSlackCaller
// }

// NewClient is a constructor to create slack client.
// this version uses go-slack client and this construction wraps the client.
func NewClient(cfg AppConfig, opts ...ClientOption) *Client {
	c := &Client{
		cfg: cfg,
	}
	for _, opt := range opts {
		opt(c)
	}

	if cfg.APIHost == "" {
		c.cfg.APIHost = defaultSlackAPIHost
	}

	// sanitize
	c.cfg.APIHost = c.cfg.APIHost + "/"

	if c.httpClient == nil {
		c.httpClient = httpclient.New(httpclient.Config{})
	}

	return c
}

// // createClient create slack client with 3 options
// // - goslack Client
// // - token
// // - client secret
// // the order that took precedence
// // goslackClient - token - client secret
// // e.g. if user passes goslackClient, it will ignore the others
// func (c *Client) createGoSlackClient(ctx context.Context, opts ...ClientCallOption) (GoSlackCaller, error) {
// 	c.data = &clientData{}
// 	for _, opt := range opts {
// 		opt(c.data)
// 	}

// 	if c.data.goslackClient != nil {
// 		return c.data.goslackClient, nil
// 	}

// 	if c.data.token != "" {
// 		goslackOpts := []goslack.Option{}
// 		if c.data.customPath != "" {
// 			goslackOpts = append(goslackOpts, goslack.OptionAPIURL(c.cfg.APIHost+c.data.customPath))
// 		}

// 		c.data.goslackClient = goslack.New(c.data.token, goslackOpts...)
// 		return c.data.goslackClient, nil
// 	}

// 	// if c.data.authCode == "" || c.data.clientID == "" || c.data.clientSecret == "" {
// 	// 	return nil, errors.New("no client id/secret credential provided")
// 	// }

// 	// creds, err := c.ExchangeAuth(ctx, c.data.authCode, c.data.clientID, c.data.clientSecret)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// c.data.creds = creds

// 	return c.data.goslackClient, nil
// }

// ExchangeAuth submits client ID, client secret, and auth code and retrieve acces token and team name
func (c *Client) ExchangeAuth(ctx context.Context, authCode, clientID, clientSecret string) (Credential, error) {
	data := url.Values{}
	data.Set("code", authCode)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	response := codeExchangeHTTPResponse{}
	req, err := http.NewRequest(http.MethodPost, c.cfg.APIHost+oAuthSlackPath, strings.NewReader(data.Encode()))
	if err != nil {
		return Credential{}, fmt.Errorf("failed to create request body: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.HTTP().Do(req)
	if err != nil {
		return Credential{}, fmt.Errorf("failure in http call: %w", err)
	}
	defer resp.Body.Close()
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

// GetWorkspaceChannels fetches list of joined channel of a client
func (c *Client) GetWorkspaceChannels(ctx context.Context, token string) ([]Channel, error) {
	// gsc, err := c.createGoSlackClient(ctx, opts...)
	gsc := goslack.New(token, goslack.OptionAPIURL(c.cfg.APIHost))
	// if err != nil {
	// 	return nil, fmt.Errorf("goslack client creation failure: %w", err)
	// }

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

// Notify sends message to a specific slack channel
func (c *Client) Notify(ctx context.Context, conf NotificationConfig, message Message) error {
	gsc := goslack.New(conf.ReceiverConfig.Token, goslack.OptionAPIURL(c.cfg.APIHost))

	var channelID string
	switch conf.ChannelType {
	case TypeChannelChannel:
		joinedChannelList, err := c.getJoinedChannelsList(ctx, gsc)
		if err != nil {
			return fmt.Errorf("failed to fetch joined channel list: %w", err)
		}
		channelID = searchChannelId(joinedChannelList, conf.ChannelName)
		if channelID == "" {
			return fmt.Errorf("app is not part of the channel %q", conf.ChannelName)
		}
	case TypeChannelUser:
		// https://api.slack.com/methods/users.lookupByEmail
		user, err := gsc.GetUserByEmailContext(ctx, conf.ChannelName)
		if err != nil {
			if err.Error() == "users_not_found" {
				return fmt.Errorf("failed to get id for %q", conf.ChannelName)
			}
			return err
		}
		channelID = user.ID
	default:
		return fmt.Errorf("unknown receiver type %q", conf.ChannelType)
	}

	msgOptions, err := message.BuildGoSlackMessageOptions()
	if err != nil {
		return err
	}

	if c.retrier != nil {
		if err := c.retrier.Run(ctx, func(ctx context.Context) error {
			return c.sendMessageContext(ctx, gsc, channelID, conf.ChannelName, msgOptions...)
		}); err != nil {
			return err
		}
		return nil
	} else {
		return c.sendMessageContext(ctx, gsc, channelID, conf.ChannelName, msgOptions...)
	}
}

func (c *Client) sendMessageContext(ctx context.Context, gsc GoSlackCaller, channelID string, channelName string, msgOpts ...goslack.MsgOption) error {
	_, _, _, err := gsc.SendMessageContext(
		ctx,
		channelID,
		msgOpts...,
	)
	if err != nil {
		// if 429 or 5xx do retry
		var scErr goslack.StatusCodeError
		isit := errors.As(err, &scErr)
		if isit {
			if scErr.Retryable() {
				return retry.RetryableError{Err: err}
			}
		}
		var rlErr *goslack.RateLimitedError
		if errors.As(err, &rlErr) {
			if rlErr.Retryable() {
				return retry.RetryableError{Err: err}
			}
		}
		return fmt.Errorf("failed to send message to %q: %w", channelName, err)
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
