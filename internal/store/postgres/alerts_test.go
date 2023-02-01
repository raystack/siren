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
	"github.com/odpf/siren/pkg/pgc"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
)

type AlertsRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	client     *pgc.Client
	repository *postgres.AlertRepository
}

func (s *AlertsRepositoryTestSuite) SetupSuite() {
	var err error

	logger := log.NewZap()
	dpg, err := dockertestx.CreatePostgres(
		dockertestx.PostgresWithDetail(
			pgUser, pgPass, pgDBName,
		),
		dockertestx.PostgresWithVersionTag("13"),
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

	s.client, err = pgc.NewClient(logger, dbc)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	s.Require().NoError(migrate(s.ctx, logger, s.client, dbConfig))

	s.repository = postgres.NewAlertRepository(s.client)

	_, err = bootstrapProvider(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *AlertsRepositoryTestSuite) SetupTest() {
	_, err := bootstrapAlert(s.client)
	if err != nil {
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
		AlertToCreate alert.Alert
		ExpectedID    uint64
		ErrString     string
	}

	var testCases = []testCase{
		{
			Description: "should create an alert",
			AlertToCreate: alert.Alert{
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
			AlertToCreate: alert.Alert{
				ProviderID:   1000,
				ResourceName: "odpf-kafka-stream",
				MetricName:   "cpu_usage_user",
				MetricValue:  "88.88",
				Severity:     "CRITICAL",
				Rule:         "cpu-usage",
			},
			ErrString: "provider id does not exist",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			_, err := s.repository.Create(s.ctx, tc.AlertToCreate)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
		})
	}
}

func (s *AlertsRepositoryTestSuite) TestBulkUpdateSilence() {
	type testCase struct {
		Description    string
		SilenceStatus  string
		ExpectedAlerts []alert.Alert
		ErrString      string
	}

	var testCases = []testCase{
		{
			Description:   "should update 2 alerts to silence",
			SilenceStatus: alert.SilenceStatusTotal,
			ExpectedAlerts: []alert.Alert{
				{
					ID:            2,
					ProviderID:    1,
					ResourceName:  "odpf-kafka-2",
					MetricName:    "cpu_usage_user",
					MetricValue:   "97.95",
					Severity:      "WARNING",
					Rule:          "cpu-usage",
					SilenceStatus: alert.SilenceStatusTotal,
				},
				{
					ID:            3,
					ProviderID:    1,
					ResourceName:  "odpf-kafka-1",
					MetricName:    "cpu_usage_user",
					MetricValue:   "98.30",
					Severity:      "CRITICAL",
					Rule:          "cpu-usage",
					SilenceStatus: alert.SilenceStatusTotal,
				},
			},
		},
		{
			Description: "should return error foreign key if provider id does not exist",
			ErrString:   "err",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			err := s.repository.BulkUpdateSilence(s.ctx, []int64{2, 3}, tc.SilenceStatus)
			if err != nil {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if len(tc.ExpectedAlerts) != 0 {
				alerts, err := s.repository.List(s.ctx, alert.Filter{
					SilenceID: "",
				})
				s.Assert().NoError(err)

				var silencedAlerts []alert.Alert
				for _, al := range alerts {
					if al.SilenceStatus != "" {
						silencedAlerts = append(silencedAlerts, al)
					}
				}

				if diff := cmp.Diff(silencedAlerts, tc.ExpectedAlerts, cmpopts.IgnoreFields(alert.Alert{}, "TriggeredAt", "CreatedAt", "UpdatedAt")); diff != "" {
					s.T().Fatalf("got diff %v", diff)
				}
			}
		})
	}
}

func TestAlertsRepository(t *testing.T) {
	suite.Run(t, new(AlertsRepositoryTestSuite))
}
