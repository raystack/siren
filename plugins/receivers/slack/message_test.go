package slack

import (
	"testing"
)

func TestMessage_BuildGoSlackMessageOptions(t *testing.T) {
	tests := []struct {
		name    string
		message Message
		wantErr bool
	}{
		{
			name: "should return error if failed to parse message attachment",
			message: Message{
				Attachments: []MessageAttachment{
					{"blocks": "test"},
				},
			},
			wantErr: true,
		},
		{
			name: "should build all message options if all fields in message present",
			message: Message{
				Channel:   "channel", // won't be included
				Text:      "text",
				Username:  "username",
				IconEmoji: ":emoji:",
				IconURL:   "icon_url",
				LinkNames: true, // won't be included
				Attachments: []MessageAttachment{
					{
						"color": "#f2c744",
						"blocks": []map[string]interface{}{
							{
								"type": "section",
								"text": map[string]interface{}{
									"type": "mrkdwn",
									"text": "this is markdown task",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.message.BuildGoSlackMessageOptions()
			if (err != nil) != tt.wantErr {
				t.Errorf("Message.BuildGoSlackMessageOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
