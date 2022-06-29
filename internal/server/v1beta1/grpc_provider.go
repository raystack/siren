package v1beta1

import (
	"context"
	"fmt"

	"github.com/odpf/siren/core/provider"
	sirenv1beta1 "github.com/odpf/siren/internal/server/proto/odpf/siren/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate mockery --name=ProviderService -r --case underscore --with-expecter --structname ProviderService --filename provider_service.go --output=./mocks
type ProviderService interface {
	ListProviders(map[string]interface{}) ([]*provider.Provider, error)
	CreateProvider(*provider.Provider) (*provider.Provider, error)
	GetProvider(uint64) (*provider.Provider, error)
	UpdateProvider(*provider.Provider) (*provider.Provider, error)
	DeleteProvider(uint64) error
}

func (s *GRPCServer) ListProviders(_ context.Context, req *sirenv1beta1.ListProvidersRequest) (*sirenv1beta1.ListProvidersResponse, error) {
	providers, err := s.providerService.ListProviders(map[string]interface{}{
		"urn":  req.GetUrn(),
		"type": req.GetType(),
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

func (s *GRPCServer) CreateProvider(_ context.Context, req *sirenv1beta1.CreateProviderRequest) (*sirenv1beta1.CreateProviderResponse, error) {
	prov, err := s.providerService.CreateProvider(&provider.Provider{
		Host:        req.GetHost(),
		URN:         req.GetUrn(),
		Name:        req.GetName(),
		Type:        req.GetType(),
		Credentials: req.GetCredentials().AsMap(),
		Labels:      req.GetLabels(),
	})
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.CreateProviderResponse{
		Id: prov.ID,
	}, nil
}

func (s *GRPCServer) GetProvider(_ context.Context, req *sirenv1beta1.GetProviderRequest) (*sirenv1beta1.GetProviderResponse, error) {
	provider, err := s.providerService.GetProvider(req.GetId())
	if provider == nil {
		return nil, status.Errorf(codes.NotFound, "provider not found")
	}
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	grpcCredentials, err := structpb.NewStruct(provider.Credentials)
	if err != nil {
		return nil, s.generateRPCErr(fmt.Errorf("failed to fetch provider credentials: %w", err))
	}

	return &sirenv1beta1.GetProviderResponse{
		Provider: &sirenv1beta1.Provider{
			Id:          provider.ID,
			Host:        provider.Host,
			Urn:         provider.URN,
			Name:        provider.Name,
			Type:        provider.Type,
			Credentials: grpcCredentials,
			Labels:      provider.Labels,
			CreatedAt:   timestamppb.New(provider.CreatedAt),
			UpdatedAt:   timestamppb.New(provider.UpdatedAt),
		},
	}, nil
}

func (s *GRPCServer) UpdateProvider(_ context.Context, req *sirenv1beta1.UpdateProviderRequest) (*sirenv1beta1.UpdateProviderResponse, error) {
	provider, err := s.providerService.UpdateProvider(&provider.Provider{
		ID:          req.GetId(),
		Host:        req.GetHost(),
		Name:        req.GetName(),
		Type:        req.GetType(),
		Credentials: req.GetCredentials().AsMap(),
		Labels:      req.GetLabels(),
	})
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.UpdateProviderResponse{
		Id: provider.ID,
	}, nil
}

func (s *GRPCServer) DeleteProvider(_ context.Context, req *sirenv1beta1.DeleteProviderRequest) (*sirenv1beta1.DeleteProviderResponse, error) {
	err := s.providerService.DeleteProvider(uint64(req.GetId()))
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.DeleteProviderResponse{}, nil
}
