package namespace

import (
	"context"
	"encoding/json"

	"github.com/goto/siren/core/provider"
	"github.com/goto/siren/pkg/errors"
	"github.com/goto/siren/pkg/secret"
)

type ProviderService interface {
	Get(ctx context.Context, id uint64) (*provider.Provider, error)
}

// Service handles business logic
type Service struct {
	repository      Repository
	cryptoClient    Encryptor
	providerService ProviderService
	registry        map[string]ConfigSyncer
}

// NewService returns secure service struct
func NewService(cryptoClient Encryptor, repository Repository, providerService ProviderService, registry map[string]ConfigSyncer) *Service {
	return &Service{
		repository:      repository,
		providerService: providerService,
		cryptoClient:    cryptoClient,
		registry:        registry,
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

	prov, err := s.providerService.Get(ctx, ns.Provider.ID)
	if err != nil {
		return err
	}

	pluginService, err := s.getProviderPluginService(prov.Type)
	if err != nil {
		return err
	}

	encryptedNamespace, err := s.encrypt(ns)
	if err != nil {
		return err
	}

	// this is without transaction to decouple creation in siren and runtime config sync in remote
	// to cover case where the provider plugin needs to validate webhook to siren
	if err = s.repository.Create(ctx, encryptedNamespace); err != nil {
		if errors.Is(err, ErrDuplicate) {
			return errors.ErrConflict.WithMsgf(err.Error())
		}
		if errors.Is(err, ErrRelation) {
			return errors.ErrNotFound.WithMsgf(err.Error())
		}
		return err
	}

	encryptedNamespace.Provider = *prov

	labels, err := pluginService.SyncRuntimeConfig(ctx, encryptedNamespace.ID, encryptedNamespace.URN, encryptedNamespace.Labels, encryptedNamespace.Provider)
	if err != nil {
		return err
	}

	if encryptedNamespace.Labels == nil {
		encryptedNamespace.Labels = make(map[string]string)
	}

	for k, v := range labels {
		encryptedNamespace.Labels[k] = v
	}

	if err = s.repository.UpdateLabels(ctx, encryptedNamespace.ID, encryptedNamespace.Labels); err != nil {
		if errors.As(err, new(NotFoundError)) {
			return errors.ErrNotFound.WithMsgf(err.Error())
		}
		return err
	}

	ns.ID = encryptedNamespace.ID

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

	return s.decrypt(encryptedNamespace)
}

func (s *Service) Update(ctx context.Context, ns *Namespace) error {
	if ns == nil {
		return errors.ErrInvalid.WithCausef("namespace is nil").WithMsgf("incoming namespace is empty")
	}

	encryptedNS, err := s.repository.Get(ctx, ns.ID)
	if err != nil {
		return err
	}

	existingNS, err := s.decrypt(encryptedNS)
	if err != nil {
		return err
	}

	pluginService, err := s.getProviderPluginService(existingNS.Provider.Type)
	if err != nil {
		return err
	}

	// urn is immutable
	ns.URN = existingNS.URN

	encryptedNamespace, err := s.encrypt(ns)
	if err != nil {
		return err
	}

	ctx = s.repository.WithTransaction(ctx)
	// merge existing labels
	for k, v := range encryptedNS.Labels {
		ns.Labels[k] = v
	}

	if err = s.repository.Update(ctx, encryptedNamespace); err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return err
		}
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

	encryptedNamespace.Provider = existingNS.Provider

	labels, err := pluginService.SyncRuntimeConfig(ctx, encryptedNamespace.ID, encryptedNamespace.URN, encryptedNamespace.Labels, encryptedNamespace.Provider)
	if err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return err
		}
		return err
	}

	if encryptedNamespace.Labels == nil {
		encryptedNamespace.Labels = make(map[string]string)
	}

	for k, v := range labels {
		encryptedNamespace.Labels[k] = v
	}

	if err = s.repository.UpdateLabels(ctx, encryptedNamespace.ID, encryptedNamespace.Labels); err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return err
		}
		if errors.As(err, new(NotFoundError)) {
			return errors.ErrNotFound.WithMsgf(err.Error())
		}
		return err
	}

	if err := s.repository.Commit(ctx); err != nil {
		return err
	}

	ns.ID = encryptedNamespace.ID

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

	encryptedCredentials, err := s.cryptoClient.Encrypt(secret.MaskableString(plainTextCredentials))
	if err != nil {
		return nil, err
	}

	return &EncryptedNamespace{
		Namespace:        ns,
		CredentialString: encryptedCredentials.UnmaskedString(),
	}, nil
}

func (s *Service) decrypt(ens *EncryptedNamespace) (*Namespace, error) {
	decryptedCredentialsStr, err := s.cryptoClient.Decrypt(secret.MaskableString(ens.CredentialString))
	if err != nil {
		return nil, err
	}

	var decryptedCredentials map[string]any
	if err := json.Unmarshal([]byte(decryptedCredentialsStr), &decryptedCredentials); err != nil {
		return nil, err

	}

	ens.Namespace.Credentials = decryptedCredentials
	return ens.Namespace, nil
}

func (s *Service) getProviderPluginService(providerType string) (ConfigSyncer, error) {
	pluginService, exist := s.registry[providerType]
	if !exist {
		return nil, errors.ErrInvalid.WithMsgf("unsupported provider type: %q", providerType)
	}
	return pluginService, nil
}
