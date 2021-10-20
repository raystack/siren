package namespace

import (
	"bytes"
	"encoding/base64"
	"github.com/gtank/cryptopasta"
	"github.com/odpf/siren/domain"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"io"
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
	repository  NamespaceRepository
	transformer EncryptorDecryptor
}

// NewService returns service struct
func NewService(db *gorm.DB, encryptionKey string) (domain.NamespaceService, error) {
	transformer, err := NewTransformer(encryptionKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transformer")
	}

	return &Service{NewRepository(db), transformer}, nil
}

func (s Service) ListNamespaces() ([]*domain.Namespace, error) {
	namespaces, err := s.repository.List()
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.List")
	}

	domainNamespaces := make([]*domain.Namespace, 0, len(namespaces))
	for i := 0; i < len(namespaces); i++ {
		decryptedCredentials, err := s.transformer.Decrypt(namespaces[i].Credentials)
		if err != nil {
			return nil, errors.Wrap(err, "s.transformer.Decrypt")
		}
		namespaces[i].Credentials = decryptedCredentials
		namespace, err := namespaces[i].toDomain()
		if err != nil {
			return nil, errors.Wrap(err, "namespaces[i].toDomain()")
		}
		domainNamespaces = append(domainNamespaces, namespace)
	}

	return domainNamespaces, nil
}

func (s Service) CreateNamespace(namespace *domain.Namespace) (*domain.Namespace, error) {
	n := &Namespace{}
	_, err := n.fromDomain(namespace)
	if err != nil {
		return nil, errors.Wrap(err, "n.fromDomain()")
	}
	plainTextCredentials := n.Credentials
	encryptedCredentials, err := s.transformer.Encrypt(n.Credentials)
	if err != nil {
		return nil, errors.Wrap(err, "s.transformer.Encrypt")
	}
	n.Credentials = encryptedCredentials
	newNamespace, err := s.repository.Create(n)
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.Create")
	}
	newNamespace.Credentials = plainTextCredentials
	return newNamespace.toDomain()
}

func (s Service) GetNamespace(id uint64) (*domain.Namespace, error) {
	namespace, err := s.repository.Get(id)
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.Get")
	}
	if namespace == nil {
		return nil, nil
	}
	decryptedCredentials, err := s.transformer.Decrypt(namespace.Credentials)
	if err != nil {
		return nil, errors.Wrap(err, "s.transformer.Decrypt")
	}
	namespace.Credentials = decryptedCredentials
	return namespace.toDomain()
}

func (s Service) UpdateNamespace(namespace *domain.Namespace) (*domain.Namespace, error) {
	w := &Namespace{}
	_, err := w.fromDomain(namespace)
	if err != nil {
		return nil, errors.Wrap(err, "n.fromDomain")
	}
	plainTextCredentials := w.Credentials
	encryptedCredentials, err := s.transformer.Encrypt(w.Credentials)
	if err != nil {
		return nil, errors.Wrap(err, "s.transformer.Encrypt")
	}
	w.Credentials = encryptedCredentials
	updatedNamespace, err := s.repository.Update(w)
	if err != nil {
		return nil, errors.Wrap(err, "s.repository.Update")
	}
	updatedNamespace.Credentials = plainTextCredentials
	return updatedNamespace.toDomain()
}

func (s Service) DeleteNamespace(id uint64) error {
	return s.repository.Delete(id)
}

func (s Service) Migrate() error {
	return s.repository.Migrate()
}
