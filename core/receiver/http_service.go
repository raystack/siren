package receiver

import (
	"context"

	"github.com/odpf/siren/pkg/errors"
)

type HTTPService struct{}

// NewHTTPService returns slack service struct
func NewHTTPService() *HTTPService {
	return &HTTPService{}
}

func (s *HTTPService) Notify(ctx context.Context, rcv *Receiver, payloadMessage NotificationMessage) error {
	return ErrNotImplemented
}

func (s *HTTPService) Encrypt(r *Receiver) error {
	return nil
}

func (s *HTTPService) Decrypt(r *Receiver) error {
	return nil
}

func (s *HTTPService) PopulateReceiver(ctx context.Context, rcv *Receiver) (*Receiver, error) {
	return rcv, nil
}

func (s *HTTPService) ValidateConfiguration(rcv *Receiver) error {
	if rcv == nil {
		return errors.New("receiver to validate is nil")
	}
	_, err := rcv.Configurations.GetString("url")
	if err != nil {
		return err
	}
	return nil
}

func (s *HTTPService) GetSubscriptionConfig(subsConfs map[string]string, receiverConfs Configurations) (map[string]string, error) {
	mapConf := make(map[string]string)
	if val, ok := receiverConfs["url"]; ok {
		if mapConf["url"], ok = val.(string); !ok {
			return nil, errors.New("url config from receiver should be in string")
		}
	}
	return mapConf, nil
}
