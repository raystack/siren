package v1beta1

import (
	"context"
	"fmt"

	"github.com/odpf/siren/core/provider"
	sirenv1beta1 "github.com/odpf/siren/proto/odpf/siren/v1beta1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate mockery --name=ProviderService -r --case underscore --with-expecter --structname ProviderService --filename provider_service.go --output=./mocks
type ProviderService interface {
	List(context.Context, provider.Filter) ([]provider.Provider, error)
	Create(context.Context, *provider.Provider) error
	Get(context.Context, uint64) (*provider.Provider, error)
	Update(context.Context, *provider.Provider) error
	Delete(context.Context, uint64) error
}

func (s *GRPCServer) ListProviders(ctx context.Context, req *sirenv1beta1.ListProvidersRequest) (*sirenv1beta1.ListProvidersResponse, error) {
	providers, err := s.providerService.List(ctx, provider.Filter{
		URN:  req.GetUrn(),
		Type: req.GetType(),
	})
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	items := []*sirenv1beta1.Provider{}
	for _, provider := range providers {
		credentials, err := structpb.NewStruct(provider.Credentials)
		if err != nil {
			return nil, s.generateRPCErr(fmt.Errorf("failed to fetch provider credentials: %w", err))
		}

		item := &sirenv1beta1.Provider{
			Id:          provider.ID,
			Urn:         provider.URN,
			Host:        provider.Host,
			Type:        provider.Type,
			Name:        provider.Name,
			Credentials: credentials,
			Labels:      provider.Labels,
			CreatedAt:   timestamppb.New(provider.CreatedAt),
			UpdatedAt:   timestamppb.New(provider.UpdatedAt),
		}
		items = append(items, item)
	}
	return &sirenv1beta1.ListProvidersResponse{
		Providers: items,
	}, nil
}

func (s *GRPCServer) CreateProvider(ctx context.Context, req *sirenv1beta1.CreateProviderRequest) (*sirenv1beta1.CreateProviderResponse, error) {
	prv := &provider.Provider{
		Host:        req.GetHost(),
		URN:         req.GetUrn(),
		Name:        req.GetName(),
		Type:        req.GetType(),
		Credentials: req.GetCredentials().AsMap(),
		Labels:      req.GetLabels(),
	}

	if err := s.providerService.Create(ctx, prv); err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.CreateProviderResponse{
		Id: prv.ID,
	}, nil
}

func (s *GRPCServer) GetProvider(ctx context.Context, req *sirenv1beta1.GetProviderRequest) (*sirenv1beta1.GetProviderResponse, error) {
	fetchedProvider, err := s.providerService.Get(ctx, req.GetId())
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	grpcCredentials, err := structpb.NewStruct(fetchedProvider.Credentials)
	if err != nil {
		return nil, s.generateRPCErr(fmt.Errorf("failed to fetch provider credentials: %w", err))
	}

	return &sirenv1beta1.GetProviderResponse{
		Provider: &sirenv1beta1.Provider{
			Id:          fetchedProvider.ID,
			Host:        fetchedProvider.Host,
			Urn:         fetchedProvider.URN,
			Name:        fetchedProvider.Name,
			Type:        fetchedProvider.Type,
			Credentials: grpcCredentials,
			Labels:      fetchedProvider.Labels,
			CreatedAt:   timestamppb.New(fetchedProvider.CreatedAt),
			UpdatedAt:   timestamppb.New(fetchedProvider.UpdatedAt),
		},
	}, nil
}

func (s *GRPCServer) UpdateProvider(ctx context.Context, req *sirenv1beta1.UpdateProviderRequest) (*sirenv1beta1.UpdateProviderResponse, error) {
	prv := &provider.Provider{
		ID:          req.GetId(),
		Host:        req.GetHost(),
		Name:        req.GetName(),
		Type:        req.GetType(),
		Credentials: req.GetCredentials().AsMap(),
		Labels:      req.GetLabels(),
	}

	if err := s.providerService.Update(ctx, prv); err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.UpdateProviderResponse{
		Id: prv.ID,
	}, nil
}

func (s *GRPCServer) DeleteProvider(ctx context.Context, req *sirenv1beta1.DeleteProviderRequest) (*sirenv1beta1.DeleteProviderResponse, error) {
	err := s.providerService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.DeleteProviderResponse{}, nil
}
