package v1

import (
	"context"
	"errors"
	sirenv1 "github.com/odpf/siren/api/proto/odpf/siren/v1"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/mocks"
	"github.com/odpf/siren/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"strings"
	"testing"
	"time"
)

func TestGRPCServer_ListProvider(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	t.Run("should return list of all provider", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyResult := []*domain.Provider{
			{
				Id:          1,
				Host:        "foo",
				Type:        "bar",
				Name:        "foo",
				Credentials: credentials,
				Labels:      labels,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		mockedProviderService.
			On("ListProviders").
			Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.ListProviders(context.Background(), &emptypb.Empty{})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res.GetProviders()))
		assert.Equal(t, "foo", res.GetProviders()[0].GetHost())
		assert.Equal(t, "bar", res.GetProviders()[0].GetType())
		assert.Equal(t, "foo", res.GetProviders()[0].GetName())
	})

	t.Run("should return error code 13 if getting providers failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedProviderService.
			On("ListProviders").
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.ListProviders(context.Background(), &emptypb.Empty{})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})

	t.Run("should return error code 13 if NewStruct conversion failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		credentials["bar"] = string([]byte{0xff})
		dummyResult := []*domain.Provider{
			{
				Id:          1,
				Host:        "foo",
				Type:        "bar",
				Name:        "foo",
				Credentials: credentials,
				Labels:      labels,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		mockedProviderService.
			On("ListProviders").
			Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.ListProviders(context.Background(), &emptypb.Empty{})
		assert.Nil(t, res)
		assert.Equal(t, strings.Replace(err.Error(), "\u00a0", " ", -1),
			"rpc error: code = Internal desc = proto: invalid UTF-8 in string: \"\\xff\"")
	})
}

func TestGRPCServer_CreateProvider(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	credentialsData, _ := structpb.NewStruct(credentials)

	payload := &domain.Provider{
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
	}
	dummyReq := &sirenv1.CreateProviderRequest{
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentialsData,
		Labels:      labels,
	}

	t.Run("should create provider object", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedProviderService.
			On("CreateProvider", payload).
			Return(payload, nil).Once()
		res, err := dummyGRPCServer.CreateProvider(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, "foo", res.GetName())
		assert.Equal(t, "bar", res.GetType())
		assert.Equal(t, "foo", res.GetName())
		assert.Equal(t, "bar", res.GetCredentials().GetFields()["foo"].GetStringValue())
	})

	t.Run("should return error code 13 if creating providers failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedProviderService.
			On("CreateProvider", payload).
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.CreateProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})

	t.Run("should return error code 13 if NewStruct conversion failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		credentials["bar"] = string([]byte{0xff})
		newPayload := &domain.Provider{
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
		}

		mockedProviderService.
			On("CreateProvider", mock.Anything).
			Return(newPayload, nil).Once()
		res, err := dummyGRPCServer.CreateProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.Equal(t, strings.Replace(err.Error(), "\u00a0", " ", -1),
			"rpc error: code = Internal desc = proto: invalid UTF-8 in string: \"\\xff\"")
	})
}

func TestGRPCServer_GetProvider(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	providerId := uint64(1)
	dummyReq := &sirenv1.GetProviderRequest{
		Id: 1,
	}

	t.Run("should return a provider", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyResult := &domain.Provider{
			Id:          1,
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockedProviderService.
			On("GetProvider", providerId).
			Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.GetProvider(context.Background(), dummyReq)
		assert.Nil(t, err)

		assert.Equal(t, "foo", res.GetHost())
		assert.Equal(t, "bar", res.GetType())
		assert.Equal(t, "foo", res.GetName())
		assert.Equal(t, "bar", res.GetCredentials().GetFields()["foo"].GetStringValue())
	})

	t.Run("should return error code 5 if no provider found", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedProviderService.
			On("GetProvider", providerId).
			Return(nil, nil).Once()

		res, err := dummyGRPCServer.GetProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = NotFound desc = provider not found")
	})

	t.Run("should return error code 13 if getting provider failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyResult := &domain.Provider{
			Id:          1,
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockedProviderService.
			On("GetProvider", providerId).
			Return(dummyResult, errors.New("random error")).Once()

		res, err := dummyGRPCServer.GetProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})

	t.Run("should return error code 13 if NewStruct conversion failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		credentials["bar"] = string([]byte{0xff})
		dummyResult := &domain.Provider{
			Id:          1,
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockedProviderService.
			On("GetProvider", providerId).
			Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.GetProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.Equal(t, strings.Replace(err.Error(), "\u00a0", " ", -1),
			"rpc error: code = Internal desc = proto: invalid UTF-8 in string: \"\\xff\"")
	})
}

func TestGRPCServer_UpdateProvider(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	credentialsData, _ := structpb.NewStruct(credentials)

	payload := &domain.Provider{
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
	}
	dummyReq := &sirenv1.UpdateProviderRequest{
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentialsData,
		Labels:      labels,
	}

	t.Run("should update provider object", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedProviderService.
			On("UpdateProvider", payload).
			Return(payload, nil).Once()
		res, err := dummyGRPCServer.UpdateProvider(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, "foo", res.GetName())
		assert.Equal(t, "bar", res.GetType())
		assert.Equal(t, "foo", res.GetName())
		assert.Equal(t, "bar", res.GetCredentials().GetFields()["foo"].GetStringValue())
	})

	t.Run("should return error code 13 if updating providers failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedProviderService.
			On("UpdateProvider", payload).
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.UpdateProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})

	t.Run("should return error code 13 if NewStruct conversion failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		credentials["bar"] = string([]byte{0xff})
		newPayload := &domain.Provider{
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
		}

		mockedProviderService.
			On("UpdateProvider", mock.Anything).
			Return(newPayload, nil).Once()
		res, err := dummyGRPCServer.UpdateProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.Equal(t, strings.Replace(err.Error(), "\u00a0", " ", -1),
			"rpc error: code = Internal desc = proto: invalid UTF-8 in string: \"\\xff\"")
	})
}

func TestGRPCServer_DeleteProvider(t *testing.T) {
	providerId := uint64(10)
	dummyReq := &sirenv1.DeleteProviderRequest{
		Id: uint64(10),
	}

	t.Run("should delete provider object", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedProviderService.
			On("DeleteProvider", providerId).
			Return(nil).Once()
		res, err := dummyGRPCServer.DeleteProvider(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, "", res.String())
	})

	t.Run("should return error code 13 if deleting providers failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedProviderService.
			On("DeleteProvider", providerId).
			Return(errors.New("random error")).Once()
		res, err := dummyGRPCServer.DeleteProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})
}
