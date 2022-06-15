package v1beta1

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/internal/server/v1beta1/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	sirenv1beta1 "go.buf.build/odpf/gw/odpf/proton/odpf/siren/v1beta1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestGRPCServer_ListNamespaces(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	t.Run("should return list of all namespaces", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}
		dummyResult := []*namespace.Namespace{
			{
				Id:          1,
				Provider:    2,
				Name:        "foo",
				Credentials: credentials,
				Labels:      labels,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		mockedNamespaceService.EXPECT().ListNamespaces().Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.ListNamespaces(context.Background(), &emptypb.Empty{})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res.GetNamespaces()))
		assert.Equal(t, "foo", res.GetNamespaces()[0].GetName())
		assert.Equal(t, uint64(1), res.GetNamespaces()[0].GetId())
		assert.Equal(t, uint64(2), res.GetNamespaces()[0].GetProvider())
		assert.Equal(t, "bar", res.GetNamespaces()[0].GetCredentials().GetFields()["foo"].GetStringValue())
	})

	t.Run("should return code 13 if getting namespaces failed", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}
		mockedNamespaceService.EXPECT().ListNamespaces().
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.ListNamespaces(context.Background(), &emptypb.Empty{})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})

	t.Run("should return code 13 if NewStruct conversion failed", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}
		credentials["bar"] = string([]byte{0xff})
		dummyResult := []*namespace.Namespace{
			{
				Id:          1,
				Provider:    2,
				Name:        "foo",
				Credentials: credentials,
				Labels:      labels,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}
		mockedNamespaceService.EXPECT().ListNamespaces().Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.ListNamespaces(context.Background(), &emptypb.Empty{})
		assert.Nil(t, res)
		assert.Equal(t, strings.Replace(err.Error(), "\u00a0", " ", -1),
			"rpc error: code = Internal desc = proto: invalid UTF-8 in string: \"\\xff\"")
	})
}

func TestGRPCServer_CreateNamespaces(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	credentialsData, _ := structpb.NewStruct(credentials)
	payload := &namespace.Namespace{
		Provider:    2,
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
	}
	request := &sirenv1beta1.CreateNamespaceRequest{
		Provider:    2,
		Name:        "foo",
		Credentials: credentialsData,
		Labels:      labels,
	}

	t.Run("should create a namespace", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}
		mockedNamespaceService.EXPECT().CreateNamespace(payload).Return(nil).Once()
		res, err := dummyGRPCServer.CreateNamespace(context.Background(), request)
		assert.Nil(t, err)
		assert.Equal(t, "foo", res.GetName())
	})

	t.Run("should return error code 13 if creating namespaces failed", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}
		mockedNamespaceService.EXPECT().CreateNamespace(payload).
			Return(errors.New("random error")).Once()
		res, err := dummyGRPCServer.CreateNamespace(context.Background(), request)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})

	t.Run("should return error code 5 if namespace urn conflict within a provider", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}

		mockedNamespaceService.EXPECT().CreateNamespace(payload).Return(
			errors.New(`violates unique constraint "urn_provider_id_unique"`)).Once()
		res, err := dummyGRPCServer.CreateNamespace(context.Background(), request)
		assert.Nil(t, res)
		assert.EqualError(t, err,
			"rpc error: code = InvalidArgument desc = urn and provider pair already exist")
	})

	t.Run("should return error code 13 if NewStruct conversion failed", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}
		mockedNamespaceService.EXPECT().CreateNamespace(mock.AnythingOfType("*namespace.Namespace")).Return(nil).
			Run(func(n *namespace.Namespace) {
				credentials["bar"] = string([]byte{0xff})
				n.Credentials = credentials
			}).Once()
		res, err := dummyGRPCServer.CreateNamespace(context.Background(), request)
		assert.Nil(t, res)
		assert.Equal(t, strings.Replace(err.Error(), "\u00a0", " ", -1),
			"rpc error: code = Internal desc = proto: invalid UTF-8 in string: \"\\xff\"")
	})
}

func TestGRPCServer_GetNamespace(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	t.Run("should get the namespace", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}
		dummyResult := &namespace.Namespace{
			Id:          1,
			Provider:    2,
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockedNamespaceService.EXPECT().GetNamespace(uint64(1)).Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.GetNamespace(context.Background(),
			&sirenv1beta1.GetNamespaceRequest{Id: uint64(1)})
		assert.Nil(t, err)
		assert.Equal(t, "foo", res.GetName())
		assert.Equal(t, uint64(1), res.GetId())
		assert.Equal(t, uint64(2), res.GetProvider())
		assert.Equal(t, "bar", res.GetCredentials().GetFields()["foo"].GetStringValue())
	})

	t.Run("should return error code 5 if no namespace found", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}
		mockedNamespaceService.EXPECT().GetNamespace(uint64(1)).Return(nil, nil).Once()
		res, err := dummyGRPCServer.GetNamespace(context.Background(),
			&sirenv1beta1.GetNamespaceRequest{Id: uint64(1)})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = NotFound desc = namespace not found")
	})

	t.Run("should return error code 13 if getting namespace fails", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}
		mockedNamespaceService.EXPECT().GetNamespace(uint64(1)).
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.GetNamespace(context.Background(),
			&sirenv1beta1.GetNamespaceRequest{Id: uint64(1)})
		assert.Nil(t, res)
		assert.EqualError(t, err, `rpc error: code = Internal desc = random error`)
	})

	t.Run("should return error code 13 if NewStruct conversion failed", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}
		credentials["bar"] = string([]byte{0xff})
		dummyResult := &namespace.Namespace{
			Id:          1,
			Provider:    2,
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		mockedNamespaceService.EXPECT().GetNamespace(uint64(1)).Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.GetNamespace(context.Background(),
			&sirenv1beta1.GetNamespaceRequest{Id: uint64(1)})
		assert.Nil(t, res)
		assert.Equal(t, strings.Replace(err.Error(), "\u00a0", " ", -1),
			"rpc error: code = Internal desc = proto: invalid UTF-8 in string: \"\\xff\"")
	})
}

func TestGRPCServer_UpdateNamespace(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	credentialsData, _ := structpb.NewStruct(credentials)
	payload := &namespace.Namespace{
		Id:          1,
		Provider:    2,
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
	}
	request := &sirenv1beta1.UpdateNamespaceRequest{
		Id:          1,
		Provider:    2,
		Name:        "foo",
		Credentials: credentialsData,
		Labels:      labels,
	}

	t.Run("should update a namespace", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}
		mockedNamespaceService.EXPECT().UpdateNamespace(payload).Return(nil).Once()
		res, err := dummyGRPCServer.UpdateNamespace(context.Background(), request)
		assert.Nil(t, err)
		assert.Equal(t, "foo", res.GetName())
		mockedNamespaceService.AssertExpectations(t)
	})

	t.Run("should return error code 5 if namespace urn conflict within a provider", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}
		mockedNamespaceService.EXPECT().UpdateNamespace(payload).Return(
			errors.New(`violates unique constraint "urn_provider_id_unique"`)).Once()

		res, err := dummyGRPCServer.UpdateNamespace(context.Background(), request)
		assert.Nil(t, res)
		assert.EqualError(t, err,
			"rpc error: code = InvalidArgument desc = urn and provider pair already exist")
	})

	t.Run("should return error code 13 if updating namespaces failed", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}
		mockedNamespaceService.EXPECT().UpdateNamespace(payload).
			Return(errors.New("random error")).Once()
		res, err := dummyGRPCServer.UpdateNamespace(context.Background(), request)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
		mockedNamespaceService.AssertExpectations(t)
	})

	t.Run("should return error code 13 if NewStruct conversion failed", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}
		mockedNamespaceService.EXPECT().UpdateNamespace(mock.AnythingOfType("*namespace.Namespace")).Return(nil).
			Run(func(n *namespace.Namespace) {
				credentials["foo"] = string([]byte{0xff})
				n.Credentials = credentials
			}).Once()
		res, err := dummyGRPCServer.UpdateNamespace(context.Background(), request)
		assert.Nil(t, res)
		assert.Equal(t, strings.Replace(err.Error(), "\u00a0", " ", -1),
			"rpc error: code = Internal desc = proto: invalid UTF-8 in string: \"\\xff\"")
		mockedNamespaceService.AssertExpectations(t)
	})
}

func TestGRPCServer_DeleteNamespace(t *testing.T) {
	namespaceId := uint64(10)
	dummyReq := &sirenv1beta1.DeleteNamespaceRequest{
		Id: uint64(10),
	}

	t.Run("should delete namespace object", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}
		mockedNamespaceService.EXPECT().DeleteNamespace(namespaceId).Return(nil).Once()
		res, err := dummyGRPCServer.DeleteNamespace(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, "", res.String())
		mockedNamespaceService.AssertExpectations(t)
	})

	t.Run("should return error code 13 if deleting namespace failed", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}
		mockedNamespaceService.EXPECT().DeleteNamespace(namespaceId).Return(errors.New("random error")).Once()
		res, err := dummyGRPCServer.DeleteNamespace(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})
}
