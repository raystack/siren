package namespace

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"

	"github.com/gtank/cryptopasta"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/store"
	"github.com/pkg/errors"
)

type EncryptorDecryptor interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

type Transformer struct {
	encryptionKey *[32]byte
}

func NewTransformer(encryptionKey string) (*Transformer, error) {
	secretKey := &[32]byte{}
	if len(encryptionKey) < 32 {
		return nil, errors.New("random hash should be 32 chars in length")
	}
	_, err := io.ReadFull(bytes.NewBufferString(encryptionKey), secretKey[:])
	if err != nil {
		return nil, err
	}

	return &Transformer{
		encryptionKey: secretKey,
	}, nil
}

func (t *Transformer) Encrypt(s string) (string, error) {
	cipher, err := cryptopasta.Encrypt([]byte(s), t.encryptionKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipher), nil
}

func (t *Transformer) Decrypt(s string) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	decryptedToken, err := cryptopasta.Decrypt(encrypted, t.encryptionKey)
	if err != nil {
		return "", err
	}
	return string(decryptedToken), nil
}

// Service handles business logic
type Service struct {
	repository  store.NamespaceRepository
	transformer EncryptorDecryptor
}

// NewService returns service struct
func NewService(repository store.NamespaceRepository, encryptionKey string) (domain.NamespaceService, error) {
	transformer, err := NewTransformer(encryptionKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transformer")
	}

	return &Service{repository, transformer}, nil
}

func (s Service) ListNamespaces() ([]*domain.Namespace, error) {
	encrytpedNamespaces, err := s.repository.List()
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.List")
	}

	namespaces := make([]*domain.Namespace, 0, len(encrytpedNamespaces))
	for _, en := range encrytpedNamespaces {
		ns, err := s.decrypt(en)
		if err != nil {
			return nil, err
		}
		namespaces = append(namespaces, ns)
	}
	return namespaces, nil
}

func (s Service) CreateNamespace(namespace *domain.Namespace) error {
	encryptedNamespace, err := s.encrypt(namespace)
	if err != nil {
		return err
	}

	if err := s.repository.Create(encryptedNamespace); err != nil {
		return errors.Wrap(err, "s.repository.Create")
	}

	encryptedNamespace.Namespace.Credentials = namespace.Credentials
	*namespace = *encryptedNamespace.Namespace

	return nil
}

func (s Service) GetNamespace(id uint64) (*domain.Namespace, error) {
	encryptedNamespace, err := s.repository.Get(id)
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.Get")
	}
	if encryptedNamespace == nil {
		return nil, nil
	}

	return s.decrypt(encryptedNamespace)
}

func (s Service) UpdateNamespace(namespace *domain.Namespace) error {
	encryptedNamespace, err := s.encrypt(namespace)
	if err != nil {
		return err
	}

	if err := s.repository.Update(encryptedNamespace); err != nil {
		return errors.Wrap(err, "s.repository.Update")
	}

	encryptedNamespace.Namespace.Credentials = namespace.Credentials
	*namespace = *encryptedNamespace.Namespace
	return nil
}

func (s Service) DeleteNamespace(id uint64) error {
	return s.repository.Delete(id)
}

func (s Service) Migrate() error {
	return s.repository.Migrate()
}

func (s Service) encrypt(ns *domain.Namespace) (*domain.EncryptedNamespace, error) {
	plainTextCredentials, err := json.Marshal(ns.Credentials)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	encryptedCredentials, err := s.transformer.Encrypt(string(plainTextCredentials))
	if err != nil {
		return nil, errors.Wrap(err, "s.transformer.Encrypt")
	}

	return &domain.EncryptedNamespace{
		Namespace:   ns,
		Credentials: encryptedCredentials,
	}, nil
}

func (s Service) decrypt(ens *domain.EncryptedNamespace) (*domain.Namespace, error) {
	decryptedCredentialsStr, err := s.transformer.Decrypt(ens.Credentials)
	if err != nil {
		return nil, errors.Wrap(err, "s.transformer.Decrypt")
	}

	var decryptedCredentials map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedCredentialsStr), &decryptedCredentials); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}

	ens.Namespace.Credentials = decryptedCredentials
	return ens.Namespace, nil
}
