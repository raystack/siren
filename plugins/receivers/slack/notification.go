package slack

import (
	"encoding/json"
	"fmt"
)

// ToSlackMessage
//
//	{
//		"receiver_name": "",
//		"receiver_type": "",
//		"message": "",
//		"blocks": [
//				{
//					"": ""
//				}
//			]
//	}
func GetSlackMessage(payloadMessage map[string]interface{}) (*Message, error) {
	jsonByte, err := json.Marshal(payloadMessage)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal notification message: %w", err)
	}

	sm := &Message{}
	if err := json.Unmarshal(jsonByte, sm); err != nil {
		return nil, fmt.Errorf("unable to unmarshal notification message byte to slack message: %w", err)
	}

	if err := sm.Validate(); err != nil {
		return nil, err
	}

	return sm, nil
}
