package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/odpf/salt/db"
	"github.com/odpf/salt/dockertestx"
	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/internal/store/postgres"
	"github.com/stretchr/testify/suite"
)

type AlertsRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	pool       *dockertestx.Pool
	resource   *dockertestx.Resource
	client     *postgres.Client
	repository *postgres.AlertRepository
}

func (s *AlertsRepositoryTestSuite) SetupSuite() {
	var err error

	logger := log.NewZap()
	dpg, err := dockertestx.CreatePostgres(
		dockertestx.PostgresWithDetail(
			pgUser, pgPass, pgDBName,
		),
	)
	if err != nil {
		s.T().Fatal(err)
	}

	s.pool = dpg.GetPool()
	s.resource = dpg.GetResource()

	dbConfig.URL = dpg.GetExternalConnString()
	dbc, err := db.New(dbConfig)
	if err != nil {
		s.T().Fatal(err)
	}

	s.client, err = postgres.NewClient(logger, dbc)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	migrate(s.ctx, logger, s.client, dbConfig)
	s.repository = postgres.NewAlertRepository(s.client)

	_, err = bootstrapProvider(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *AlertsRepositoryTestSuite) SetupTest() {
	var err error
	if err = bootstrapAlert(s.client); err != nil {
		s.T().Fatal(err)
	}
}

func (s *AlertsRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *AlertsRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *AlertsRepositoryTestSuite) cleanup() error {
	queries := []string{
		"TRUNCATE TABLE alerts RESTART IDENTITY CASCADE",
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *AlertsRepositoryTestSuite) TestList() {
	type testCase struct {
		Description    string
		Filter         alert.Filter
		ExpectedAlerts []alert.Alert
		ErrString      string
	}

	var testCases = []testCase{
		{
			Description: "should get all filtered alerts with correct filter",
			Filter: alert.Filter{
				ResourceName: "odpf-kafka-1",
				ProviderID:   1,
				StartTime:    int64(time.Date(2021, time.January, 0, 0, 0, 0, 0, time.UTC).Unix()),
				EndTime:      int64(time.Date(2022, time.January, 0, 0, 0, 0, 0, time.UTC).Unix()),
			},
			ExpectedAlerts: []alert.Alert{
				{
					ID:           1,
					ProviderID:   1,
					ResourceName: "odpf-kafka-1",
					MetricName:   "cpu_usage_user",
					MetricValue:  "97.30",
					Severity:     "CRITICAL",
					Rule:         "cpu-usage",
				},
				{
					ID:           3,
					ProviderID:   1,
					ResourceName: "odpf-kafka-1",
					MetricName:   "cpu_usage_user",
					MetricValue:  "98.30",
					Severity:     "CRITICAL",
					Rule:         "cpu-usage",
				},
			},
		},
		{
			Description: "should get empty alerts if out of range",
			Filter: alert.Filter{
				ResourceName: "odpf-kafka-1",
				ProviderID:   1,
				StartTime:    int64(time.Date(1980, time.January, 0, 0, 0, 0, 0, time.UTC).Unix()),
				EndTime:      int64(time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC).Unix()),
			},
			ExpectedAlerts: []alert.Alert{},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.List(s.ctx, tc.Filter)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedAlerts, cmpopts.IgnoreFields(alert.Alert{}, "TriggeredAt", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedAlerts)
			}
		})
	}
}

func (s *AlertsRepositoryTestSuite) TestCreate() {
	type testCase struct {
		Description   string
		AlertToCreate *alert.Alert
		ExpectedID    uint64
		ErrString     string
	}

	var testCases = []testCase{
		{
			Description: "should create an alert",
			AlertToCreate: &alert.Alert{
				ProviderID:   1,
				ResourceName: "odpf-kafka-stream",
				MetricName:   "cpu_usage_user",
				MetricValue:  "88.88",
				Severity:     "CRITICAL",
				Rule:         "cpu-usage",
			},
			ExpectedID: uint64(4), // autoincrement in db side
		},
		{
			Description: "should return error foreign key if provider id does not exist",
			AlertToCreate: &alert.Alert{
				ProviderID:   1000,
				ResourceName: "odpf-kafka-stream",
				MetricName:   "cpu_usage_user",
				MetricValue:  "88.88",
				Severity:     "CRITICAL",
				Rule:         "cpu-usage",
			},
			ErrString: "provider id does not exist",
		},
		{
			Description: "should return error if alert is nil",
			ErrString:   "alert domain is nil",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			err := s.repository.Create(s.ctx, tc.AlertToCreate)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
		})
	}
}

func TestAlertsRepository(t *testing.T) {
	suite.Run(t, new(AlertsRepositoryTestSuite))
}
