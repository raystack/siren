package v1beta1

import (
	"context"

	"github.com/goto/siren/core/provider"
	sirenv1beta1 "github.com/goto/siren/proto/gotocompany/siren/v1beta1"
)

func (s *GRPCServer) ListProviders(ctx context.Context, req *sirenv1beta1.ListProvidersRequest) (*sirenv1beta1.ListProvidersResponse, error) {
	providers, err := s.providerService.List(ctx, provider.Filter{
		URN:  req.GetUrn(),
		Type: req.GetType(),
	})
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	items := []*sirenv1beta1.Provider{}
	for _, prov := range providers {
		item, err := prov.ToV1beta1Proto()
		if err != nil {
			return nil, s.generateRPCErr(err)
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

	protoProv, err := fetchedProvider.ToV1beta1Proto()
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.GetProviderResponse{
		Provider: protoProv,
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
