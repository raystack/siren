package v1beta1

import (
	"context"
	"fmt"

	"github.com/odpf/siren/core/namespace"
	sirenv1beta1 "github.com/odpf/siren/internal/server/proto/odpf/siren/v1beta1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate mockery --name=NamespaceService -r --case underscore --with-expecter --structname NamespaceService --filename namespace_service.go --output=./mocks
type NamespaceService interface {
	List(context.Context) ([]namespace.Namespace, error)
	Create(context.Context, *namespace.Namespace) error
	Get(context.Context, uint64) (*namespace.Namespace, error)
	Update(context.Context, *namespace.Namespace) error
	Delete(context.Context, uint64) error
}

func (s *GRPCServer) ListNamespaces(ctx context.Context, _ *sirenv1beta1.ListNamespacesRequest) (*sirenv1beta1.ListNamespacesResponse, error) {
	namespaces, err := s.namespaceService.List(ctx)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	items := []*sirenv1beta1.Namespace{}
	for _, namespace := range namespaces {
		credentials, err := structpb.NewStruct(namespace.Credentials)
		if err != nil {
			return nil, s.generateRPCErr(fmt.Errorf("failed to fetch namespace credentials: %w", err))
		}

		item := &sirenv1beta1.Namespace{
			Id:          namespace.ID,
			Urn:         namespace.URN,
			Name:        namespace.Name,
			Credentials: credentials,
			Labels:      namespace.Labels,
			Provider:    namespace.Provider,
			CreatedAt:   timestamppb.New(namespace.CreatedAt),
			UpdatedAt:   timestamppb.New(namespace.UpdatedAt),
		}
		items = append(items, item)
	}
	return &sirenv1beta1.ListNamespacesResponse{
		Namespaces: items,
	}, nil
}

func (s *GRPCServer) CreateNamespace(ctx context.Context, req *sirenv1beta1.CreateNamespaceRequest) (*sirenv1beta1.CreateNamespaceResponse, error) {
	ns := &namespace.Namespace{
		Provider:    req.GetProvider(),
		URN:         req.GetUrn(),
		Name:        req.GetName(),
		Credentials: req.GetCredentials().AsMap(),
		Labels:      req.GetLabels(),
	}
	err := s.namespaceService.Create(ctx, ns)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.CreateNamespaceResponse{
		Id: ns.ID,
	}, nil
}

func (s *GRPCServer) GetNamespace(ctx context.Context, req *sirenv1beta1.GetNamespaceRequest) (*sirenv1beta1.GetNamespaceResponse, error) {
	namespace, err := s.namespaceService.Get(ctx, req.GetId())
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	credentials, err := structpb.NewStruct(namespace.Credentials)
	if err != nil {
		return nil, s.generateRPCErr(fmt.Errorf("failed to fetch namespace credentials: %w", err))
	}

	return &sirenv1beta1.GetNamespaceResponse{
		Namespace: &sirenv1beta1.Namespace{
			Id:          namespace.ID,
			Urn:         namespace.URN,
			Name:        namespace.Name,
			Credentials: credentials,
			Labels:      namespace.Labels,
			Provider:    namespace.Provider,
			CreatedAt:   timestamppb.New(namespace.CreatedAt),
			UpdatedAt:   timestamppb.New(namespace.UpdatedAt),
		},
	}, nil
}

func (s *GRPCServer) UpdateNamespace(ctx context.Context, req *sirenv1beta1.UpdateNamespaceRequest) (*sirenv1beta1.UpdateNamespaceResponse, error) {
	ns := &namespace.Namespace{
		ID:          req.GetId(),
		Provider:    req.GetProvider(),
		Name:        req.GetName(),
		Credentials: req.GetCredentials().AsMap(),
		Labels:      req.GetLabels(),
	}
	err := s.namespaceService.Update(ctx, ns)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.UpdateNamespaceResponse{
		Id: ns.ID,
	}, nil
}

func (s *GRPCServer) DeleteNamespace(ctx context.Context, req *sirenv1beta1.DeleteNamespaceRequest) (*sirenv1beta1.DeleteNamespaceResponse, error) {
	err := s.namespaceService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.DeleteNamespaceResponse{}, nil
}
