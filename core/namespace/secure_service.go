package namespace

import (
	"encoding/json"
	"fmt"
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

func (ss *SecureService) ListNamespaces() ([]*Namespace, error) {
	encrytpedNamespaces, err := ss.repository.List()
	if err != nil {
		return nil, fmt.Errorf("secureService.repository.List: %w", err)
	}

	namespaces := make([]*Namespace, 0, len(encrytpedNamespaces))
	for _, en := range encrytpedNamespaces {
		ns, err := ss.decrypt(en)
		if err != nil {
			return nil, err
		}
		namespaces = append(namespaces, ns)
	}
	return namespaces, nil
}

func (ss *SecureService) CreateNamespace(namespace *Namespace) error {
	//TODO check if namespace is nil
	encryptedNamespace, err := ss.encrypt(namespace)
	if err != nil {
		return err
	}

	if err := ss.repository.Create(encryptedNamespace); err != nil {
		return fmt.Errorf("secureService.repository.Create: %w", err)
	}

	encryptedNamespace.Namespace.Credentials = namespace.Credentials
	*namespace = *encryptedNamespace.Namespace

	return nil
}

func (ss *SecureService) GetNamespace(id uint64) (*Namespace, error) {
	encryptedNamespace, err := ss.repository.Get(id)
	if err != nil {
		return nil, fmt.Errorf("secureService.repository.Get: %w", err)
	}
	if encryptedNamespace == nil {
		return nil, nil
	}

	return ss.decrypt(encryptedNamespace)
}

func (ss *SecureService) UpdateNamespace(namespace *Namespace) error {
	encryptedNamespace, err := ss.encrypt(namespace)
	if err != nil {
		return err
	}

	if err := ss.repository.Update(encryptedNamespace); err != nil {
		return fmt.Errorf("secureService.repository.Update: %w", err)
	}

	encryptedNamespace.Namespace.Credentials = namespace.Credentials
	*namespace = *encryptedNamespace.Namespace
	return nil
}

func (ss *SecureService) DeleteNamespace(id uint64) error {
	return ss.repository.Delete(id)
}

func (ss *SecureService) Migrate() error {
	return ss.repository.Migrate()
}

func (ss *SecureService) encrypt(ns *Namespace) (*EncryptedNamespace, error) {
	plainTextCredentials, err := json.Marshal(ns.Credentials)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	encryptedCredentials, err := ss.cryptoClient.Encrypt(string(plainTextCredentials))
	if err != nil {
		return nil, fmt.Errorf("secureService.cryptoClient.Encrypt: %w", err)
	}

	return &EncryptedNamespace{
		Namespace:   ns,
		Credentials: encryptedCredentials,
	}, nil
}

func (ss *SecureService) decrypt(ens *EncryptedNamespace) (*Namespace, error) {
	decryptedCredentialsStr, err := ss.cryptoClient.Decrypt(ens.Credentials)
	if err != nil {
		return nil, fmt.Errorf("secureService.cryptoClient.Decrypt: %w", err)
	}

	var decryptedCredentials map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedCredentialsStr), &decryptedCredentials); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)

	}

	ens.Namespace.Credentials = decryptedCredentials
	return ens.Namespace, nil
}
