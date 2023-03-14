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

func groupByLabels(alerts []alert.Alert) (map[uint64][]alert.Alert, error) {
	var alertsMap = map[uint64][]alert.Alert{}

	for _, a := range alerts {
		hash, err := hashstructure.Hash(a.Labels, hashstructure.FormatV2, nil)
		if err != nil {
			return nil, fmt.Errorf("cannot get hash from alert %v", a)
		}
		alertsMap[hash] = append(alertsMap[hash], a)
	}

	return alertsMap, nil
}

// hashGroupKey hash groupKey from alert and hashKey from labels
func hashGroupKey(groupKey string, hashKey uint64) string {
	h := sha256.New()
	// hash.Hash.Write never returns an error.
	//nolint: errcheck
	h.Write([]byte(fmt.Sprintf("%s%d", groupKey, hashKey)))
	return fmt.Sprintf("%x", h.Sum(nil))
}
