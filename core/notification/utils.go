package notification

import (
	"crypto/sha256"
	"fmt"

	"github.com/goto/siren/core/alert"
	"github.com/mitchellh/hashstructure/v2"
)

func removeDuplicateStringValues(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, v := range strSlice {
		if _, value := keys[v]; !value {
			keys[v] = true
			list = append(list, v)
		}
	}
	return list
}

func groupByLabels(alerts []alert.Alert, groupBy []string) (map[uint64][]alert.Alert, error) {
	var alertsMap = map[uint64][]alert.Alert{}

	for _, a := range alerts {
		var groupLabels = buildGroupLabels(a.Labels, groupBy)
		if len(groupLabels) == 0 {
			groupLabels = a.Labels
		}
		hash, err := hashstructure.Hash(groupLabels, hashstructure.FormatV2, nil)
		if err != nil {
			return nil, fmt.Errorf("cannot get hash from alert %v", a)
		}
		alertsMap[hash] = append(alertsMap[hash], a)
	}

	return alertsMap, nil
}

func buildGroupLabels(alertLabels map[string]string, groupBy []string) map[string]string {
	var groupLabels = map[string]string{}

	for _, g := range groupBy {
		if v, ok := alertLabels[g]; ok {
			groupLabels[g] = v
		}
	}

	return groupLabels
}

// hashGroupKey hash groupKey from alert and hashKey from labels
func hashGroupKey(groupKey string, hashKey uint64) string {
	h := sha256.New()
	// hash.Hash.Write never returns an error.
	//nolint: errcheck
	h.Write([]byte(fmt.Sprintf("%s%d", groupKey, hashKey)))
	return fmt.Sprintf("%x", h.Sum(nil))
}
