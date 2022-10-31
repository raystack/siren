package postgres_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/odpf/salt/db"
	"github.com/odpf/salt/dockertest"
	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/internal/store/postgres"
	"github.com/stretchr/testify/suite"
)

type SubscriptionRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	client     *postgres.Client
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.SubscriptionRepository
}

func (s *SubscriptionRepositoryTestSuite) SetupSuite() {
	var err error

	logger := log.NewZap()
	dpg, err := dockertest.CreatePostgres(
		dockertest.PostgresWithDetail(
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
	s.repository = postgres.NewSubscriptionRepository(s.client)

	_, err = bootstrapProvider(s.client)
	if err != nil {
		s.T().Fatal(err)
	}

	_, err = bootstrapNamespace(s.client)
	if err != nil {
		s.T().Fatal(err)
	}

	_, err = bootstrapReceiver(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *SubscriptionRepositoryTestSuite) SetupTest() {
	var err error
	_, err = bootstrapSubscription(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *SubscriptionRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *SubscriptionRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *SubscriptionRepositoryTestSuite) cleanup() error {
	queries := []string{
		"TRUNCATE TABLE subscriptions RESTART IDENTITY CASCADE",
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *SubscriptionRepositoryTestSuite) TestList() {
	type testCase struct {
		Description           string
		Filter                subscription.Filter
		ExpectedSubscriptions []subscription.Subscription
		ErrString             string
	}

	var testCases = []testCase{
		{
			Description: "should get all subscriptions",
			ExpectedSubscriptions: []subscription.Subscription{
				{
					ID:        1,
					URN:       "alert-history-odpf",
					Namespace: 2,
					Receivers: []subscription.Receiver{
						{
							ID: 1,
						},
					},
					Match: map[string]string{},
				},
				{
					ID:        2,
					URN:       "odpf-data-warning",
					Namespace: 1,
					Receivers: []subscription.Receiver{
						{
							ID: 1,
							Configuration: map[string]interface{}{
								"channel_name": "odpf-data",
							},
						},
					},
					Match: map[string]string{
						"environment": "integration",
						"team":        "odpf-data",
					},
				},
				{
					ID:        3,
					URN:       "odpf-pd",
					Namespace: 2,
					Receivers: []subscription.Receiver{
						{
							ID: 1,
							Configuration: map[string]interface{}{
								"channel_name": "odpf-data",
							},
						},
					},
					Match: map[string]string{
						"environment": "production",
						"severity":    "CRITICAL",
						"team":        "odpf-app",
					},
				},
			},
		},
		{
			Description: "should get filtered subscriptions by namespace id",
			Filter: subscription.Filter{
				NamespaceID: 1,
			},
			ExpectedSubscriptions: []subscription.Subscription{
				{
					ID:        2,
					URN:       "odpf-data-warning",
					Namespace: 1,
					Receivers: []subscription.Receiver{
						{
							ID: 1,
							Configuration: map[string]interface{}{
								"channel_name": "odpf-data",
							},
						},
					},
					Match: map[string]string{
						"environment": "integration",
						"team":        "odpf-data",
					},
				},
			},
		},
		{
			Description: "should get filtered subscriptions by match labels",
			Filter: subscription.Filter{
				Labels: map[string]string{
					"environment": "production",
					"severity":    "CRITICAL",
					"team":        "odpf-app",
				},
			},
			ExpectedSubscriptions: []subscription.Subscription{
				{
					ID:        3,
					URN:       "odpf-pd",
					Namespace: 2,
					Receivers: []subscription.Receiver{
						{
							ID: 1,
							Configuration: map[string]interface{}{
								"channel_name": "odpf-data",
							},
						},
					},
					Match: map[string]string{
						"environment": "production",
						"severity":    "CRITICAL",
						"team":        "odpf-app",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.List(s.ctx, tc.Filter)
			if err != nil && err.Error() != tc.ErrString {
				s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
			}
			if !cmp.Equal(got, tc.ExpectedSubscriptions, cmpopts.IgnoreFields(subscription.Subscription{}, "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedSubscriptions)
			}
		})
	}
}

func (s *SubscriptionRepositoryTestSuite) TestCreate() {
	type testCase struct {
		Description          string
		SubscriptionToUpsert *subscription.Subscription
		ExpectedID           uint64
		ErrString            string
	}

	var testCases = []testCase{
		{
			Description: "should create a subscription if doesn't exist",
			SubscriptionToUpsert: &subscription.Subscription{
				Namespace: 1,
				URN:       "foo",
				Match: map[string]string{
					"foo": "bar",
				},
				Receivers: []subscription.Receiver{
					{
						ID:            2,
						Configuration: map[string]interface{}{},
					},
					{
						ID: 1,
						Configuration: map[string]interface{}{
							"channel_name": "test",
						},
					},
				},
			},
			ExpectedID: uint64(4), // autoincrement in db side
		},
		{
			Description: "should return duplicate error if urn already exist",
			SubscriptionToUpsert: &subscription.Subscription{
				Namespace: 1,
				URN:       "foo",
				Match: map[string]string{
					"foo": "bar",
				},
				Receivers: []subscription.Receiver{
					{
						ID:            2,
						Configuration: map[string]interface{}{},
					},
					{
						ID: 1,
						Configuration: map[string]interface{}{
							"channel_name": "test",
						},
					},
				},
			},
			ErrString: "urn already exist",
		}, {
			Description: "should return relation error if namespace id does not exist",
			SubscriptionToUpsert: &subscription.Subscription{
				Namespace: 1000,
				URN:       "new-foo",
				Match: map[string]string{
					"foo": "bar",
				},
				Receivers: []subscription.Receiver{
					{
						ID:            2,
						Configuration: map[string]interface{}{},
					},
					{
						ID: 1,
						Configuration: map[string]interface{}{
							"channel_name": "test",
						},
					},
				},
			},
			ErrString: "namespace id does not exist",
		},
		{
			Description: "should return error if subscription is nil",
			ErrString:   "subscription domain is nil",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			err := s.repository.Create(s.ctx, tc.SubscriptionToUpsert)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
		})
	}
}

func (s *SubscriptionRepositoryTestSuite) TestGet() {
	type testCase struct {
		Description          string
		PassedID             uint64
		ExpectedSubscription *subscription.Subscription
		ErrString            string
	}

	var testCases = []testCase{
		{
			Description: "should get a subscription",
			PassedID:    uint64(3),
			ExpectedSubscription: &subscription.Subscription{
				ID:        3,
				URN:       "odpf-pd",
				Namespace: 2,
				Receivers: []subscription.Receiver{
					{
						ID: 1,
						Configuration: map[string]interface{}{
							"channel_name": "odpf-data",
						},
					},
				},
				Match: map[string]string{
					"environment": "production",
					"severity":    "CRITICAL",
					"team":        "odpf-app",
				},
			},
		},
		{
			Description: "should return not found if id not found",
			PassedID:    uint64(1000),
			ErrString:   "subscription with id 1000 not found",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.Get(s.ctx, tc.PassedID)
			if err != nil && err.Error() != tc.ErrString {
				s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
			}
			if !cmp.Equal(got, tc.ExpectedSubscription, cmpopts.IgnoreFields(subscription.Subscription{}, "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedSubscription)
			}
		})
	}
}

func (s *SubscriptionRepositoryTestSuite) TestUpdate() {
	type testCase struct {
		Description          string
		SubscriptionToUpsert *subscription.Subscription
		ExpectedID           uint64
		ErrString            string
	}

	var testCases = []testCase{
		{
			Description: "should update a subscription",
			SubscriptionToUpsert: &subscription.Subscription{
				ID:        3,
				URN:       "odpf-pd",
				Namespace: 2,
				Receivers: []subscription.Receiver{
					{
						ID: 3100,
					},
				},
				Match: map[string]string{
					"key": "label",
				},
			},
			ExpectedID: uint64(3),
		},
		{
			Description: "should return duplicate error if urn already exist",
			SubscriptionToUpsert: &subscription.Subscription{
				ID:        1,
				URN:       "odpf-pd",
				Namespace: 2,
				Receivers: []subscription.Receiver{
					{
						ID: 3100,
					},
				},
				Match: map[string]string{
					"key": "label",
				},
			},
			ErrString: "urn already exist",
		},
		{
			Description: "should return relation error if namespace id does not exist",
			SubscriptionToUpsert: &subscription.Subscription{
				ID:        3,
				URN:       "odpf-pd",
				Namespace: 1000,
				Receivers: []subscription.Receiver{
					{
						ID: 3100,
					},
				},
				Match: map[string]string{
					"key": "label",
				},
			},
			ErrString: "namespace id does not exist",
		},
		{
			Description: "should return not found error if id not found",
			SubscriptionToUpsert: &subscription.Subscription{
				ID:        3000,
				URN:       "odpf-pd",
				Namespace: 1,
				Receivers: []subscription.Receiver{
					{
						ID: 3100,
					},
				},
				Match: map[string]string{
					"key": "label",
				},
			},
			ErrString: "subscription with id 3000 not found",
		},
		{
			Description: "should return error if subscription is nil",
			ErrString:   "subscription domain is nil",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			err := s.repository.Update(s.ctx, tc.SubscriptionToUpsert)
			if err != nil && err.Error() != tc.ErrString {
				s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
			}
		})
	}
}

func (s *SubscriptionRepositoryTestSuite) TestDelete() {
	type testCase struct {
		Description string
		IDToDelete  uint64
		ErrString   string
	}

	var testCases = []testCase{
		{
			Description: "should delete a subscription",
			IDToDelete:  1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			err := s.repository.Delete(s.ctx, tc.IDToDelete)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
		})
	}
}

func TestSubscriptionRepository(t *testing.T) {
	suite.Run(t, new(SubscriptionRepositoryTestSuite))
}
