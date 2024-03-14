package postgres_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/goto/salt/dockertestx"
	"github.com/goto/salt/log"
	"github.com/goto/siren/core/receiver"
	"github.com/goto/siren/internal/store/postgres"
	"github.com/goto/siren/pkg/pgc"
	"github.com/ory/dockertest/v3"
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
	s.client, err = pgc.NewClient(logger, dbConfig)
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
					Name: "gotocompany-slack",
					Type: "slack",
					Labels: map[string]string{
						"entity":   "gotocompany,org-a,org-b",
						"severity": "warning",
						"team":     "infra",
					},
					Configurations: map[string]any{
						"token":     "xxxxxxxxxx",
						"workspace": "gotocompany",
					},
				},
				{
					ID:   2,
					Name: "alert-history",
					Type: "http",
					Labels: map[string]string{
						"entity": "gotocompany,org-a,org-b,org-c",
						"team":   "infra",
					},
					Configurations: map[string]any{
						"url": "http://siren.gotocompany.com/v1beta1/alerts/cortex/1",
					},
				},
				{
					ID:   3,
					Name: "gotocompany_pagerduty",
					Type: "pagerduty",
					Labels: map[string]string{
						"entity": "gotocompany",
						"team":   "siren-gotocompany",
					},
					Configurations: map[string]any{
						"service_key": "1212121212121212121212121",
					},
				},
				{
					ID:   4,
					Name: "gotocompany-slack",
					Type: "slack_channel",
					Labels: map[string]string{
						"org":      "gotocompany,org-a,org-b",
						"team":     "infra",
						"severity": "critical",
					},
					Configurations: map[string]any{
						"channel_name": "test-pilot-alert",
					},
					ParentID: 1,
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
						"entity": "gotocompany,org-a,org-b,org-c",
						"team":   "infra",
					},
					Configurations: map[string]any{
						"url": "http://siren.gotocompany.com/v1beta1/alerts/cortex/1",
					},
				},
				{
					ID:   3,
					Name: "gotocompany_pagerduty",
					Type: "pagerduty",
					Labels: map[string]string{
						"entity": "gotocompany",
						"team":   "siren-gotocompany",
					},
					Configurations: map[string]any{
						"service_key": "1212121212121212121212121",
					},
				},
			},
		},
		{
			Description: "should get filtered receivers with labels",
			Filter: receiver.Filter{
				MultipleLabels: []map[string]string{
					{
						"team":     "infra",
						"severity": "warning",
					},
				},
			},
			ExpectedReceivers: []receiver.Receiver{
				{
					ID:   1,
					Name: "gotocompany-slack",
					Type: "slack",
					Labels: map[string]string{
						"entity":   "gotocompany,org-a,org-b",
						"team":     "infra",
						"severity": "warning",
					},
					Configurations: map[string]any{
						"token":     "xxxxxxxxxx",
						"workspace": "gotocompany",
					},
				},
			},
		},
		{
			Description: "should get all receivers with parent configs merged if expanded",
			Filter: receiver.Filter{
				Expanded: true,
			},
			ExpectedReceivers: []receiver.Receiver{
				{
					ID:   1,
					Name: "gotocompany-slack",
					Type: "slack",
					Labels: map[string]string{
						"entity":   "gotocompany,org-a,org-b",
						"severity": "warning",
						"team":     "infra",
					},
					Configurations: map[string]any{
						"token":     "xxxxxxxxxx",
						"workspace": "gotocompany",
					},
				},
				{
					ID:   2,
					Name: "alert-history",
					Type: "http",
					Labels: map[string]string{
						"entity": "gotocompany,org-a,org-b,org-c",
						"team":   "infra",
					},
					Configurations: map[string]any{
						"url": "http://siren.gotocompany.com/v1beta1/alerts/cortex/1",
					},
				},
				{
					ID:   3,
					Name: "gotocompany_pagerduty",
					Type: "pagerduty",
					Labels: map[string]string{
						"entity": "gotocompany",
						"team":   "siren-gotocompany",
					},
					Configurations: map[string]any{
						"service_key": "1212121212121212121212121",
					},
				},
				{
					ID:   4,
					Name: "gotocompany-slack",
					Type: "slack_channel",
					Labels: map[string]string{
						"org":      "gotocompany,org-a,org-b",
						"team":     "infra",
						"severity": "critical",
					},
					Configurations: map[string]any{
						"token":        "xxxxxxxxxx",
						"workspace":    "gotocompany",
						"channel_name": "test-pilot-alert",
					},
					ParentID: 1,
				},
			},
		},
		{
			Description: "should get filtered receivers with list of ids with parent configs merged",
			Filter: receiver.Filter{
				ReceiverIDs: []uint64{3, 4},
				Expanded:    true,
			},
			ExpectedReceivers: []receiver.Receiver{
				{
					ID:   4,
					Name: "gotocompany-slack",
					Type: "slack_channel",
					Labels: map[string]string{
						"org":      "gotocompany,org-a,org-b",
						"team":     "infra",
						"severity": "critical",
					},
					Configurations: map[string]any{
						"token":        "xxxxxxxxxx",
						"workspace":    "gotocompany",
						"channel_name": "test-pilot-alert",
					},
					ParentID: 1,
				},
				{
					ID:   3,
					Name: "gotocompany_pagerduty",
					Type: "pagerduty",
					Labels: map[string]string{
						"entity": "gotocompany",
						"team":   "siren-gotocompany",
					},
					Configurations: map[string]any{
						"service_key": "1212121212121212121212121",
					},
				},
			},
		},
		{
			Description: "should get filtered receivers with labels with parent configs merged",
			Filter: receiver.Filter{
				MultipleLabels: []map[string]string{
					{
						"org":      "gotocompany,org-a,org-b",
						"team":     "infra",
						"severity": "critical",
					},
				},
				Expanded: true,
			},
			ExpectedReceivers: []receiver.Receiver{
				{
					ID:   4,
					Name: "gotocompany-slack",
					Type: "slack_channel",
					Labels: map[string]string{
						"org":      "gotocompany,org-a,org-b",
						"team":     "infra",
						"severity": "critical",
					},
					Configurations: map[string]any{
						"token":        "xxxxxxxxxx",
						"workspace":    "gotocompany",
						"channel_name": "test-pilot-alert",
					},
					ParentID: 1,
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
			if diff := cmp.Diff(got, tc.ExpectedReceivers, cmpopts.IgnoreFields(receiver.Receiver{}, "CreatedAt", "UpdatedAt")); diff != "" {
				s.T().Fatalf("got diff %+v", diff)
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
		WithExpand       bool
	}

	var testCases = []testCase{
		{
			Description: "should get a receivers",
			PassedID:    3,
			ExpectedReceiver: &receiver.Receiver{
				ID:   3,
				Name: "gotocompany_pagerduty",
				Type: "pagerduty",
				Labels: map[string]string{
					"entity": "gotocompany",
					"team":   "siren-gotocompany",
				},
				Configurations: map[string]any{
					"service_key": "1212121212121212121212121",
				},
			},
			WithExpand: false,
		},
		{
			Description: "should return not found if id not found",
			PassedID:    1000,
			ErrString:   "receiver with id 1000 not found",
			WithExpand:  false,
		},
		{
			Description: "should get a receiver with parent configurations when with expand set to true",
			PassedID:    4,
			ExpectedReceiver: &receiver.Receiver{
				ID:   4,
				Name: "gotocompany-slack",
				Type: "slack_channel",
				Labels: map[string]string{
					"org":      "gotocompany,org-a,org-b",
					"team":     "infra",
					"severity": "critical",
				},
				Configurations: map[string]any{
					"token":        "xxxxxxxxxx",
					"workspace":    "gotocompany",
					"channel_name": "test-pilot-alert",
				},
				ParentID: 1,
			},
			WithExpand: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			filter := receiver.Filter{
				Expanded: tc.WithExpand,
			}
			got, err := s.repository.Get(s.ctx, tc.PassedID, filter)
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
				Configurations: map[string]any{
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
					"entity": "gotocompany",
				},
				Configurations: map[string]any{
					"url": "http://siren.gotocompany.com/v2/alerts/cortex",
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

func (s *ReceiverRepositoryTestSuite) TestPatchLabels() {
	var (
		receiverID uint64 = 2
	)
	type testCase struct {
		Description      string
		ReceiverToPatch  *receiver.Receiver
		ExpectedReceiver *receiver.Receiver
		ErrString        string
	}

	var testCases = []testCase{
		{
			Description: "should patch existing receiver",
			ReceiverToPatch: &receiver.Receiver{
				ID: receiverID,
				Labels: map[string]string{
					"foo": "newbar",
					"a":   "b",
				},
			},
			ExpectedReceiver: &receiver.Receiver{
				ID:   receiverID,
				Type: receiver.TypeHTTP,
				Name: "alert-history",
				Labels: map[string]string{
					"foo": "newbar",
					"a":   "b",
				},
				Configurations: map[string]any{
					"url": "http://siren.gotocompany.com/v1beta1/alerts/cortex/1",
				},
			},
		},
		{
			Description:     "should return error not found if receiver is nil",
			ReceiverToPatch: nil,
			ErrString:       "request is not valid",
		},
	}

	for _, tc := range testCases {
		rcv := tc.ReceiverToPatch
		s.Run(tc.Description, func() {
			err := s.repository.PatchLabels(s.ctx, rcv)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if diff := cmp.Diff(rcv, tc.ExpectedReceiver, cmpopts.IgnoreFields(receiver.Receiver{}, "CreatedAt", "UpdatedAt")); diff != "" {
				s.T().Fatalf("got diff %+v", diff)
			}
		})
	}
}

func TestReceiverRepository(t *testing.T) {
	suite.Run(t, new(ReceiverRepositoryTestSuite))
}
