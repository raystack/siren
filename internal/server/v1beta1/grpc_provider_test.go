package v1beta1

import (
	"context"
	"testing"
	"time"

	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/provider"
	sirenv1beta1 "github.com/odpf/siren/internal/server/proto/odpf/siren/v1beta1"
	"github.com/odpf/siren/internal/server/v1beta1/mocks"
	"github.com/odpf/siren/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestGRPCServer_ListProvider(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	t.Run("should return list of all provider", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			providerService: mockedProviderService,
			logger:          log.NewNoop(),
		}
		dummyResult := []*provider.Provider{
			{
				ID:          1,
				Host:        "foo",
				Type:        "bar",
				Name:        "foo",
				Credentials: credentials,
				Labels:      labels,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		mockedProviderService.EXPECT().ListProviders(map[string]interface{}{
			"type": "",
			"urn":  "",
		}).
			Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.ListProviders(context.Background(), &sirenv1beta1.ListProvidersRequest{})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res.GetProviders()))
		assert.Equal(t, "foo", res.GetProviders()[0].GetHost())
		assert.Equal(t, "bar", res.GetProviders()[0].GetType())
		assert.Equal(t, "foo", res.GetProviders()[0].GetName())
	})

	t.Run("should return error Internal if getting providers failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			providerService: mockedProviderService,
			logger:          log.NewNoop(),
		}

		mockedProviderService.EXPECT().ListProviders(map[string]interface{}{
			"type": "",
			"urn":  "",
		}).
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.ListProviders(context.Background(), &sirenv1beta1.ListProvidersRequest{})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
	})

	t.Run("should return error Internal if NewStruct conversion failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			providerService: mockedProviderService,
			logger:          log.NewNoop(),
		}

		credentials["bar"] = string([]byte{0xff})
		dummyResult := []*provider.Provider{
			{
				ID:          1,
				Host:        "foo",
				Type:        "bar",
				Name:        "foo",
				Credentials: credentials,
				Labels:      labels,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		mockedProviderService.EXPECT().ListProviders(map[string]interface{}{
			"type": "",
			"urn":  "",
		}).
			Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.ListProviders(context.Background(), &sirenv1beta1.ListProvidersRequest{})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
	})
}

func TestGRPCServer_CreateProvider(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	testID := uint64(88)
	credentialsData, _ := structpb.NewStruct(credentials)

	payload := &provider.Provider{
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
	}
	returnedPayload := &provider.Provider{
		ID:          testID,
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
	}
	dummyReq := &sirenv1beta1.CreateProviderRequest{
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentialsData,
		Labels:      labels,
	}

	t.Run("should create provider object", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			providerService: mockedProviderService,
			logger:          log.NewNoop(),
		}

		mockedProviderService.EXPECT().CreateProvider(payload).Return(returnedPayload, nil).Once()
		res, err := dummyGRPCServer.CreateProvider(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, testID, res.GetId())
	})

	t.Run("should return error Invalid Argument if provider return error invalid", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			providerService: mockedProviderService,
			logger:          log.NewNoop(),
		}

		mockedProviderService.EXPECT().CreateProvider(payload).
			Return(nil, errors.ErrInvalid).Once()
		res, err := dummyGRPCServer.CreateProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = request is not valid")
	})

	t.Run("should return error AlreadyExist if provider return error conflict", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			providerService: mockedProviderService,
			logger:          log.NewNoop(),
		}

		mockedProviderService.EXPECT().CreateProvider(payload).
			Return(nil, errors.ErrConflict).Once()
		res, err := dummyGRPCServer.CreateProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = AlreadyExists desc = an entity with conflicting identifier exists")
	})

	t.Run("should return error Internal if creating providers failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			providerService: mockedProviderService,
			logger:          log.NewNoop(),
		}

		mockedProviderService.EXPECT().CreateProvider(payload).
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.CreateProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
	})
}

func TestGRPCServer_GetProvider(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	providerId := uint64(1)
	dummyReq := &sirenv1beta1.GetProviderRequest{
		Id: 1,
	}

	t.Run("should return a provider", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			providerService: mockedProviderService,
			logger:          log.NewNoop(),
		}
		dummyResult := &provider.Provider{
			ID:          1,
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockedProviderService.EXPECT().GetProvider(providerId).Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.GetProvider(context.Background(), dummyReq)
		assert.Nil(t, err)

		assert.Equal(t, "foo", res.GetProvider().GetHost())
		assert.Equal(t, "bar", res.GetProvider().GetType())
		assert.Equal(t, "foo", res.GetProvider().GetName())
		assert.Equal(t, "bar", res.GetProvider().GetCredentials().GetFields()["foo"].GetStringValue())
	})

	t.Run("should return error Not Found if no provider found", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			providerService: mockedProviderService,
			logger:          log.NewNoop(),
		}

		mockedProviderService.EXPECT().GetProvider(providerId).
			Return(nil, nil).Once()

		res, err := dummyGRPCServer.GetProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = NotFound desc = provider not found")
	})

	t.Run("should return error Internal if getting provider failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			providerService: mockedProviderService,
			logger:          log.NewNoop(),
		}
		dummyResult := &provider.Provider{
			ID:          1,
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockedProviderService.EXPECT().GetProvider(providerId).
			Return(dummyResult, errors.New("random error")).Once()

		res, err := dummyGRPCServer.GetProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
	})

	t.Run("should return error Internal if NewStruct conversion failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			providerService: mockedProviderService,
			logger:          log.NewNoop(),
		}

		credentials["bar"] = string([]byte{0xff})
		dummyResult := &provider.Provider{
			ID:          1,
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockedProviderService.EXPECT().GetProvider(providerId).
			Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.GetProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
	})
}

func TestGRPCServer_UpdateProvider(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	testID := uint64(88)
	credentialsData, _ := structpb.NewStruct(credentials)

	payload := &provider.Provider{
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
	}
	returnedPayload := &provider.Provider{
		ID:          testID,
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
	}
	dummyReq := &sirenv1beta1.UpdateProviderRequest{
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentialsData,
		Labels:      labels,
	}

	t.Run("should update provider object", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			providerService: mockedProviderService,
			logger:          log.NewNoop(),
		}

		mockedProviderService.EXPECT().UpdateProvider(payload).Return(returnedPayload, nil).Once()
		res, err := dummyGRPCServer.UpdateProvider(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, testID, res.GetId())
	})

	t.Run("should return error Invalid Argument if updating providers return err invalid", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			providerService: mockedProviderService,
			logger:          log.NewNoop(),
		}

		mockedProviderService.EXPECT().UpdateProvider(payload).
			Return(nil, errors.ErrInvalid).Once()
		res, err := dummyGRPCServer.UpdateProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = request is not valid")
	})

	t.Run("should return error AlreadyExist if updating providers return err conflict", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			providerService: mockedProviderService,
			logger:          log.NewNoop(),
		}

		mockedProviderService.EXPECT().UpdateProvider(payload).
			Return(nil, errors.ErrConflict).Once()
		res, err := dummyGRPCServer.UpdateProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = AlreadyExists desc = an entity with conflicting identifier exists")
	})

	t.Run("should return error Internal if updating providers failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			providerService: mockedProviderService,
			logger:          log.NewNoop(),
		}

		mockedProviderService.EXPECT().UpdateProvider(payload).
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.UpdateProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
	})
}

func TestGRPCServer_DeleteProvider(t *testing.T) {
	providerId := uint64(10)
	dummyReq := &sirenv1beta1.DeleteProviderRequest{
		Id: uint64(10),
	}

	t.Run("should delete provider object", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			providerService: mockedProviderService,
			logger:          log.NewNoop(),
		}

		mockedProviderService.EXPECT().DeleteProvider(providerId).Return(nil).Once()
		res, err := dummyGRPCServer.DeleteProvider(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, "", res.String())
	})

	t.Run("should return error Internal if deleting providers failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			providerService: mockedProviderService,
			logger:          log.NewNoop(),
		}

		mockedProviderService.EXPECT().DeleteProvider(providerId).Return(errors.New("random error")).Once()
		res, err := dummyGRPCServer.DeleteProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
	})
}