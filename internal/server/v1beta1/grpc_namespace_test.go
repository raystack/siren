package v1beta1

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/namespace"
	sirenv1beta1 "github.com/odpf/siren/internal/server/proto/odpf/siren/v1beta1"
	"github.com/odpf/siren/internal/server/v1beta1/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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
				ID:          1,
				Provider:    2,
				Name:        "foo",
				Credentials: credentials,
				Labels:      labels,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		mockedNamespaceService.EXPECT().ListNamespaces().Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.ListNamespaces(context.Background(), &sirenv1beta1.ListNamespacesRequest{})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res.GetData()))
		assert.Equal(t, "foo", res.GetData()[0].GetName())
		assert.Equal(t, uint64(1), res.GetData()[0].GetId())
		assert.Equal(t, uint64(2), res.GetData()[0].GetProvider())
		assert.Equal(t, "bar", res.GetData()[0].GetCredentials().GetFields()["foo"].GetStringValue())
	})

	t.Run("should return Internal if getting namespaces failed", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}
		mockedNamespaceService.EXPECT().ListNamespaces().
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.ListNamespaces(context.Background(), &sirenv1beta1.ListNamespacesRequest{})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})

	t.Run("should return Internal if NewStruct conversion failed", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}
		credentials["bar"] = string([]byte{0xff})
		dummyResult := []*namespace.Namespace{
			{
				ID:          1,
				Provider:    2,
				Name:        "foo",
				Credentials: credentials,
				Labels:      labels,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}
		mockedNamespaceService.EXPECT().ListNamespaces().Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.ListNamespaces(context.Background(), &sirenv1beta1.ListNamespacesRequest{})
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
	generatedID := uint64(77)

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
		mockedNamespaceService.EXPECT().CreateNamespace(payload).Run(func(ns *namespace.Namespace) {
			ns.ID = generatedID
		}).Return(nil).Once()
		res, err := dummyGRPCServer.CreateNamespace(context.Background(), request)
		assert.Nil(t, err)
		assert.Equal(t, generatedID, res.GetId())
	})

	t.Run("should return error Internal if creating namespaces failed", func(t *testing.T) {
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

	t.Run("should return error Invalid Argument if namespace urn conflict within a provider", func(t *testing.T) {
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
			ID:          1,
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
		assert.Equal(t, "foo", res.GetData().GetName())
		assert.Equal(t, uint64(1), res.GetData().GetId())
		assert.Equal(t, uint64(2), res.GetData().GetProvider())
		assert.Equal(t, "bar", res.GetData().GetCredentials().GetFields()["foo"].GetStringValue())
	})

	t.Run("should return error Invalid Argument if no namespace found", func(t *testing.T) {
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

	t.Run("should return error Internal if getting namespace fails", func(t *testing.T) {
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

	t.Run("should return error Internal if NewStruct conversion failed", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := GRPCServer{
			namespaceService: mockedNamespaceService,
			logger:           log.NewNoop(),
		}
		credentials["bar"] = string([]byte{0xff})
		dummyResult := &namespace.Namespace{
			ID:          1,
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
		ID:          1,
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
		assert.Equal(t, payload.ID, res.GetId())
		mockedNamespaceService.AssertExpectations(t)
	})

	t.Run("should return error Invalid Argument if namespace urn conflict within a provider", func(t *testing.T) {
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

	t.Run("should return error Internal if updating namespaces failed", func(t *testing.T) {
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

	t.Run("should return error Internal if deleting namespace failed", func(t *testing.T) {
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
