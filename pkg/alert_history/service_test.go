package alert_history

import (
	"errors"
	"github.com/odpf/siren/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestService_Get(t *testing.T) {
	t.Run("should call repository Get method with proper arguments and return result in domain's type", func(t *testing.T) {
		repositoryMock := &AlertHistoryRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		dummyAlerts := []Alert{
			{ID: 1, Resource: "foo", Template: "bar", MetricName: "baz", MetricValue: "20", Level: "CRITICAL"},
			{ID: 2, Resource: "foo", Template: "bar", MetricName: "baz", MetricValue: "0", Level: "RESOLVED"},
		}
		expectedAlerts := []domain.AlertHistoryObject{
			{ID: 1, Name: "foo", TemplateID: "bar", MetricName: "baz", MetricValue: "20", Level: "CRITICAL"},
			{ID: 2, Name: "foo", TemplateID: "bar", MetricName: "baz", MetricValue: "0", Level: "RESOLVED"},
		}
		repositoryMock.On("Get", "foo", uint32(0), uint32(100)).Return(dummyAlerts, nil)
		actualAlerts, err := dummyService.Get("foo", 0, 100)
		assert.Nil(t, err)
		assert.Equal(t, expectedAlerts, actualAlerts)
		repositoryMock.AssertCalled(t, "Get", "foo", uint32(0), uint32(100))
	})

	t.Run("should call repository Get method with proper arguments if endtime is zero", func(t *testing.T) {
		repositoryMock := &AlertHistoryRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		dummyAlerts := []Alert{
			{ID: 1, Resource: "foo", Template: "bar", MetricName: "baz", MetricValue: "20", Level: "CRITICAL"},
			{ID: 2, Resource: "foo", Template: "bar", MetricName: "baz", MetricValue: "0", Level: "RESOLVED"},
		}
		expectedAlerts := []domain.AlertHistoryObject{
			{ID: 1, Name: "foo", TemplateID: "bar", MetricName: "baz", MetricValue: "20", Level: "CRITICAL"},
			{ID: 2, Name: "foo", TemplateID: "bar", MetricName: "baz", MetricValue: "0", Level: "RESOLVED"},
		}
		repositoryMock.On("Get", "foo", uint32(0), mock.Anything).Return(dummyAlerts, nil)
		actualAlerts, err := dummyService.Get("foo", 0, 0)
		assert.Nil(t, err)
		assert.Equal(t, expectedAlerts, actualAlerts)
		repositoryMock.AssertNotCalled(t, "Get", "foo", uint32(0), uint32(0))
	})

	t.Run("should call repository Get method and handle errors", func(t *testing.T) {
		repositoryMock := &AlertHistoryRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Get", "foo", uint32(0), uint32(100)).Return(nil, errors.New("random error"))
		actualAlerts, err := dummyService.Get("foo", 0, 100)
		assert.EqualError(t, err, "random error")
		assert.Nil(t, actualAlerts)
	})
}

func TestService_Create(t *testing.T) {
	t.Run("should call repository Create method with proper arguments for firing alerts", func(t *testing.T) {
		repositoryMock := &AlertHistoryRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		dummyAlerts := []Alert{
			{ID: 1, Resource: "foo", Template: "lagHigh", MetricName: "lag", MetricValue: "20", Level: "CRITICAL"},
		}
		alertsToBeCreated := &domain.Alerts{Alerts: []domain.Alert{
			{Status: "firing", Labels: domain.Labels{Severity: "CRITICAL"},
				Annotations: domain.Annotations{Resource: "foo", Template: "lagHigh", MetricName: "lag", MetricValue: "20"}},
		}}
		expectedAlerts := []domain.AlertHistoryObject{
			{ID: 1, Name: "foo", TemplateID: "lagHigh", MetricName: "lag", MetricValue: "20", Level: "CRITICAL"},
		}
		repositoryMock.On("Create", mock.Anything).Return(&dummyAlerts[0], nil)
		actualAlerts, err := dummyService.Create(alertsToBeCreated)
		assert.Nil(t, err)
		assert.Equal(t, expectedAlerts, actualAlerts)
		repositoryMock.AssertNumberOfCalls(t, "Create", 1)
	})

	t.Run("should call repository Create method with proper arguments for resolved alerts", func(t *testing.T) {
		repositoryMock := &AlertHistoryRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		dummyAlerts := []Alert{
			{ID: 1, Resource: "foo", Template: "lagHigh", MetricName: "lag", MetricValue: "20", Level: "RESOLVED"},
		}
		alertsToBeCreated := &domain.Alerts{Alerts: []domain.Alert{
			{Status: "resolved", Labels: domain.Labels{Severity: "CRITICAL"},
				Annotations: domain.Annotations{Resource: "foo", Template: "lagHigh", MetricName: "lag", MetricValue: "20"}},
		}}
		expectedAlerts := []domain.AlertHistoryObject{
			{ID: 1, Name: "foo", TemplateID: "lagHigh", MetricName: "lag", MetricValue: "20", Level: "RESOLVED"},
		}
		repositoryMock.On("Create", mock.Anything).Return(&dummyAlerts[0], nil)
		actualAlerts, err := dummyService.Create(alertsToBeCreated)
		assert.Nil(t, err)
		assert.Equal(t, expectedAlerts, actualAlerts)
		repositoryMock.AssertNumberOfCalls(t, "Create", 1)
	})

	t.Run("should handle errors from repository", func(t *testing.T) {
		repositoryMock := &AlertHistoryRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		alertsToBeCreated := &domain.Alerts{Alerts: []domain.Alert{
			{Status: "resolved", Labels: domain.Labels{Severity: "CRITICAL"},
				Annotations: domain.Annotations{Resource: "foo", Template: "lagHigh", MetricName: "lag", MetricValue: "20"}},
		}}
		repositoryMock.On("Create", mock.Anything).Return(nil, errors.New("random error"))
		actualAlerts, err := dummyService.Create(alertsToBeCreated)
		assert.EqualError(t, err, "random error")
		assert.Nil(t, actualAlerts)
	})
}

func TestService_Migrate(t *testing.T) {
	t.Run("should call repository Migrate method and return result", func(t *testing.T) {
		repositoryMock := &AlertHistoryRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Migrate").Return(nil).Once()
		err := dummyService.Migrate()
		assert.Nil(t, err)
		repositoryMock.AssertCalled(t, "Migrate")
	})
}
