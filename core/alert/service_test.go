package alert_test

import (
	"context"
	"testing"
	"time"

	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/core/alert/mocks"
	"github.com/odpf/siren/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Get(t *testing.T) {
	ctx := context.TODO()

	t.Run("should call repository List method with proper arguments and return result in domain's type", func(t *testing.T) {
		repositoryMock := &mocks.AlertRepository{}
		dummyService := alert.NewService(repositoryMock)
		timenow := time.Now()
		dummyAlerts := []alert.Alert{
			{ID: 1, ProviderID: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "baz", MetricValue: "20",
				Rule: "bar", TriggeredAt: timenow},
			{ID: 2, ProviderID: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "baz", MetricValue: "0",
				Rule: "bar", TriggeredAt: timenow},
		}
		repositoryMock.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), alert.Filter{
			ProviderID:   1,
			ResourceName: "foo",
			StartTime:    0,
			EndTime:      100,
		}).Return(dummyAlerts, nil)
		actualAlerts, err := dummyService.List(ctx, alert.Filter{
			ProviderID:   1,
			ResourceName: "foo",
			StartTime:    0,
			EndTime:      100,
		})
		assert.Nil(t, err)
		assert.NotEmpty(t, actualAlerts)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should call repository List method with proper arguments if endtime is zero", func(t *testing.T) {
		repositoryMock := &mocks.AlertRepository{}
		dummyService := alert.NewService(repositoryMock)
		timenow := time.Now()
		dummyAlerts := []alert.Alert{
			{ID: 1, ProviderID: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "baz", MetricValue: "20",
				Rule: "bar", TriggeredAt: timenow},
			{ID: 2, ProviderID: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "baz", MetricValue: "0",
				Rule: "bar", TriggeredAt: timenow},
		}
		repositoryMock.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(dummyAlerts, nil)
		actualAlerts, err := dummyService.List(ctx, alert.Filter{
			ProviderID:   1,
			ResourceName: "foo",
			StartTime:    0,
			EndTime:      0,
		})
		assert.Nil(t, err)
		assert.NotEmpty(t, actualAlerts)
		repositoryMock.AssertNotCalled(t, "Get", "foo", uint64(1), uint64(0), uint64(0))
	})

	t.Run("should call repository List method and handle errors", func(t *testing.T) {
		repositoryMock := &mocks.AlertRepository{}
		dummyService := alert.NewService(repositoryMock)
		repositoryMock.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).
			Return(nil, errors.New("random error"))
		actualAlerts, err := dummyService.List(ctx, alert.Filter{
			ProviderID:   1,
			ResourceName: "foo",
			StartTime:    0,
			EndTime:      0,
		})
		assert.EqualError(t, err, "random error")
		assert.Nil(t, actualAlerts)
	})
}

func TestService_Create(t *testing.T) {
	ctx := context.TODO()
	timenow := time.Now()

	t.Run("should call repository Create method with proper arguments ", func(t *testing.T) {
		repositoryMock := &mocks.AlertRepository{}
		dummyService := alert.NewService(repositoryMock)
		alertsToBeCreated := []*alert.Alert{
			{ID: 1, ProviderID: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "lag", MetricValue: "20",
				Rule: "lagHigh", TriggeredAt: timenow},
		}

		repositoryMock.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(alertsToBeCreated[0], nil)
		actualAlerts, err := dummyService.Create(ctx, alertsToBeCreated)
		assert.Nil(t, err)
		assert.NotEmpty(t, actualAlerts)
		repositoryMock.AssertNumberOfCalls(t, "Create", 1)
	})

	t.Run("should return error not found if repository return err relation", func(t *testing.T) {
		repositoryMock := &mocks.AlertRepository{}
		dummyService := alert.NewService(repositoryMock)
		alertsToBeCreated := []*alert.Alert{
			{ID: 1, ProviderID: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "lag", MetricValue: "20",
				Rule: "lagHigh", TriggeredAt: timenow},
		}
		repositoryMock.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(nil, alert.ErrRelation)
		actualAlerts, err := dummyService.Create(ctx, alertsToBeCreated)
		assert.EqualError(t, err, "provider id does not exist")
		assert.Nil(t, actualAlerts)
	})

	t.Run("should handle errors from repository", func(t *testing.T) {
		repositoryMock := &mocks.AlertRepository{}
		dummyService := alert.NewService(repositoryMock)
		alertsToBeCreated := []*alert.Alert{
			{ID: 1, ProviderID: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "lag", MetricValue: "20",
				Rule: "lagHigh", TriggeredAt: timenow},
		}
		repositoryMock.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(nil, errors.New("random error"))
		actualAlerts, err := dummyService.Create(ctx, alertsToBeCreated)
		assert.EqualError(t, err, "random error")
		assert.Nil(t, actualAlerts)
	})
}
