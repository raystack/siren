package namespace

import (
	"encoding/json"

	"github.com/odpf/siren/pkg/errors"
)

// Service handles business logic
type Service struct {
	repository   Repository
	cryptoClient Encryptor
}

// NewService returns secure service struct
func NewService(cryptoClient Encryptor, repository Repository) *Service {
	return &Service{
		repository:   repository,
		cryptoClient: cryptoClient,
	}
}

func (s *Service) ListNamespaces() ([]*Namespace, error) {
	encrytpedNamespaces, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	namespaces := make([]*Namespace, 0, len(encrytpedNamespaces))
	for _, en := range encrytpedNamespaces {
		ns, err := s.decrypt(en)
		if err != nil {
			return nil, err
		}
		namespaces = append(namespaces, ns)
	}
	return namespaces, nil
}

func (s *Service) CreateNamespace(ns *Namespace) error {
	if ns == nil {
		return errors.ErrInvalid.WithCausef("namespace is nil").WithMsgf("incoming namespace is empty")
	}
	encryptedNamespace, err := s.encrypt(ns)
	if err != nil {
		return err
	}

	if err := s.repository.Create(encryptedNamespace); err != nil {
		if errors.Is(err, ErrDuplicate) {
			return errors.ErrConflict.WithMsgf(err.Error())
		}
		return err
	}

	encryptedNamespace.Namespace.Credentials = ns.Credentials
	*ns = *encryptedNamespace.Namespace
	return nil
}

func (s *Service) GetNamespace(id uint64) (*Namespace, error) {
	encryptedNamespace, err := s.repository.Get(id)
	if err != nil {
		if errors.As(err, new(NotFoundError)) {
			return nil, errors.ErrNotFound.WithMsgf(err.Error())
		}
		return nil, err
	}
	if encryptedNamespace == nil {
		return nil, nil
	}

	return s.decrypt(encryptedNamespace)
}

func (s *Service) UpdateNamespace(namespace *Namespace) error {
	encryptedNamespace, err := s.encrypt(namespace)
	if err != nil {
		return err
	}

	if err := s.repository.Update(encryptedNamespace); err != nil {
		if errors.Is(err, ErrDuplicate) {
			return errors.ErrConflict.WithMsgf(err.Error())
		}
		if errors.As(err, new(NotFoundError)) {
			return errors.ErrNotFound.WithMsgf(err.Error())
		}
		return err
	}

	encryptedNamespace.Namespace.Credentials = namespace.Credentials
	*namespace = *encryptedNamespace.Namespace
	return nil
}

func (s *Service) DeleteNamespace(id uint64) error {
	return s.repository.Delete(id)
}

func (s *Service) encrypt(ns *Namespace) (*EncryptedNamespace, error) {
	plainTextCredentials, err := json.Marshal(ns.Credentials)
	if err != nil {
		return nil, err
	}

	encryptedCredentials, err := s.cryptoClient.Encrypt(string(plainTextCredentials))
	if err != nil {
		return nil, err
	}

	return &EncryptedNamespace{
		Namespace:   ns,
		Credentials: encryptedCredentials,
	}, nil
}

func (s *Service) decrypt(ens *EncryptedNamespace) (*Namespace, error) {
	decryptedCredentialsStr, err := s.cryptoClient.Decrypt(ens.Credentials)
	if err != nil {
		return nil, err
	}

	var decryptedCredentials map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedCredentialsStr), &decryptedCredentials); err != nil {
		return nil, err

	}

	ens.Namespace.Credentials = decryptedCredentials
	return ens.Namespace, nil
}
