package v1

import (
	"context"
	sirenv1 "github.com/odpf/siren/api/proto/odpf/siren/v1"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/helper"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GRPCServer) ListProviders(_ context.Context, _ *emptypb.Empty) (*sirenv1.ListProvidersResponse, error) {
	providers, err := s.container.ProviderService.ListProviders()
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	res := &sirenv1.ListProvidersResponse{
		Providers: make([]*sirenv1.Provider, 0),
	}
	for _, provider := range providers {
		credentials, err := structpb.NewStruct(provider.Credentials)
		if err != nil {
			s.logger.Error("handler", zap.Error(err))
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		item := &sirenv1.Provider{
			Id:          provider.Id,
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

func (s *GRPCServer) CreateProvider(_ context.Context, req *sirenv1.CreateProviderRequest) (*sirenv1.Provider, error) {
	provider, err := s.container.ProviderService.CreateProvider(&domain.Provider{
		Host:        req.GetHost(),
		Name:        req.GetName(),
		Type:        req.GetType(),
		Credentials: req.GetCredentials().AsMap(),
		Labels:      req.GetLabels(),
	})
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	grpcCredentials, err := structpb.NewStruct(provider.Credentials)
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	return &sirenv1.Provider{
		Id:          provider.Id,
		Host:        provider.Host,
		Name:        provider.Name,
		Type:        provider.Type,
		Credentials: grpcCredentials,
		Labels:      provider.Labels,
		CreatedAt:   timestamppb.New(provider.CreatedAt),
		UpdatedAt:   timestamppb.New(provider.UpdatedAt),
	}, nil
}

func (s *GRPCServer) GetProvider(_ context.Context, req *sirenv1.GetProviderRequest) (*sirenv1.Provider, error) {
	provider, err := s.container.ProviderService.GetProvider(req.GetId())
	if provider == nil {
		return nil, status.Errorf(codes.NotFound, "provider not found")
	}
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	grpcCredentials, err := structpb.NewStruct(provider.Credentials)
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	return &sirenv1.Provider{
		Id:          provider.Id,
		Host:        provider.Host,
		Name:        provider.Name,
		Type:        provider.Type,
		Credentials: grpcCredentials,
		Labels:      provider.Labels,
		CreatedAt:   timestamppb.New(provider.CreatedAt),
		UpdatedAt:   timestamppb.New(provider.UpdatedAt),
	}, nil
}

func (s *GRPCServer) UpdateProvider(_ context.Context, req *sirenv1.UpdateProviderRequest) (*sirenv1.Provider, error) {
	provider, err := s.container.ProviderService.UpdateProvider(&domain.Provider{
		Id:          req.GetId(),
		Host:        req.GetHost(),
		Name:        req.GetName(),
		Type:        req.GetType(),
		Credentials: req.GetCredentials().AsMap(),
		Labels:      req.GetLabels(),
	})
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	grpcCredentials, err := structpb.NewStruct(provider.Credentials)
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	return &sirenv1.Provider{
		Id:          provider.Id,
		Host:        provider.Host,
		Name:        provider.Name,
		Type:        provider.Type,
		Credentials: grpcCredentials,
		Labels:      provider.Labels,
		CreatedAt:   timestamppb.New(provider.CreatedAt),
		UpdatedAt:   timestamppb.New(provider.UpdatedAt),
	}, nil
}

func (s *GRPCServer) DeleteProvider(_ context.Context, req *sirenv1.DeleteProviderRequest) (*emptypb.Empty, error) {
	err := s.container.ProviderService.DeleteProvider(uint64(req.GetId()))
	if err != nil {
		return nil, helper.GRPCLogError(s.logger, codes.Internal, err)
	}

	return &emptypb.Empty{}, nil
}
