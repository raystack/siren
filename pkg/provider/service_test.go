package provider

import (
	"errors"
	"github.com/odpf/siren/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestListProviders(t *testing.T) {
	credentials := make(StringInterfaceMap)
	credentials["foo"] = "bar"
	labels := make(StringStringMap)
	labels["foo"] = "bar"

	t.Run("should call repository List method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &MockProviderRepository{}
		dummyService := Service{repository: repositoryMock}
		dummyProviders := []*domain.Provider{
			{
				Id:          10,
				Host:        "foo",
				Type:        "bar",
				Name:        "foo",
				Credentials: credentials,
				Labels:      labels,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}
		providers := []*Provider{
			{
				Id:          10,
				Host:        "foo",
				Type:        "bar",
				Name:        "foo",
				Credentials: credentials,
				Labels:      labels,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}
		repositoryMock.On("List").Return(providers, nil).Once()
		result, err := dummyService.ListProviders()
		assert.Nil(t, err)
		assert.Equal(t, len(dummyProviders), len(result))
		assert.Equal(t, dummyProviders[0].Name, result[0].Name)
		repositoryMock.AssertCalled(t, "List")
	})

	t.Run("should call repository List method and return error if any", func(t *testing.T) {
		repositoryMock := &MockProviderRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("List").
			Return(nil, errors.New("random error")).Once()
		result, err := dummyService.ListProviders()
		assert.Nil(t, result)
		assert.EqualError(t, err, "service.repository.List: random error")
		repositoryMock.AssertCalled(t, "List")
	})
}

func TestCreateProviders(t *testing.T) {
	credentials := make(StringInterfaceMap)
	credentials["foo"] = "bar"
	labels := make(StringStringMap)
	labels["foo"] = "bar"
	timenow := time.Now()
	dummyProvider := &domain.Provider{
		Id:          10,
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
		CreatedAt:   timenow,
		UpdatedAt:   timenow,
	}
	provider := &Provider{
		Id:          10,
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
		CreatedAt:   timenow,
		UpdatedAt:   timenow,
	}

	t.Run("should call repository Create method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &MockProviderRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Create", provider).Return(provider, nil).Once()
		result, err := dummyService.CreateProvider(dummyProvider)
		assert.Nil(t, err)
		assert.Equal(t, dummyProvider, result)
		repositoryMock.AssertCalled(t, "Create", provider)
	})

	t.Run("should call repository Create method and return error if any", func(t *testing.T) {
		repositoryMock := &MockProviderRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Create", provider).
			Return(nil, errors.New("random error")).Once()
		result, err := dummyService.CreateProvider(dummyProvider)
		assert.Nil(t, result)
		assert.EqualError(t, err, "service.repository.Create: random error")
		repositoryMock.AssertCalled(t, "Create", provider)
	})
}

func TestGetProviders(t *testing.T) {
	providerID := uint64(10)
	credentials := make(StringInterfaceMap)
	credentials["foo"] = "bar"
	labels := make(StringStringMap)
	labels["foo"] = "bar"
	timenow := time.Now()
	dummyProvider := &domain.Provider{
		Id:          10,
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
		CreatedAt:   timenow,
		UpdatedAt:   timenow,
	}
	provider := &Provider{
		Id:          10,
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
		CreatedAt:   timenow,
		UpdatedAt:   timenow,
	}

	t.Run("should call repository Get method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &MockProviderRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Get", providerID).Return(provider, nil).Once()
		result, err := dummyService.GetProvider(providerID)
		assert.Nil(t, err)
		assert.Equal(t, dummyProvider, result)
		repositoryMock.AssertCalled(t, "Get", providerID)
	})

	t.Run("should call repository Get method and return error if any", func(t *testing.T) {
		repositoryMock := &MockProviderRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Get", providerID).
			Return(nil, errors.New("random error")).Once()
		result, err := dummyService.GetProvider(providerID)
		assert.Nil(t, result)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertCalled(t, "Get", providerID)
	})
}

func TestUpdateProviders(t *testing.T) {
	timenow := time.Now()
	credentials := make(StringInterfaceMap)
	credentials["foo"] = "bar"
	labels := make(StringStringMap)
	labels["foo"] = "bar"
	dummyProvider := &domain.Provider{
		Id:          10,
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
		CreatedAt:   timenow,
		UpdatedAt:   timenow,
	}
	provider := &Provider{
		Id:          10,
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
		CreatedAt:   timenow,
		UpdatedAt:   timenow,
	}

	t.Run("should call repository Update method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &MockProviderRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Update", provider).Return(provider, nil).Once()
		result, err := dummyService.UpdateProvider(dummyProvider)
		assert.Nil(t, err)
		assert.Equal(t, dummyProvider, result)
		repositoryMock.AssertCalled(t, "Update", provider)
	})

	t.Run("should call repository Update method and return error if any", func(t *testing.T) {
		repositoryMock := &MockProviderRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Update", provider).
			Return(nil, errors.New("random error")).Once()
		result, err := dummyService.UpdateProvider(dummyProvider)
		assert.Nil(t, result)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertCalled(t, "Update", provider)
	})
}

func TestDeleteProviders(t *testing.T) {
	credentials := make(StringInterfaceMap)
	credentials["foo"] = "bar"
	labels := make(StringStringMap)
	labels["foo"] = "bar"
	providerID := uint64(10)

	t.Run("should call repository Delete method and return nil if no error", func(t *testing.T) {
		repositoryMock := &MockProviderRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Delete", providerID).Return(nil).Once()
		err := dummyService.DeleteProvider(providerID)
		assert.Nil(t, err)
		repositoryMock.AssertCalled(t, "Delete", providerID)
	})

	t.Run("should call repository Delete method and return error if any", func(t *testing.T) {
		repositoryMock := &MockProviderRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Delete", providerID).
			Return(errors.New("random error")).Once()
		err := dummyService.DeleteProvider(providerID)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertCalled(t, "Delete", providerID)
	})
}

func TestService_Migrate(t *testing.T) {
	t.Run("should call repository Migrate method and return result", func(t *testing.T) {
		repositoryMock := &MockProviderRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Migrate").Return(nil).Once()
		err := dummyService.Migrate()
		assert.Nil(t, err)
		repositoryMock.AssertCalled(t, "Migrate")
	})
}
