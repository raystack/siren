package slack

import (
	"encoding/json"
	"fmt"

	goslack "github.com/slack-go/slack"
)

// TODO support block-kit messages
type Message struct {
	Channel     string              `yaml:"channel,omitempty" json:"channel,omitempty"  mapstructure:"channel"`
	Text        string              `yaml:"text,omitempty" json:"text,omitempty"  mapstructure:"text"`
	Username    string              `yaml:"username,omitempty" json:"username,omitempty"  mapstructure:"username"`
	IconEmoji   string              `yaml:"icon_emoji,omitempty" json:"icon_emoji,omitempty" mapstructure:"icon_emoji"`
	IconURL     string              `yaml:"icon_url,omitempty" json:"icon_url,omitempty"  mapstructure:"icon_url"`
	LinkNames   bool                `yaml:"link_names,omitempty" json:"link_names,omitempty"  mapstructure:"link_names"`
	Attachments []MessageAttachment `yaml:"attachments,omitempty" json:"attachments,omitempty" mapstructure:"attachments"`
}

func (m Message) BuildGoSlackMessageOptions() ([]goslack.MsgOption, error) {
	goslackAttachments := []goslack.Attachment{}
	for _, a := range m.Attachments {
		attachment, err := a.ToGoSlack()
		if err != nil {
			return nil, fmt.Errorf("failed to parse slack attachment: %w", err)
		}
		goslackAttachments = append(goslackAttachments, *attachment)
	}

	msgOptions := []goslack.MsgOption{}

	if m.Text != "" {
		msgOptions = append(msgOptions, goslack.MsgOptionText(m.Text, false))
	}

	if m.Username != "" {
		msgOptions = append(msgOptions, goslack.MsgOptionUsername(m.Username))
	}

	if m.IconEmoji != "" {
		msgOptions = append(msgOptions, goslack.MsgOptionIconEmoji(m.IconEmoji))
	}

	if m.IconURL != "" {
		msgOptions = append(msgOptions, goslack.MsgOptionIconURL(m.IconURL))
	}

	if len(goslackAttachments) != 0 {
		msgOptions = append(msgOptions, goslack.MsgOptionAttachments(goslackAttachments...))
	}

	return msgOptions, nil
}

type MessageAttachment map[string]any

func (ma MessageAttachment) ToGoSlack() (*goslack.Attachment, error) {
	// TODO might want to use more performant JSON marshaller
	// Can't use mapstructure here because goslack.Attachment
	// structure is conflicted in its Blocks field
	gaBlob, err := json.Marshal(ma)
	if err != nil {
		return nil, err
	}

	ga := &goslack.Attachment{}
	if err := json.Unmarshal(gaBlob, &ga); err != nil {
		return nil, err
	}

	return ga, nil
}
