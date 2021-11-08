package alerts

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/odpf/siren/mocks"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
	"time"
)

// AnyTime is used to expect arbitrary time value
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type RepositoryTestSuite struct {
	suite.Suite
	sqldb      *sql.DB
	dbmock     sqlmock.Sqlmock
	repository AlertRepository
}

func (s *RepositoryTestSuite) SetupTest() {
	db, mock, _ := mocks.NewStore()
	s.sqldb, _ = db.DB()
	s.dbmock = mock
	s.repository = NewRepository(db)
}

func (s *RepositoryTestSuite) TearDownTest() {
	s.sqldb.Close()
}

func (s *RepositoryTestSuite) TestGet() {
	timenow := time.Now()
	s.Run("should fetch matching alert history objects", func() {
		expectedQuery := regexp.QuoteMeta(`select * from alerts where resource_name = 'foo' AND provider_id = '1' AND triggered_at BETWEEN to_timestamp('0') AND to_timestamp('1000')`)
		expectedAlert := Alert{
			Id: 1, ProviderId: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "baz", MetricValue: "20",
			Rule: "bar", TriggeredAt: timenow, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		expectedAlerts := []Alert{expectedAlert}
		expectedRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "resource_name", "provider_id",
			"severity", "metric_name", "metric_value", "rule", "triggered_at"}).
			AddRow(expectedAlert.Id, expectedAlert.CreatedAt,
				expectedAlert.UpdatedAt, expectedAlert.ResourceName, expectedAlert.ProviderId, expectedAlert.Severity,
				expectedAlert.MetricName, expectedAlert.MetricValue, expectedAlert.Rule, expectedAlert.TriggeredAt)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		actualAlerts, err := s.repository.Get("foo", 1, 0, 1000)
		s.Equal(expectedAlerts, actualAlerts)
		s.Nil(err)
	})

	s.Run("should return error if any in fetching alert history objects", func() {
		expectedQuery := regexp.QuoteMeta(`select * from alerts where resource_name = 'foo' AND provider_id = '1' AND triggered_at BETWEEN to_timestamp('0') AND to_timestamp('1000')`)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnError(errors.New("random error"))
		actualAlerts, err := s.repository.Get("foo", 1, 0, 1000)
		s.Nil(actualAlerts)
		s.EqualError(err, "random error")
	})
}

func (s *RepositoryTestSuite) TestCreate() {
	timenow := time.Now()
	s.Run("should create alert object", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "alerts" ("provider_id","resource_name","metric_name","metric_value","severity","rule","triggered_at","created_at","updated_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "id"`)
		expectedAlerts := &Alert{
			Id: 1, ProviderId: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "baz", MetricValue: "20",
			Rule: "bar", TriggeredAt: timenow, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedAlerts.ProviderId,
			expectedAlerts.ResourceName, expectedAlerts.MetricName, expectedAlerts.MetricValue, expectedAlerts.Severity,
			expectedAlerts.Rule, expectedAlerts.TriggeredAt, expectedAlerts.CreatedAt, expectedAlerts.UpdatedAt,
			expectedAlerts.Id).
			WillReturnRows(sqlmock.NewRows(nil))
		actualAlert, err := s.repository.Create(expectedAlerts)
		s.Equal(expectedAlerts, actualAlert)
		s.Nil(err)
	})

	s.Run("should return error in alert history creation", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "alerts" ("provider_id","resource_name","metric_name","metric_value","severity","rule","triggered_at","created_at","updated_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "id"`)
		expectedAlerts := &Alert{
			Id: 1, ProviderId: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "baz", MetricValue: "20",
			Rule: "bar", TriggeredAt: timenow, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedAlerts.ProviderId,
			expectedAlerts.ResourceName, expectedAlerts.MetricName, expectedAlerts.MetricValue, expectedAlerts.Severity,
			expectedAlerts.Rule, expectedAlerts.TriggeredAt, expectedAlerts.CreatedAt, expectedAlerts.UpdatedAt,
			expectedAlerts.Id).
			WillReturnError(errors.New("random error"))
		actualAlert, err := s.repository.Create(expectedAlerts)
		s.Nil(actualAlert)
		s.EqualError(err, "random error")
	})
}

func TestRepository(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
