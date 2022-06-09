package alert_test

import (
	"errors"
	"testing"
	"time"

	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/core/alert/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Get(t *testing.T) {
	t.Run("should call repository Get method with proper arguments and return result in domain's type", func(t *testing.T) {
		repositoryMock := &mocks.AlertRepository{}
		dummyService := alert.NewService(repositoryMock)
		timenow := time.Now()
		dummyAlerts := []alert.Alert{
			{Id: 1, ProviderId: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "baz", MetricValue: "20",
				Rule: "bar", TriggeredAt: timenow},
			{Id: 2, ProviderId: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "baz", MetricValue: "0",
				Rule: "bar", TriggeredAt: timenow},
		}
		expectedAlerts := []alert.Alert{
			{Id: 1, ProviderId: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "baz", MetricValue: "20",
				Rule: "bar", TriggeredAt: timenow},
			{Id: 2, ProviderId: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "baz", MetricValue: "0",
				Rule: "bar", TriggeredAt: timenow},
		}
		repositoryMock.On("Get", "foo", uint64(1), uint64(0), uint64(100)).Return(dummyAlerts, nil)
		actualAlerts, err := dummyService.Get("foo", 1, 0, 100)
		assert.Nil(t, err)
		assert.Equal(t, expectedAlerts, actualAlerts)
		repositoryMock.AssertCalled(t, "Get", "foo", uint64(1), uint64(0), uint64(100))
	})

	t.Run("should call repository Get method with proper arguments if endtime is zero", func(t *testing.T) {
		repositoryMock := &mocks.AlertRepository{}
		dummyService := alert.NewService(repositoryMock)
		timenow := time.Now()
		dummyAlerts := []alert.Alert{
			{Id: 1, ProviderId: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "baz", MetricValue: "20",
				Rule: "bar", TriggeredAt: timenow},
			{Id: 2, ProviderId: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "baz", MetricValue: "0",
				Rule: "bar", TriggeredAt: timenow},
		}
		expectedAlerts := []alert.Alert{
			{Id: 1, ProviderId: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "baz", MetricValue: "20",
				Rule: "bar", TriggeredAt: timenow},
			{Id: 2, ProviderId: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "baz", MetricValue: "0",
				Rule: "bar", TriggeredAt: timenow},
		}
		repositoryMock.On("Get", "foo", uint64(1), uint64(0), mock.Anything).
			Return(dummyAlerts, nil)
		actualAlerts, err := dummyService.Get("foo", 1, 0, 0)
		assert.Nil(t, err)
		assert.Equal(t, expectedAlerts, actualAlerts)
		repositoryMock.AssertNotCalled(t, "Get", "foo", uint64(1), uint64(0), uint64(0))
	})

	t.Run("should call repository Get method and handle errors", func(t *testing.T) {
		repositoryMock := &mocks.AlertRepository{}
		dummyService := alert.NewService(repositoryMock)
		repositoryMock.On("Get", "foo", uint64(1), uint64(0), uint64(100)).
			Return(nil, errors.New("random error"))
		actualAlerts, err := dummyService.Get("foo", 1, 0, 100)
		assert.EqualError(t, err, "random error")
		assert.Nil(t, actualAlerts)
	})
}

func TestService_Create(t *testing.T) {
	timenow := time.Now()

	t.Run("should call repository Create method with proper arguments for firing alerts", func(t *testing.T) {
		repositoryMock := &mocks.AlertRepository{}
		dummyService := alert.NewService(repositoryMock)
		alertsToBeCreated := &alert.Alerts{Alerts: []alert.Alert{
			{Id: 1, ProviderId: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "lag", MetricValue: "20",
				Rule: "lagHigh", TriggeredAt: timenow},
		}}
		expectedAlerts := []alert.Alert{
			{Id: 1, ProviderId: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "lag", MetricValue: "20",
				Rule: "lagHigh", TriggeredAt: timenow},
		}
		repositoryMock.On("Create", mock.Anything).Return(nil)
		actualAlerts, err := dummyService.Create(alertsToBeCreated)
		assert.Nil(t, err)
		assert.Equal(t, expectedAlerts, actualAlerts)
		repositoryMock.AssertNumberOfCalls(t, "Create", 1)
	})

	t.Run("should call repository Create method with proper arguments for resolved alerts", func(t *testing.T) {
		repositoryMock := &mocks.AlertRepository{}
		dummyService := alert.NewService(repositoryMock)
		alertsToBeCreated := &alert.Alerts{Alerts: []alert.Alert{
			{Id: 1, ProviderId: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "lag", MetricValue: "20",
				Rule: "lagHigh", TriggeredAt: timenow},
		}}
		expectedAlerts := []alert.Alert{
			{Id: 1, ProviderId: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "lag", MetricValue: "20",
				Rule: "lagHigh", TriggeredAt: timenow},
		}
		repositoryMock.On("Create", mock.Anything).Return(nil)
		actualAlerts, err := dummyService.Create(alertsToBeCreated)
		assert.Nil(t, err)
		assert.Equal(t, expectedAlerts, actualAlerts)
		repositoryMock.AssertNumberOfCalls(t, "Create", 1)
	})

	t.Run("should handle errors from repository", func(t *testing.T) {
		repositoryMock := &mocks.AlertRepository{}
		dummyService := alert.NewService(repositoryMock)
		alertsToBeCreated := &alert.Alerts{Alerts: []alert.Alert{
			{Id: 1, ProviderId: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "lag", MetricValue: "20",
				Rule: "lagHigh", TriggeredAt: timenow},
		}}
		repositoryMock.On("Create", mock.Anything).Return(errors.New("random error"))
		actualAlerts, err := dummyService.Create(alertsToBeCreated)
		assert.EqualError(t, err, "random error")
		assert.Nil(t, actualAlerts)
	})
}

func TestService_Migrate(t *testing.T) {
	t.Run("should call repository Migrate method and return result", func(t *testing.T) {
		repositoryMock := &mocks.AlertRepository{}
		dummyService := alert.NewService(repositoryMock)
		repositoryMock.On("Migrate").Return(nil).Once()
		err := dummyService.Migrate()
		assert.Nil(t, err)
		repositoryMock.AssertCalled(t, "Migrate")
	})
}
