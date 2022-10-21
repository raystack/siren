package slack_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/odpf/siren/pkg/retry"
	"github.com/odpf/siren/pkg/secret"
	"github.com/odpf/siren/plugins/receivers/slack"
	goslack "github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetWorkspaceChannels(t *testing.T) {
	var token = secret.MaskableString("test-token")

	t.Run("return error when failed to fetch joined channel list", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
		}))

		c := slack.NewClient(slack.AppConfig{APIHost: testServer.URL})
		channels, err := c.GetWorkspaceChannels(context.Background(), token)

		assert.EqualError(t, err, "failed to fetch joined channel list: slack server error: 502 Bad Gateway")
		assert.Empty(t, channels)

		testServer.Close()
	})

	t.Run("return channels when GetWorkspaceChannels succeed", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respStruct := struct {
				Channels []goslack.Channel `json:"channels"`
			}{
				Channels: []goslack.Channel{
					{
						GroupConversation: goslack.GroupConversation{
							Conversation: goslack.Conversation{
								ID: "123",
							},
							Name: "test",
						},
						IsChannel: true,
					},
				},
			}

			respByte, _ := json.Marshal(respStruct)

			w.Write(respByte)
		}))

		c := slack.NewClient(slack.AppConfig{APIHost: testServer.URL})
		channels, err := c.GetWorkspaceChannels(context.Background(), token)

		assert.NoError(t, err)
		assert.Equal(t, []slack.Channel{{
			ID:   "123",
			Name: "test",
		}}, channels)

		testServer.Close()
	})
}

func TestClient_NotifyChannel(t *testing.T) {
	var token = secret.MaskableString("test-token")

	t.Run("return error when message receiver type is wrong", func(t *testing.T) {
		c := slack.NewClient(slack.AppConfig{})
		err := c.Notify(context.Background(),
			slack.NotificationConfig{
				ReceiverConfig: slack.ReceiverConfig{
					Token: token,
				},
			},
			slack.Message{})

		assert.EqualError(t, err, "unknown receiver type \"\"")
	})

	t.Run("return error when failed to fetch joined channel list", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
		}))

		c := slack.NewClient(slack.AppConfig{APIHost: testServer.URL})
		err := c.Notify(
			context.Background(),
			slack.NotificationConfig{
				ReceiverConfig: slack.ReceiverConfig{
					Token: token,
				},
				SubscriptionConfig: slack.SubscriptionConfig{
					ChannelType: slack.TypeChannelChannel,
				},
			},
			slack.Message{})

		assert.EqualError(t, err, "failed to fetch joined channel list: slack server error: 502 Bad Gateway")

		testServer.Close()
	})

	t.Run("return error when app is not part of the channel", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users.conversations" {
				respStruct := struct {
					Channels []goslack.Channel `json:"channels"`
				}{
					Channels: []goslack.Channel{
						{
							GroupConversation: goslack.GroupConversation{
								Conversation: goslack.Conversation{
									ID: "123",
								},
								Name: "test",
							},
							IsChannel: true,
						},
					},
				}

				respByte, _ := json.Marshal(respStruct)

				w.Write(respByte)
				return
			}
		}))

		c := slack.NewClient(slack.AppConfig{APIHost: testServer.URL})
		err := c.Notify(
			context.Background(),
			slack.NotificationConfig{
				ReceiverConfig: slack.ReceiverConfig{
					Token: token,
				},
				SubscriptionConfig: slack.SubscriptionConfig{
					ChannelType: slack.TypeChannelChannel,
				},
			},
			slack.Message{})

		assert.EqualError(t, err, "app is not part of the channel \"\"")

		testServer.Close()
	})

	t.Run("return error when app is not part of the channel", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users.conversations" {
				respStruct := struct {
					Channels []goslack.Channel `json:"channels"`
				}{
					Channels: []goslack.Channel{
						{
							GroupConversation: goslack.GroupConversation{
								Conversation: goslack.Conversation{
									ID: "123",
								},
								Name: "test",
							},
							IsChannel: true,
						},
					},
				}

				respByte, _ := json.Marshal(respStruct)

				w.Write(respByte)
				return
			}
		}))

		c := slack.NewClient(slack.AppConfig{APIHost: testServer.URL})
		err := c.Notify(
			context.Background(),
			slack.NotificationConfig{
				ReceiverConfig: slack.ReceiverConfig{
					Token: token,
				},
				SubscriptionConfig: slack.SubscriptionConfig{
					ChannelType: slack.TypeChannelChannel,
				},
			},
			slack.Message{})

		assert.EqualError(t, err, "app is not part of the channel \"\"")

		testServer.Close()
	})

	t.Run("return nil error when notify is succeed through channel", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users.conversations" {
				respStruct := struct {
					Channels []goslack.Channel `json:"channels"`
				}{
					Channels: []goslack.Channel{
						{
							GroupConversation: goslack.GroupConversation{
								Conversation: goslack.Conversation{
									ID: "123",
								},
								Name: "test",
							},
							IsChannel: true,
						},
					},
				}

				respByte, _ := json.Marshal(respStruct)
				w.Write(respByte)
				return
			} else {
				w.Write([]byte(`{"ok":true}`))
			}
		}))

		c := slack.NewClient(slack.AppConfig{APIHost: testServer.URL})
		err := c.Notify(
			context.Background(),
			slack.NotificationConfig{
				ReceiverConfig: slack.ReceiverConfig{
					Token: token,
				},
				SubscriptionConfig: slack.SubscriptionConfig{
					ChannelType: slack.TypeChannelChannel,
				},
			},
			slack.Message{
				Channel: "test",
			})

		assert.NoError(t, err)

		testServer.Close()
	})
}

func TestClient_NotifyUser(t *testing.T) {
	var token = secret.MaskableString("test-token")

	t.Run("return error when failed to get user for an email", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"ok": false,"error": "users_not_found"}`))
		}))

		c := slack.NewClient(slack.AppConfig{APIHost: testServer.URL})
		err := c.Notify(
			context.Background(),
			slack.NotificationConfig{
				ReceiverConfig: slack.ReceiverConfig{
					Token: token,
				},
				SubscriptionConfig: slack.SubscriptionConfig{
					ChannelType: slack.TypeChannelUser,
				},
			},
			slack.Message{})

		assert.EqualError(t, err, "failed to get id for \"\"")

		testServer.Close()
	})

	t.Run("return error when GetUserByEmailContext return error", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
		}))

		c := slack.NewClient(slack.AppConfig{APIHost: testServer.URL})
		err := c.Notify(
			context.Background(),
			slack.NotificationConfig{
				ReceiverConfig: slack.ReceiverConfig{
					Token: token,
				},
				SubscriptionConfig: slack.SubscriptionConfig{
					ChannelType: slack.TypeChannelUser,
				},
			},
			slack.Message{})

		assert.EqualError(t, err, "slack server error: 502 Bad Gateway")

		testServer.Close()
	})

	t.Run("return nil error when notify is succeed through user", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users.lookupByEmail" {
				w.Write([]byte(`{"ok": true,"user": {"id": "123123","name":"email@email.com"}}`))
			} else {
				w.Write([]byte(`{"ok":true}`))
			}
		}))

		c := slack.NewClient(slack.AppConfig{APIHost: testServer.URL})
		err := c.Notify(
			context.Background(),
			slack.NotificationConfig{
				ReceiverConfig: slack.ReceiverConfig{
					Token: token,
				},
				SubscriptionConfig: slack.SubscriptionConfig{
					ChannelType: slack.TypeChannelUser,
					ChannelName: "email@email.com",
				},
			},
			slack.Message{})

		assert.NoError(t, err)

		testServer.Close()
	})
}

func TestClient_NotifyWithRetrier(t *testing.T) {
	var (
		expectedCounter = 4
		token           = secret.MaskableString("test-token")
	)

	t.Run("when 429 is returned", func(t *testing.T) {
		var counter = 0
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users.lookupByEmail" {
				w.Write([]byte(`{"ok": true,"user": {"id": "123123","name":"email@email.com"}}`))
			} else {
				counter = counter + 1
				w.Header().Set("Retry-After", "10")
				w.WriteHeader(http.StatusTooManyRequests)

				w.Write([]byte(`{"ok":false}`))
			}
		}))

		c := slack.NewClient(slack.AppConfig{APIHost: testServer.URL}, slack.ClientWithRetrier(retry.New(retry.Config{Enable: true})))
		_ = c.Notify(
			context.Background(),
			slack.NotificationConfig{
				ReceiverConfig: slack.ReceiverConfig{
					Token: token,
				},
				SubscriptionConfig: slack.SubscriptionConfig{
					ChannelType: slack.TypeChannelUser,
					ChannelName: "email@email.com",
				},
			},
			slack.Message{})

		assert.Equal(t, expectedCounter, counter)

		testServer.Close()
	})

	t.Run("when 5xx is returned", func(t *testing.T) {
		var counter = 0
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/users.lookupByEmail" {
				w.Write([]byte(`{"ok": true,"user": {"id": "123123","name":"email@email.com"}}`))
			} else {
				counter = counter + 1
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"ok":false}`))
			}
		}))

		c := slack.NewClient(slack.AppConfig{APIHost: testServer.URL}, slack.ClientWithRetrier(retry.New(retry.Config{Enable: true})))
		_ = c.Notify(
			context.Background(),
			slack.NotificationConfig{
				ReceiverConfig: slack.ReceiverConfig{
					Token: token,
				},
				SubscriptionConfig: slack.SubscriptionConfig{
					ChannelType: slack.TypeChannelUser,
					ChannelName: "email@email.com",
				},
			},
			slack.Message{})

		assert.Equal(t, expectedCounter, counter)

		testServer.Close()
	})

}
