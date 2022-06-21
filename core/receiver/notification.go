package receiver

import (
	"encoding/json"
	"fmt"

	"github.com/odpf/siren/pkg/slack"
)

// NotificationMessage is an abstraction of receiver's notification message
type NotificationMessage map[string]interface{}

// ToSlackMessage
// {
// 	"receiver_name": "",
// 	"receiver_type": "",
// 	"message": "",
// 	"blocks": [
//			{
// 				"": ""
// 			}
//		]
// }
func (nm NotificationMessage) ToSlackMessage() (*slack.Message, error) {
	jsonByte, err := json.Marshal(nm)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal notification message: %w", err)
	}

	sm := &slack.Message{}
	if err := json.Unmarshal(jsonByte, sm); err != nil {
		return nil, fmt.Errorf("unable to unmarshal notification message byte to slack message: %w", err)
	}

	if err := sm.Validate(); err != nil {
		return nil, err
	}

	return sm, nil
}
