package v1

import (
	"context"
	sirenv1beta1 "github.com/odpf/siren/api/proto/odpf/siren/v1beta1"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/helper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GRPCServer) ListProviders(_ context.Context, req *sirenv1beta1.ListProvidersRequest) (*sirenv1beta1.ListProvidersResponse, error) {
	providers, err := s.container.ProviderService.ListProviders(map[string]interface{}{
		"urn":  req.GetUrn(),
		"type": req.GetType(),
	})
	if err != nil {
		return nil, helper.GRPCLogError(s.log, codes.Internal, err)
	}

	res := &sirenv1beta1.ListProvidersResponse{
		Providers: make([]*sirenv1beta1.Provider, 0),
	}
	for _, provider := range providers {
		credentials, err := structpb.NewStruct(provider.Credentials)
		if err != nil {
			s.log.Error("failed to fetch provider credentials", "error", err)
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		item := &sirenv1beta1.Provider{
			Id:          provider.Id,
			Urn:         provider.Urn,
			Host:        provider.Host,
			Type:        provider.Type,
			Name:        provider.Name,
			Credentials: credentials,
			Labels:      provider.Labels,
			CreatedAt:   timestamppb.New(provider.CreatedAt),
			UpdatedAt:   timestamppb.New(provider.UpdatedAt),
		}
		res.Providers = append(res.Providers, item)
	}
	return res, nil
}

func (s *GRPCServer) CreateProvider(_ context.Context, req *sirenv1beta1.CreateProviderRequest) (*sirenv1beta1.Provider, error) {
	provider, err := s.container.ProviderService.CreateProvider(&domain.Provider{
		Host:        req.GetHost(),
		Urn:         req.GetUrn(),
		Name:        req.GetName(),
		Type:        req.GetType(),
		Credentials: req.GetCredentials().AsMap(),
		Labels:      req.GetLabels(),
	})
	if err != nil {
		return nil, helper.GRPCLogError(s.log, codes.Internal, err)
	}

	grpcCredentials, err := structpb.NewStruct(provider.Credentials)
	if err != nil {
		return nil, helper.GRPCLogError(s.log, codes.Internal, err)
	}

	return &sirenv1beta1.Provider{
		Id:          provider.Id,
		Host:        provider.Host,
		Urn:         provider.Urn,
		Name:        provider.Name,
		Type:        provider.Type,
		Credentials: grpcCredentials,
		Labels:      provider.Labels,
		CreatedAt:   timestamppb.New(provider.CreatedAt),
		UpdatedAt:   timestamppb.New(provider.UpdatedAt),
	}, nil
}

func (s *GRPCServer) GetProvider(_ context.Context, req *sirenv1beta1.GetProviderRequest) (*sirenv1beta1.Provider, error) {
	provider, err := s.container.ProviderService.GetProvider(req.GetId())
	if provider == nil {
		return nil, status.Errorf(codes.NotFound, "provider not found")
	}
	if err != nil {
		return nil, helper.GRPCLogError(s.log, codes.Internal, err)
	}

	grpcCredentials, err := structpb.NewStruct(provider.Credentials)
	if err != nil {
		return nil, helper.GRPCLogError(s.log, codes.Internal, err)
	}

	return &sirenv1beta1.Provider{
		Id:          provider.Id,
		Host:        provider.Host,
		Urn:         provider.Urn,
		Name:        provider.Name,
		Type:        provider.Type,
		Credentials: grpcCredentials,
		Labels:      provider.Labels,
		CreatedAt:   timestamppb.New(provider.CreatedAt),
		UpdatedAt:   timestamppb.New(provider.UpdatedAt),
	}, nil
}

func (s *GRPCServer) UpdateProvider(_ context.Context, req *sirenv1beta1.UpdateProviderRequest) (*sirenv1beta1.Provider, error) {
	provider, err := s.container.ProviderService.UpdateProvider(&domain.Provider{
		Id:          req.GetId(),
		Host:        req.GetHost(),
		Name:        req.GetName(),
		Type:        req.GetType(),
		Credentials: req.GetCredentials().AsMap(),
		Labels:      req.GetLabels(),
	})
	if err != nil {
		return nil, helper.GRPCLogError(s.log, codes.Internal, err)
	}

	grpcCredentials, err := structpb.NewStruct(provider.Credentials)
	if err != nil {
		return nil, helper.GRPCLogError(s.log, codes.Internal, err)
	}

	return &sirenv1beta1.Provider{
		Id:          provider.Id,
		Host:        provider.Host,
		Urn:         provider.Urn,
		Name:        provider.Name,
		Type:        provider.Type,
		Credentials: grpcCredentials,
		Labels:      provider.Labels,
		CreatedAt:   timestamppb.New(provider.CreatedAt),
		UpdatedAt:   timestamppb.New(provider.UpdatedAt),
	}, nil
}

func (s *GRPCServer) DeleteProvider(_ context.Context, req *sirenv1beta1.DeleteProviderRequest) (*emptypb.Empty, error) {
	err := s.container.ProviderService.DeleteProvider(uint64(req.GetId()))
	if err != nil {
		return nil, helper.GRPCLogError(s.log, codes.Internal, err)
	}

	return &emptypb.Empty{}, nil
}
