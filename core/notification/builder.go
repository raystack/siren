package notification

import (
	"fmt"
	"time"

	"github.com/goto/siren/pkg/errors"
	"github.com/mitchellh/mapstructure"
)

// BuildTypeReceiver builds a notification struct with receiver type flow
func BuildTypeReceiver(receiverID uint64, payloadMap map[string]any) (Notification, error) {
	n := Notification{}
	if err := mapstructure.Decode(payloadMap, &n); err != nil {
		return Notification{}, errors.ErrInvalid.WithMsgf("failed to parse payload to notification: %s", err.Error())
	}

	if val, ok := payloadMap[ValidDurationRequestKey]; ok {
		valString, ok := val.(string)
		if !ok {
			return Notification{}, fmt.Errorf("cannot parse %s value: %v", ValidDurationRequestKey, val)
		}
		parsedDur, err := time.ParseDuration(valString)
		if err != nil {
			return Notification{}, err
		}
		n.ValidDuration = parsedDur
	}

	n.Type = TypeReceiver

	if len(n.Labels) == 0 {
		n.Labels = map[string]string{}
	}

	n.Labels[ReceiverIDLabelKey] = fmt.Sprintf("%d", receiverID)

	return n, nil
}
