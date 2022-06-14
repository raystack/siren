package receiver

import (
	"fmt"

	"github.com/pkg/errors"
)

// Service handles business logic
type SecureService struct {
	repository   Repository
	cryptoClient Encryptor
}

// NewSecureService returns secure service struct
func NewSecureService(cryptoClient Encryptor, repository Repository) *SecureService {
	return &SecureService{
		repository:   repository,
		cryptoClient: cryptoClient,
	}
}

func (ss *SecureService) ListReceivers() ([]*Receiver, error) {
	receivers, err := ss.repository.List()
	if err != nil {
		return nil, err
	}

	domainReceivers := make([]*Receiver, 0, len(receivers))
	for i := 0; i < len(receivers); i++ {
		rcv := receivers[i]

		if rcv.Type == Slack {
			if err = ss.postTransform(rcv); err != nil {
				return nil, err
			}
		}

		domainReceivers = append(domainReceivers, rcv)
	}

	return domainReceivers, nil

}

func (ss *SecureService) CreateReceiver(rcv *Receiver) error {
	if rcv.Type == Slack {
		if err := ss.preTransform(rcv); err != nil {
			return err
		}
	}

	if err := ss.repository.Create(rcv); err != nil {
		return err
	}

	if rcv.Type == Slack {
		if err := ss.postTransform(rcv); err != nil {
			return err
		}
	}

	return nil
}

func (ss *SecureService) GetReceiver(id uint64) (*Receiver, error) {
	rcv, err := ss.repository.Get(id)
	if err != nil {
		return nil, err
	}

	if rcv.Type == Slack {
		if err := ss.postTransform(rcv); err != nil {
			return nil, err
		}
	}

	return rcv, nil
}

func (ss *SecureService) UpdateReceiver(rcv *Receiver) error {
	if rcv.Type == Slack {
		if err := ss.preTransform(rcv); err != nil {
			return err
		}
	}

	if err := ss.repository.Update(rcv); err != nil {
		return err
	}

	return nil
}

func (ss *SecureService) DeleteReceiver(id uint64) error {
	return ss.repository.Delete(id)
}

func (ss *SecureService) Migrate() error {
	return ss.repository.Migrate()
}

func (ss *SecureService) preTransform(r *Receiver) error {
	var token string
	var ok bool
	if token, ok = r.Configurations["token"].(string); !ok {
		return errors.New("no token field found")
	}
	chiperText, err := ss.cryptoClient.Encrypt(token)
	if err != nil {
		return fmt.Errorf("pre transform encrypt failed: %w", err)
	}
	r.Configurations["token"] = chiperText

	return nil
}

func (ss *SecureService) postTransform(r *Receiver) error {
	var cipherText string
	var ok bool
	if cipherText, ok = r.Configurations["token"].(string); !ok {
		return errors.New("no token field found")
	}
	token, err := ss.cryptoClient.Decrypt(cipherText)
	if err != nil {
		return fmt.Errorf("post transform decrypt failed: %w", err)
	}
	r.Configurations["token"] = token
	return nil
}
