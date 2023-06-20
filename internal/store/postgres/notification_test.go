package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/goto/salt/db"
	"github.com/goto/salt/dockertestx"
	"github.com/goto/salt/log"
	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/internal/store/postgres"
	"github.com/goto/siren/pkg/pgc"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
)

type NotificationRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	client     *pgc.Client
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.NotificationRepository
}

func (s *NotificationRepositoryTestSuite) SetupSuite() {
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
	s.repository = postgres.NewNotificationRepository(s.client)
}

func (s *NotificationRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *NotificationRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *NotificationRepositoryTestSuite) cleanup() error {
	queries := []string{
		"TRUNCATE TABLE notifications RESTART IDENTITY CASCADE",
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *NotificationRepositoryTestSuite) TestCreate() {
	type testCase struct {
		Description          string
		NotificationToCreate notification.Notification
		ErrString            string
	}

	var testCases = []testCase{
		{
			Description: "should create a notification",
			NotificationToCreate: notification.Notification{
				NamespaceID: 1,
				Type:        notification.TypeReceiver,
				Data:        map[string]any{},
				Labels:      map[string]string{},
				CreatedAt:   time.Now(),
			},
		},
		{
			Description: "should return error if a notification is invalid",
			NotificationToCreate: notification.Notification{
				NamespaceID: 1,
				Type:        notification.TypeReceiver,
				Data: map[string]any{
					"k1": func(x chan struct{}) {
						<-x
					},
				},
				Labels:    map[string]string{},
				CreatedAt: time.Now(),
			},
			ErrString: "sql: converting argument $3 type: json: unsupported type: func(chan struct {})",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			_, err := s.repository.Create(s.ctx, tc.NotificationToCreate)
			if err != nil {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
		})
	}
}

func TestNotificationRepository(t *testing.T) {
	suite.Run(t, new(NotificationRepositoryTestSuite))
}
