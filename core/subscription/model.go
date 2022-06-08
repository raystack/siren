package subscription

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/odpf/siren/internal/store/model"
	"github.com/pkg/errors"
)

type StringStringMap map[string]string

func (m *StringStringMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed type assertion to []byte")
	}
	return json.Unmarshal(b, &m)
}

func (a StringStringMap) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

type ReceiverMetadata struct {
	Id            uint64            `json:"id"`
	Configuration map[string]string `json:"configuration"`
}

type EnrichedReceiverMetadata struct {
	Id            uint64            `json:"id"`
	Type          string            `json:"type"`
	Configuration map[string]string `json:"configuration"`
}

type EnrichedReceiverMetadataList []EnrichedReceiverMetadata

type SubscriptionEnrichedWithReceivers struct {
	Id          uint64
	Namespace   *model.Namespace
	NamespaceId uint64
	Urn         string
	Receiver    EnrichedReceiverMetadataList
	Match       StringStringMap
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
