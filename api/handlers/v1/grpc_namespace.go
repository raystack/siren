package v1

import (
	"context"
	sirenv1 "github.com/odpf/siren/api/proto/odpf/siren/v1"
	"github.com/odpf/siren/domain"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

func (s *GRPCServer) ListNamespaces(_ context.Context, _ *emptypb.Empty) (*sirenv1.ListNamespacesResponse, error) {
	namespaces, err := s.container.NamespaceService.ListNamespaces()
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res := &sirenv1.ListNamespacesResponse{
		Namespaces: make([]*sirenv1.Namespace, 0),
	}
	for _, namespace := range namespaces {
		credentials, err := structpb.NewStruct(namespace.Credentials)
		if err != nil {
			s.logger.Error("handler", zap.Error(err))
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		item := &sirenv1.Namespace{
			Id:          namespace.Id,
			Name:        namespace.Name,
			Credentials: credentials,
			Labels:      namespace.Labels,
			Provider:    namespace.Provider,
			CreatedAt:   timestamppb.New(namespace.CreatedAt),
			UpdatedAt:   timestamppb.New(namespace.UpdatedAt),
		}
		res.Namespaces = append(res.Namespaces, item)
	}
	return res, nil
}

func (s *GRPCServer) CreateNamespace(_ context.Context, req *sirenv1.CreateNamespaceRequest) (*sirenv1.Namespace, error) {
	namespace, err := s.container.NamespaceService.CreateNamespace(&domain.Namespace{
		Provider:    req.GetProvider(),
		Name:        req.GetName(),
		Credentials: req.GetCredentials().AsMap(),
		Labels:      req.GetLabels(),
	})
	if err != nil {
		if strings.Contains(err.Error(), `violates unique constraint "name_provider_id_unique"`) {
			return nil, status.Errorf(codes.InvalidArgument, "name and provider pair already exist")
		}
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	grpcCredentials, err := structpb.NewStruct(namespace.Credentials)
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &sirenv1.Namespace{
		Id:          namespace.Id,
		Provider:    namespace.Provider,
		Name:        namespace.Name,
		Credentials: grpcCredentials,
		Labels:      namespace.Labels,
		CreatedAt:   timestamppb.New(namespace.CreatedAt),
		UpdatedAt:   timestamppb.New(namespace.UpdatedAt),
	}, nil
}

func (s *GRPCServer) GetNamespace(_ context.Context, req *sirenv1.GetNamespaceRequest) (*sirenv1.Namespace, error) {
	namespace, err := s.container.NamespaceService.GetNamespace(req.GetId())
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if namespace == nil {
		return nil, status.Errorf(codes.NotFound, "namespace not found")
	}

	credentials, err := structpb.NewStruct(namespace.Credentials)
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &sirenv1.Namespace{
		Id:          namespace.Id,
		Name:        namespace.Name,
		Credentials: credentials,
		Labels:      namespace.Labels,
		Provider:    namespace.Provider,
		CreatedAt:   timestamppb.New(namespace.CreatedAt),
		UpdatedAt:   timestamppb.New(namespace.UpdatedAt),
	}, nil
}

func (s *GRPCServer) UpdateNamespace(_ context.Context, req *sirenv1.UpdateNamespaceRequest) (*sirenv1.Namespace, error) {
	namespace, err := s.container.NamespaceService.UpdateNamespace(&domain.Namespace{
		Id:          req.GetId(),
		Provider:    req.GetProvider(),
		Name:        req.GetName(),
		Credentials: req.GetCredentials().AsMap(),
		Labels:      req.GetLabels(),
	})
	if err != nil {
		if strings.Contains(err.Error(), `violates unique constraint "name_provider_id_unique"`) {
			return nil, status.Errorf(codes.InvalidArgument, "name and provider pair already exist")
		}
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	grpcCredentials, err := structpb.NewStruct(namespace.Credentials)
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &sirenv1.Namespace{
		Id:          namespace.Id,
		Name:        namespace.Name,
		Provider:    namespace.Provider,
		Credentials: grpcCredentials,
		Labels:      namespace.Labels,
		CreatedAt:   timestamppb.New(namespace.CreatedAt),
		UpdatedAt:   timestamppb.New(namespace.UpdatedAt),
	}, nil
}

func (s *GRPCServer) DeleteNamespace(_ context.Context, req *sirenv1.DeleteNamespaceRequest) (*emptypb.Empty, error) {
	err := s.container.NamespaceService.DeleteNamespace(req.GetId())
	if err != nil {
		s.logger.Error("handler", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
