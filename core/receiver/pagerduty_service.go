package receiver

import (
	"context"

	"github.com/odpf/siren/pkg/errors"
)

type PagerDutyService struct{}

// NewPagerDutyService returns slack service struct
func NewPagerDutyService() *PagerDutyService {
	return &PagerDutyService{}
}

func (s *PagerDutyService) Notify(ctx context.Context, rcv *Receiver, payloadMessage NotificationMessage) error {
	return ErrNotImplemented
}

func (s *PagerDutyService) Encrypt(r *Receiver) error {
	return nil
}

func (s *PagerDutyService) Decrypt(r *Receiver) error {
	return nil
}

func (s *PagerDutyService) PopulateReceiver(ctx context.Context, rcv *Receiver) (*Receiver, error) {
	return rcv, nil
}

func (s *PagerDutyService) ValidateConfiguration(rcv *Receiver) error {
	if rcv == nil {
		return errors.New("receiver to validate is nil")
	}
	_, err := rcv.Configurations.GetString("service_key")
	if err != nil {
		return err
	}

	return nil
}

func (s *PagerDutyService) GetSubscriptionConfig(subsConfs map[string]string, receiverConfs Configurations) (map[string]string, error) {
	mapConf := make(map[string]string)
	if val, ok := receiverConfs["service_key"]; ok {
		if mapConf["service_key"], ok = val.(string); !ok {
			return nil, errors.New("service_key config from receiver should be in string")
		}
	}
	return mapConf, nil
}
