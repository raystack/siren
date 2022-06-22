package provider_test

import (
	"errors"
	"testing"
	"time"

	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/provider/mocks"
	"github.com/odpf/siren/internal/store/model"
	"github.com/stretchr/testify/assert"
)

func TestListProviders(t *testing.T) {
	credentials := make(model.StringInterfaceMap)
	credentials["foo"] = "bar"
	labels := make(model.StringStringMap)
	labels["foo"] = "bar"

	t.Run("should call repository List method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &mocks.ProviderRepository{}
		dummyService := provider.NewService(repositoryMock)
		dummyProviders := []*provider.Provider{
			{
				ID:          10,
				Host:        "foo",
				Type:        "bar",
				Name:        "foo",
				Credentials: credentials,
				Labels:      labels,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}
		repositoryMock.EXPECT().List(map[string]interface{}{}).Return(dummyProviders, nil).Once()
		result, err := dummyService.ListProviders(map[string]interface{}{})
		assert.Nil(t, err)
		assert.Equal(t, len(dummyProviders), len(result))
		assert.Equal(t, dummyProviders[0].Name, result[0].Name)
		repositoryMock.AssertCalled(t, "List", map[string]interface{}{})
	})

	t.Run("should call repository List method and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.ProviderRepository{}
		dummyService := provider.NewService(repositoryMock)
		repositoryMock.EXPECT().List(map[string]interface{}{}).
			Return(nil, errors.New("random error")).Once()
		result, err := dummyService.ListProviders(map[string]interface{}{})
		assert.Nil(t, result)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertCalled(t, "List", map[string]interface{}{})
	})
}

func TestCreateProvider(t *testing.T) {
	credentials := make(model.StringInterfaceMap)
	credentials["foo"] = "bar"
	labels := make(model.StringStringMap)
	labels["foo"] = "bar"
	timenow := time.Now()
	dummyProvider := &provider.Provider{
		ID:          10,
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
		CreatedAt:   timenow,
		UpdatedAt:   timenow,
	}

	t.Run("should call repository Create method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &mocks.ProviderRepository{}
		dummyService := provider.NewService(repositoryMock)
		repositoryMock.EXPECT().Create(dummyProvider).Return(dummyProvider, nil).Once()
		result, err := dummyService.CreateProvider(dummyProvider)
		assert.Nil(t, err)
		assert.Equal(t, dummyProvider, result)
		repositoryMock.AssertCalled(t, "Create", dummyProvider)
	})

	t.Run("should call repository Create method and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.ProviderRepository{}
		dummyService := provider.NewService(repositoryMock)
		repositoryMock.EXPECT().Create(dummyProvider).
			Return(nil, errors.New("random error")).Once()
		result, err := dummyService.CreateProvider(dummyProvider)
		assert.Nil(t, result)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertCalled(t, "Create", dummyProvider)
	})
}

func TestGetProvider(t *testing.T) {
	providerID := uint64(10)
	credentials := make(model.StringInterfaceMap)
	credentials["foo"] = "bar"
	labels := make(model.StringStringMap)
	labels["foo"] = "bar"
	timenow := time.Now()
	dummyProvider := &provider.Provider{
		ID:          10,
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
		CreatedAt:   timenow,
		UpdatedAt:   timenow,
	}

	t.Run("should call repository Get method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &mocks.ProviderRepository{}
		dummyService := provider.NewService(repositoryMock)
		repositoryMock.EXPECT().Get(providerID).Return(dummyProvider, nil).Once()
		result, err := dummyService.GetProvider(providerID)
		assert.Nil(t, err)
		assert.Equal(t, dummyProvider, result)
		repositoryMock.AssertCalled(t, "Get", providerID)
	})

	t.Run("should call repository Get method and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.ProviderRepository{}
		dummyService := provider.NewService(repositoryMock)
		repositoryMock.EXPECT().Get(providerID).
			Return(nil, errors.New("random error")).Once()
		result, err := dummyService.GetProvider(providerID)
		assert.Nil(t, result)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertCalled(t, "Get", providerID)
	})
}

func TestUpdateProvider(t *testing.T) {
	timenow := time.Now()
	credentials := make(model.StringInterfaceMap)
	credentials["foo"] = "bar"
	labels := make(model.StringStringMap)
	labels["foo"] = "bar"
	dummyProvider := &provider.Provider{
		ID:          10,
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
		CreatedAt:   timenow,
		UpdatedAt:   timenow,
	}

	t.Run("should call repository Update method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &mocks.ProviderRepository{}
		dummyService := provider.NewService(repositoryMock)
		repositoryMock.EXPECT().Update(dummyProvider).Return(dummyProvider, nil).Once()
		result, err := dummyService.UpdateProvider(dummyProvider)
		assert.Nil(t, err)
		assert.Equal(t, dummyProvider, result)
		repositoryMock.AssertCalled(t, "Update", dummyProvider)
	})

	t.Run("should call repository Update method and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.ProviderRepository{}
		dummyService := provider.NewService(repositoryMock)
		repositoryMock.EXPECT().Update(dummyProvider).
			Return(nil, errors.New("random error")).Once()
		result, err := dummyService.UpdateProvider(dummyProvider)
		assert.Nil(t, result)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertCalled(t, "Update", dummyProvider)
	})
}

func TestDeleteProvider(t *testing.T) {
	credentials := make(model.StringInterfaceMap)
	credentials["foo"] = "bar"
	labels := make(model.StringStringMap)
	labels["foo"] = "bar"
	providerID := uint64(10)

	t.Run("should call repository Delete method and return nil if no error", func(t *testing.T) {
		repositoryMock := &mocks.ProviderRepository{}
		dummyService := provider.NewService(repositoryMock)
		repositoryMock.EXPECT().Delete(providerID).Return(nil).Once()
		err := dummyService.DeleteProvider(providerID)
		assert.Nil(t, err)
		repositoryMock.AssertCalled(t, "Delete", providerID)
	})

	t.Run("should call repository Delete method and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.ProviderRepository{}
		dummyService := provider.NewService(repositoryMock)
		repositoryMock.EXPECT().Delete(providerID).
			Return(errors.New("random error")).Once()
		err := dummyService.DeleteProvider(providerID)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertCalled(t, "Delete", providerID)
	})
}