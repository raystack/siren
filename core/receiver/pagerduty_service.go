package receiver

import "github.com/odpf/siren/pkg/errors"

type PagerDutyService struct{}

// NewPagerDutyService returns slack service struct
func NewPagerDutyService() *PagerDutyService {
	return &PagerDutyService{}
}

func (s *PagerDutyService) Notify(rcv *Receiver, payloadMessage NotificationMessage) error {
	return nil
}

func (s *PagerDutyService) Encrypt(r *Receiver) error {
	return nil
}

func (s *PagerDutyService) Decrypt(r *Receiver) error {
	return nil
}

func (s *PagerDutyService) PopulateReceiver(rcv *Receiver) (*Receiver, error) {
	return rcv, nil
}

func (s *PagerDutyService) ValidateConfiguration(configurations Configurations) error {
	_, err := configurations.GetString("service_key")
	if err != nil {
		return errors.ErrInvalid.WithMsgf(err.Error())
	}

	return nil
}
