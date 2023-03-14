package subscription

import (
	context "context"
	"fmt"
	"time"

	"github.com/goto/siren/core/silence"
)

//go:generate mockery --name=Repository -r --case underscore --with-expecter --structname SubscriptionRepository --filename subscription_repository.go --output=./mocks
type Repository interface {
	List(context.Context, Filter) ([]Subscription, error)
	Create(context.Context, *Subscription) error
	Get(context.Context, uint64) (*Subscription, error)
	Update(context.Context, *Subscription) error
	Delete(context.Context, uint64) error
}

type Receiver struct {
	ID            uint64                 `json:"id"`
	Configuration map[string]interface{} `json:"configuration"`

	// Type won't be exposed to the end-user, this is used to add more details for notification purposes
	Type string
}

type Subscription struct {
	ID        uint64            `json:"id"`
	URN       string            `json:"urn"`
	Namespace uint64            `json:"namespace"`
	Receivers []Receiver        `json:"receivers"`
	Match     map[string]string `json:"match"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

func (s Subscription) ReceiversAsMap() map[uint64]Receiver {
	var m = make(map[uint64]Receiver)
	for _, rcv := range s.Receivers {
		m[rcv.ID] = rcv
	}
	return m
}

func (s Subscription) SilenceReceivers(silences []silence.Silence) (map[uint64][]silence.Silence, []Receiver, error) {
	var (
		nonSilencedReceiversMap = map[uint64]Receiver{}
		silencedReceiversMap    = map[uint64][]silence.Silence{}
	)

	if len(silences) == 0 {
		return nil, s.Receivers, nil
	}

	// evaluate all receivers of subscribers with all matched silences
	for _, sil := range silences {
		for _, rcv := range s.Receivers {
			isSilenced, err := sil.EvaluateSubscriptionRule(rcv)
			if err != nil {
				return nil, nil, fmt.Errorf("error evaluating subscription receiver %v: %w", rcv, err)
			}

			if isSilenced {
				if len(silencedReceiversMap) == 0 {
					silencedReceiversMap = make(map[uint64][]silence.Silence)
				}
				silencedReceiversMap[rcv.ID] = append(silencedReceiversMap[rcv.ID], sil)
			} else {
				nonSilencedReceiversMap[rcv.ID] = rcv
			}
		}
	}

	var nonSilencedReceivers []Receiver
	for k, v := range nonSilencedReceiversMap {
		// remove if non silenced receivers are part of silenced receivers
		if _, ok := silencedReceiversMap[k]; !ok {
			nonSilencedReceivers = append(nonSilencedReceivers, v)
		}
	}

	return silencedReceiversMap, nonSilencedReceivers, nil
}
