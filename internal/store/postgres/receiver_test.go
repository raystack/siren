package postgres_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ory/dockertest/v3"
	"github.com/raystack/salt/db"
	"github.com/raystack/salt/dockertestx"
	"github.com/raystack/salt/log"
	"github.com/raystack/siren/core/receiver"
	"github.com/raystack/siren/internal/store/postgres"
	"github.com/raystack/siren/pkg/pgc"
	"github.com/stretchr/testify/suite"
)

type ReceiverRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	client     *pgc.Client
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.ReceiverRepository
}

func (s *ReceiverRepositoryTestSuite) SetupSuite() {
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
	s.repository = postgres.NewReceiverRepository(s.client)
}

func (s *ReceiverRepositoryTestSuite) SetupTest() {
	var err error
	_, err = bootstrapReceiver(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *ReceiverRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ReceiverRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ReceiverRepositoryTestSuite) cleanup() error {
	queries := []string{
		"TRUNCATE TABLE receivers RESTART IDENTITY CASCADE",
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *ReceiverRepositoryTestSuite) TestList() {
	type testCase struct {
		Description       string
		Filter            receiver.Filter
		ExpectedReceivers []receiver.Receiver
		ErrString         string
	}

	var testCases = []testCase{
		{
			Description: "should get all receivers",
			ExpectedReceivers: []receiver.Receiver{
				{
					ID:   1,
					Name: "raystack-slack",
					Type: "slack",
					Labels: map[string]string{
						"entity": "raystack,org-a,org-b",
					},
					Configurations: map[string]interface{}{
						"token":     "xxxxxxxxxx",
						"workspace": "Odpf",
					},
				},
				{
					ID:   2,
					Name: "alert-history",
					Type: "http",
					Labels: map[string]string{
						"entity": "raystack,org-a,org-b,org-c",
					},
					Configurations: map[string]interface{}{
						"url": "http://siren.raystack.io/v1beta1/alerts/cortex/1",
					},
				},
				{
					ID:   3,
					Name: "raystack_pagerduty",
					Type: "pagerduty",
					Labels: map[string]string{
						"entity": "raystack",
						"team":   "siren-raystack",
					},
					Configurations: map[string]interface{}{
						"service_key": "1212121212121212121212121",
					},
				},
			},
		},
		{
			Description: "should get filtered receivers with list of ids",
			Filter: receiver.Filter{
				ReceiverIDs: []uint64{2, 3},
			},
			ExpectedReceivers: []receiver.Receiver{
				{
					ID:   2,
					Name: "alert-history",
					Type: "http",
					Labels: map[string]string{
						"entity": "raystack,org-a,org-b,org-c",
					},
					Configurations: map[string]interface{}{
						"url": "http://siren.raystack.io/v1beta1/alerts/cortex/1",
					},
				},
				{
					ID:   3,
					Name: "raystack_pagerduty",
					Type: "pagerduty",
					Labels: map[string]string{
						"entity": "raystack",
						"team":   "siren-raystack",
					},
					Configurations: map[string]interface{}{
						"service_key": "1212121212121212121212121",
					},
				},
			},
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
			if !cmp.Equal(got, tc.ExpectedReceivers, cmpopts.IgnoreFields(receiver.Receiver{}, "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedReceivers)
			}
		})
	}
}

func (s *ReceiverRepositoryTestSuite) TestGet() {
	type testCase struct {
		Description      string
		PassedID         uint64
		ExpectedReceiver *receiver.Receiver
		ErrString        string
	}

	var testCases = []testCase{
		{
			Description: "should get a receivers",
			PassedID:    3,
			ExpectedReceiver: &receiver.Receiver{
				ID:   3,
				Name: "raystack_pagerduty",
				Type: "pagerduty",
				Labels: map[string]string{
					"entity": "raystack",
					"team":   "siren-raystack",
				},
				Configurations: map[string]interface{}{
					"service_key": "1212121212121212121212121",
				},
			},
		},
		{
			Description: "should return not found if id not found",
			PassedID:    1000,
			ErrString:   "receiver with id 1000 not found",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.Get(s.ctx, tc.PassedID)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedReceiver, cmpopts.IgnoreFields(receiver.Receiver{}, "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedReceiver)
			}
		})
	}
}

func (s *ReceiverRepositoryTestSuite) TestCreate() {
	type testCase struct {
		Description      string
		ReceiverToCreate *receiver.Receiver
		ExpectedID       uint64
		ErrString        string
	}

	var testCases = []testCase{
		{
			Description: "should create a provider",
			ReceiverToCreate: &receiver.Receiver{
				Name: "neworg_pagerduty",
				Type: "pagerduty",
				Labels: map[string]string{
					"entity": "neworg",
					"team":   "siren-neworg",
				},
				Configurations: map[string]interface{}{
					"service_key": "000999",
				},
			},
			ExpectedID: uint64(4), // autoincrement in db side
		},
		{
			Description: "should return error if receiver is nil",
			ErrString:   "receiver domain is nil",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			err := s.repository.Create(s.ctx, tc.ReceiverToCreate)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
		})
	}
}

func (s *ReceiverRepositoryTestSuite) TestUpdate() {
	type testCase struct {
		Description      string
		ReceiverToUpdate *receiver.Receiver
		ExpectedID       uint64
		ErrString        string
	}

	var testCases = []testCase{
		{
			Description: "should update existing receiver",
			ReceiverToUpdate: &receiver.Receiver{
				ID:   2,
				Name: "alert-history-updated",
				Type: "http",
				Labels: map[string]string{
					"entity": "raystack",
				},
				Configurations: map[string]interface{}{
					"url": "http://siren.raystack.io/v2/alerts/cortex",
				},
			},
			ExpectedID: uint64(2),
		},
		{
			Description: "should return error not found if id not found",
			ReceiverToUpdate: &receiver.Receiver{
				ID: 1000,
			},
			ErrString: "receiver with id 1000 not found",
		},
		{
			Description: "should return error if receiver is nil",
			ErrString:   "receiver domain is nil",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			err := s.repository.Update(s.ctx, tc.ReceiverToUpdate)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
		})
	}
}

func (s *ReceiverRepositoryTestSuite) TestDelete() {
	type testCase struct {
		Description string
		IDToDelete  uint64
		ErrString   string
	}

	var testCases = []testCase{
		{
			Description: "should delete a receiver",
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

func TestReceiverRepository(t *testing.T) {
	suite.Run(t, new(ReceiverRepositoryTestSuite))
}
