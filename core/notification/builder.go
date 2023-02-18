package notification

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/pkg/errors"
)

// Transform alerts and populate Data and Labels to be interpolated to the system-default template
// .Data
// - id
// - status "FIRING"/"RESOLVED"
// - resource
// - template
// - metric_value
// - metric_name
// - generator_url
// - num_alerts_firing
// - dashboard
// - playbook
// - summary
// .Labels
// - severity "WARNING"/"CRITICAL"
// - alertname
// - (others labels defined in rules)
func BuildFromAlerts(
	alerts []alert.Alert,
	firingLen int,
	createdTime time.Time,
) ([]Notification, error) {
	if len(alerts) == 0 {
		return nil, errors.New("empty alerts")
	}

	alertsMap, err := groupByLabels(alerts)
	if err != nil {
		return nil, err
	}

	var notifications []Notification

	for hashKey, groupedAlerts := range alertsMap {
		sampleAlert := groupedAlerts[0]

		data := map[string]interface{}{}

		mergedAnnotations := map[string][]string{}
		for _, a := range groupedAlerts {
			for k, v := range a.Annotations {
				mergedAnnotations[k] = append(mergedAnnotations[k], v)
			}
		}
		// make unique
		for k, v := range mergedAnnotations {
			mergedAnnotations[k] = removeDuplicateStringValues(v)
		}
		// render annotations
		for k, vSlice := range mergedAnnotations {
			for _, v := range vSlice {
				if _, ok := data[k]; ok {
					data[k] = fmt.Sprintf("%s\n%s", data[k], v)
				} else {
					data[k] = v
				}
			}
		}

		data["status"] = sampleAlert.Status
		data["generator_url"] = sampleAlert.GeneratorURL
		data["num_alerts_firing"] = firingLen

		alertIDs := []int64{}

		for _, a := range groupedAlerts {
			alertIDs = append(alertIDs, int64(a.ID))
		}

		notifications = append(notifications, Notification{
			NamespaceID: sampleAlert.NamespaceID,
			Type:        TypeSubscriber,
			Data:        data,
			Labels:      sampleAlert.Labels,
			Template:    template.ReservedName_SystemDefault,
			UniqueKey:   hashGroupKey(sampleAlert.GroupKey, hashKey),
			CreatedAt:   createdTime,
			AlertIDs:    alertIDs,
		})
	}

	return notifications, nil
}

// BuildTypeReceiver builds a notification struct with receiver type flow
func BuildTypeReceiver(receiverID uint64, payloadMap map[string]interface{}) (Notification, error) {
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
