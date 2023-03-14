package v1beta1_test

import (
	"context"
	"testing"
	"time"

	"github.com/goto/salt/log"
	"github.com/goto/siren/core/namespace"
	"github.com/goto/siren/core/provider"
	"github.com/goto/siren/internal/api"
	"github.com/goto/siren/internal/api/mocks"
	"github.com/goto/siren/internal/api/v1beta1"
	"github.com/goto/siren/pkg/errors"
	sirenv1beta1 "github.com/goto/siren/proto/gotocompany/siren/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestGRPCServer_ListNamespaces(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	t.Run("should return list of all namespaces", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{NamespaceService: mockedNamespaceService})
		dummyResult := []namespace.Namespace{
			{
				ID: 1,
				Provider: provider.Provider{
					ID: 2,
				},
				Name:        "foo",
				Credentials: credentials,
				Labels:      labels,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		mockedNamespaceService.EXPECT().List(mock.AnythingOfType("*context.emptyCtx")).Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.ListNamespaces(context.Background(), &sirenv1beta1.ListNamespacesRequest{})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res.GetNamespaces()))
		assert.Equal(t, "foo", res.GetNamespaces()[0].GetName())
		assert.Equal(t, uint64(1), res.GetNamespaces()[0].GetId())
		assert.Equal(t, uint64(2), res.GetNamespaces()[0].GetProvider())
		assert.Equal(t, "bar", res.GetNamespaces()[0].GetCredentials().GetFields()["foo"].GetStringValue())
	})

	t.Run("should return Internal if getting namespaces failed", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{NamespaceService: mockedNamespaceService})
		mockedNamespaceService.EXPECT().List(mock.AnythingOfType("*context.emptyCtx")).
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.ListNamespaces(context.Background(), &sirenv1beta1.ListNamespacesRequest{})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
	})

	t.Run("should return Internal if NewStruct conversion failed", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{NamespaceService: mockedNamespaceService})
		credentials["bar"] = string([]byte{0xff})
		dummyResult := []namespace.Namespace{
			{
				ID: 1,
				Provider: provider.Provider{
					ID: 2,
				},
				Name:        "foo",
				Credentials: credentials,
				Labels:      labels,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}
		mockedNamespaceService.EXPECT().List(mock.AnythingOfType("*context.emptyCtx")).Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.ListNamespaces(context.Background(), &sirenv1beta1.ListNamespacesRequest{})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
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
		Provider: provider.Provider{
			ID: 2,
		},
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
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{NamespaceService: mockedNamespaceService})
		mockedNamespaceService.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), payload).Run(func(_a0 context.Context, _a1 *namespace.Namespace) {
			_a1.ID = generatedID
		}).Return(nil).Once()
		res, err := dummyGRPCServer.CreateNamespace(context.Background(), request)
		assert.Nil(t, err)
		assert.Equal(t, generatedID, res.GetId())
	})

	t.Run("should return error Internal if creating namespaces failed", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{NamespaceService: mockedNamespaceService})
		mockedNamespaceService.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), payload).Return(errors.New("random error")).Once()
		res, err := dummyGRPCServer.CreateNamespace(context.Background(), request)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
	})

	t.Run("should return error Invalid Argument if create service return err invalid", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{NamespaceService: mockedNamespaceService})

		mockedNamespaceService.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), payload).Return(errors.ErrInvalid).Once()
		res, err := dummyGRPCServer.CreateNamespace(context.Background(), request)
		assert.Nil(t, res)
		assert.EqualError(t, err,
			"rpc error: code = InvalidArgument desc = request is not valid")
	})

	t.Run("should return error AlreadyExists if create service return err conflict", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{NamespaceService: mockedNamespaceService})

		mockedNamespaceService.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), payload).Return(errors.ErrConflict).Once()
		res, err := dummyGRPCServer.CreateNamespace(context.Background(), request)
		assert.Nil(t, res)
		assert.EqualError(t, err,
			"rpc error: code = AlreadyExists desc = an entity with conflicting identifier exists")
	})
}

func TestGRPCServer_GetNamespace(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	t.Run("should get the namespace", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{NamespaceService: mockedNamespaceService})
		dummyResult := &namespace.Namespace{
			ID: 1,
			Provider: provider.Provider{
				ID: 2,
			},
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockedNamespaceService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), uint64(1)).Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.GetNamespace(context.Background(),
			&sirenv1beta1.GetNamespaceRequest{Id: uint64(1)})
		assert.Nil(t, err)
		assert.Equal(t, "foo", res.GetNamespace().GetName())
		assert.Equal(t, uint64(1), res.GetNamespace().GetId())
		assert.Equal(t, uint64(2), res.GetNamespace().GetProvider())
		assert.Equal(t, "bar", res.GetNamespace().GetCredentials().GetFields()["foo"].GetStringValue())
	})

	t.Run("should return error Invalid Argument if no namespace found", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{NamespaceService: mockedNamespaceService})
		mockedNamespaceService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), uint64(1)).Return(nil, errors.ErrNotFound.WithCausef("some error")).Once()
		res, err := dummyGRPCServer.GetNamespace(context.Background(),
			&sirenv1beta1.GetNamespaceRequest{Id: uint64(1)})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = NotFound desc = requested entity not found")
	})

	t.Run("should return error Internal if getting namespace fails", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{NamespaceService: mockedNamespaceService})
		mockedNamespaceService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), uint64(1)).
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.GetNamespace(context.Background(),
			&sirenv1beta1.GetNamespaceRequest{Id: uint64(1)})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
	})

	t.Run("should return error Internal if NewStruct conversion failed", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{NamespaceService: mockedNamespaceService})
		credentials["bar"] = string([]byte{0xff})
		dummyResult := &namespace.Namespace{
			ID: 1,
			Provider: provider.Provider{
				ID: 2,
			},
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		mockedNamespaceService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), uint64(1)).Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.GetNamespace(context.Background(),
			&sirenv1beta1.GetNamespaceRequest{Id: uint64(1)})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
	})
}

func TestGRPCServer_UpdateNamespace(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	credentialsData, _ := structpb.NewStruct(credentials)
	payload := &namespace.Namespace{
		ID: 1,
		Provider: provider.Provider{
			ID: 2,
		},
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
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{NamespaceService: mockedNamespaceService})
		mockedNamespaceService.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), payload).Run(func(_a0 context.Context, _a1 *namespace.Namespace) {
			_a1.ID = payload.ID
		}).Return(nil).Once()
		res, err := dummyGRPCServer.UpdateNamespace(context.Background(), request)
		assert.Nil(t, err)
		assert.Equal(t, payload.ID, res.GetId())
		mockedNamespaceService.AssertExpectations(t)
	})

	t.Run("should return error Invalid Argument if namespace service return err invalid", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{NamespaceService: mockedNamespaceService})
		mockedNamespaceService.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), payload).Return(errors.ErrInvalid).Once()

		res, err := dummyGRPCServer.UpdateNamespace(context.Background(), request)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = request is not valid")
	})

	t.Run("should return error AlreadyExists if namespace service return err conflict", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{NamespaceService: mockedNamespaceService})
		mockedNamespaceService.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), payload).Return(errors.ErrConflict).Once()

		res, err := dummyGRPCServer.UpdateNamespace(context.Background(), request)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = AlreadyExists desc = an entity with conflicting identifier exists")
	})

	t.Run("should return error Internal if updating namespaces failed", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}

		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{NamespaceService: mockedNamespaceService})
		mockedNamespaceService.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), payload).Return(errors.New("random error")).Once()

		res, err := dummyGRPCServer.UpdateNamespace(context.Background(), request)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
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
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{NamespaceService: mockedNamespaceService})
		mockedNamespaceService.EXPECT().Delete(mock.AnythingOfType("*context.emptyCtx"), namespaceId).Return(nil).Once()
		res, err := dummyGRPCServer.DeleteNamespace(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, "", res.String())
		mockedNamespaceService.AssertExpectations(t)
	})

	t.Run("should return error Internal if deleting namespace failed", func(t *testing.T) {
		mockedNamespaceService := &mocks.NamespaceService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{NamespaceService: mockedNamespaceService})
		mockedNamespaceService.EXPECT().Delete(mock.AnythingOfType("*context.emptyCtx"), namespaceId).Return(errors.New("random error")).Once()
		res, err := dummyGRPCServer.DeleteNamespace(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
	})
}
