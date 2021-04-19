package alert_history

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
	repository AlertHistoryRepository
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

	s.Run("should fetch matching alert history objects", func() {
		expectedQuery := regexp.QuoteMeta(`select * from alerts where resource = 'foo' AND created_at BETWEEN to_timestamp('0') AND to_timestamp('1000')`)
		expectedAlert := Alert{
			ID: 10, CreatedAt: time.Now(), UpdatedAt: time.Now(),
			Resource: "foo", Template: "bar", Level: "CRITICAL", MetricName: "baz", MetricValue: "20",
		}
		expectedAlerts := []Alert{expectedAlert}
		expectedRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "resource", "template",
			"metric_name", "metric_value", "level"}).
			AddRow(expectedAlert.ID, expectedAlert.CreatedAt,
				expectedAlert.UpdatedAt, expectedAlert.Resource,
				expectedAlert.Template, expectedAlert.MetricName,
				expectedAlert.MetricValue, expectedAlert.Level)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		actualAlerts, err := s.repository.Get("foo", 0, 1000)
		s.Equal(expectedAlerts, actualAlerts)
		s.Nil(err)
	})

	s.Run("should return error if any in fetching alert history objects", func() {
		expectedQuery := regexp.QuoteMeta(`select * from alerts where resource = 'foo' AND created_at BETWEEN to_timestamp('0') AND to_timestamp('1000')`)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnError(errors.New("random error"))
		actualAlerts, err := s.repository.Get("foo", 0, 1000)
		s.Nil(actualAlerts)
		s.EqualError(err, "random error")
	})
}

func (s *RepositoryTestSuite) TestCreate() {

	s.Run("should create alert history object", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "alerts" ("created_at","updated_at","resource","template","metric_name","metric_value","level","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)
		expectedAlerts := &Alert{
			ID: 10, CreatedAt: time.Now(), UpdatedAt: time.Now(),
			Resource: "foo", Template: "bar", Level: "CRITICAL", MetricName: "baz", MetricValue: "20",
		}
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedAlerts.CreatedAt,
			expectedAlerts.UpdatedAt, expectedAlerts.Resource, expectedAlerts.Template, expectedAlerts.MetricName,
			expectedAlerts.MetricValue, expectedAlerts.Level, expectedAlerts.ID).
			WillReturnRows(sqlmock.NewRows(nil))
		actualAlert, err := s.repository.Create(expectedAlerts)
		s.Equal(expectedAlerts, actualAlert)
		s.Nil(err)
	})

	s.Run("should return error in alert history creation", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "alerts" ("created_at","updated_at","resource","template","metric_name","metric_value","level","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)
		expectedAlerts := &Alert{
			ID: 10, CreatedAt: time.Now(), UpdatedAt: time.Now(),
			Resource: "foo", Template: "bar", Level: "CRITICAL", MetricName: "baz", MetricValue: "20",
		}
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedAlerts.CreatedAt,
			expectedAlerts.UpdatedAt, expectedAlerts.Resource, expectedAlerts.Template, expectedAlerts.MetricName,
			expectedAlerts.MetricValue, expectedAlerts.Level, expectedAlerts.ID).
			WillReturnError(errors.New("random error"))
		actualAlert, err := s.repository.Create(expectedAlerts)
		s.Nil(actualAlert)
		s.EqualError(err, "random error")
	})
}

func TestRepository(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
