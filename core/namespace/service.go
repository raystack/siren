package namespace

import (
	"context"
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

func (s *Service) List(ctx context.Context) ([]Namespace, error) {
	encrytpedNamespaces, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}

	namespaces := make([]Namespace, 0, len(encrytpedNamespaces))
	for _, en := range encrytpedNamespaces {
		ns, err := s.decrypt(&en)
		if err != nil {
			return nil, err
		}
		namespaces = append(namespaces, *ns)
	}
	return namespaces, nil
}

func (s *Service) Create(ctx context.Context, ns *Namespace) error {
	if ns == nil {
		return errors.ErrInvalid.WithCausef("namespace is nil").WithMsgf("incoming namespace is empty")
	}
	encryptedNamespace, err := s.encrypt(ns)
	if err != nil {
		return err
	}

	err = s.repository.Create(ctx, encryptedNamespace)
	if err != nil {
		if errors.Is(err, ErrDuplicate) {
			return errors.ErrConflict.WithMsgf(err.Error())
		}
		if errors.Is(err, ErrRelation) {
			return errors.ErrNotFound.WithMsgf(err.Error())
		}
		return err
	}

	return nil
}

func (s *Service) Get(ctx context.Context, id uint64) (*Namespace, error) {
	encryptedNamespace, err := s.repository.Get(ctx, id)
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

func (s *Service) Update(ctx context.Context, namespace *Namespace) error {
	encryptedNamespace, err := s.encrypt(namespace)
	if err != nil {
		return err
	}

	err = s.repository.Update(ctx, encryptedNamespace)
	if err != nil {
		if errors.Is(err, ErrDuplicate) {
			return errors.ErrConflict.WithMsgf(err.Error())
		}
		if errors.Is(err, ErrRelation) {
			return errors.ErrNotFound.WithMsgf(err.Error())
		}
		if errors.As(err, new(NotFoundError)) {
			return errors.ErrNotFound.WithMsgf(err.Error())
		}
		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	return s.repository.Delete(ctx, id)
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
