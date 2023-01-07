package postgres_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/odpf/salt/db"
	"github.com/odpf/salt/dockertestx"
	saltlog "github.com/odpf/salt/log"
	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/core/log"
	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/internal/store/postgres"
	"github.com/odpf/siren/pkg/pgc"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
)

type LogRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	client     *pgc.Client
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.LogRepository

	namespaces    []namespace.EncryptedNamespace
	subscriptions []subscription.Subscription
	receivers     []receiver.Receiver
	silencesIDs   []string
	notifications []notification.Notification
	alerts        []alert.Alert
}

func (s *LogRepositoryTestSuite) SetupSuite() {
	var err error

	logger := saltlog.NewZap()
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
	s.repository = postgres.NewLogRepository(s.client)

	_, err = bootstrapProvider(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
	s.namespaces, err = bootstrapNamespace(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
	s.subscriptions, err = bootstrapSubscription(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
	s.receivers, err = bootstrapReceiver(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
	s.notifications, err = bootstrapNotification(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
	s.alerts, err = bootstrapAlert(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
	s.silencesIDs, err = bootstrapSilence(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *LogRepositoryTestSuite) SetupTest() {
	if err := bootstrapNotificationLog(
		s.client,
		s.namespaces,
		s.subscriptions,
		s.receivers,
		s.silencesIDs,
		s.notifications,
		s.alerts,
	); err != nil {
		s.T().Fatal(err)
	}
}

func (s *LogRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *LogRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *LogRepositoryTestSuite) cleanup() error {
	queries := []string{
		"TRUNCATE TABLE notification_log RESTART IDENTITY CASCADE",
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *LogRepositoryTestSuite) TestCreate() {
	type testCase struct {
		Description   string
		ItemsToCreate []log.Notification
		ErrString     string
	}

	var testCases = []testCase{
		{
			Description: "should create a notification_log",
			ItemsToCreate: []log.Notification{
				{
					NamespaceID:    1,
					NotificationID: s.notifications[0].ID,
					SubscriptionID: 1,
				},
			},
		},
		{
			Description: "should return error if a notification_log is invalid",
			ItemsToCreate: []log.Notification{{
				NamespaceID:    1111,
				NotificationID: "nid",
				SubscriptionID: 1111,
				AlertIDs:       []int64{11},
			}},
			ErrString: "pq: insert or update on table \"notification_log\" violates foreign key constraint \"notification_log_notification_id_fkey\"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			err := s.repository.BulkCreate(s.ctx, tc.ItemsToCreate)
			if err != nil {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
		})
	}
}

func (s *LogRepositoryTestSuite) TestListAlertIDsBySilenceID() {
	tests := []struct {
		name          string
		silenceID     string
		want          []int64
		wantErrString string
	}{
		{
			name:      "should return list of alert id if silence id exist",
			silenceID: s.silencesIDs[0],
			want: []int64{
				int64(s.alerts[0].ID),
				int64(s.alerts[2].ID),
			},
		},
		{
			name:      "should return nil if silence id does not exist",
			silenceID: "abc",
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			r := postgres.NewLogRepository(s.client)
			got, err := r.ListAlertIDsBySilenceID(context.TODO(), tt.silenceID)
			if err != nil {
				s.Assert().Equal(tt.wantErrString, err.Error())
				return
			}
			if diff := cmp.Diff(got, tt.want, cmpopts.SortSlices(func(x int64, y int64) bool {
				return x < y
			})); diff != "" {
				s.T().Errorf("LogRepository.ListAlertIDsBySilenceID() diff = %v", diff)
			}
		})
	}
}

func (s *LogRepositoryTestSuite) TestListSubscriptionIDsBySilenceID() {
	tests := []struct {
		name          string
		silenceID     string
		want          []int64
		wantErrString string
	}{
		{
			name:      "should return list of subscription id if silence id exist",
			silenceID: s.silencesIDs[0],
			want: []int64{
				int64(s.subscriptions[0].ID),
				int64(s.subscriptions[2].ID),
			},
		},
		{
			name:      "should return nil if silence id does not exist",
			silenceID: "abc",
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			r := postgres.NewLogRepository(s.client)
			got, err := r.ListSubscriptionIDsBySilenceID(context.TODO(), tt.silenceID)
			if err != nil {
				s.Assert().Equal(tt.wantErrString, err.Error())
				return
			}
			if diff := cmp.Diff(got, tt.want, cmpopts.SortSlices(func(x int64, y int64) bool {
				return x < y
			})); diff != "" {
				s.T().Errorf("LogRepository.ListSubscriptionIDsBySilenceID() diff = %v", diff)
			}
		})
	}
}

func TestLogRepository(t *testing.T) {
	suite.Run(t, new(LogRepositoryTestSuite))
}
