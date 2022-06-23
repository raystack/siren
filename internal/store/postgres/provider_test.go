package postgres_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/internal/store/postgres"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ProviderRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	db         *gorm.DB
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.ProviderRepository
}

func (s *ProviderRepositoryTestSuite) SetupSuite() {
	var err error

	logger := log.NewZap()
	s.db, s.pool, s.resource, err = newTestClient(logger)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	s.repository = postgres.NewProviderRepository(s.db)
}

func (s *ProviderRepositoryTestSuite) SetupTest() {
	var err error
	_, err = bootstrapProvider(s.db)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *ProviderRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ProviderRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ProviderRepositoryTestSuite) cleanup() error {
	queries := []string{
		"TRUNCATE TABLE providers RESTART IDENTITY CASCADE",
	}
	return execQueries(context.TODO(), s.db, queries)
}

func (s *ProviderRepositoryTestSuite) TestList() {
	type testCase struct {
		Description       string
		Filter            provider.Filter
		ExpectedProviders []*provider.Provider
		ErrString         string
	}

	var testCases = []testCase{
		{
			Description: "should get all providers",
			ExpectedProviders: []*provider.Provider{
				{
					ID:   1,
					Host: "http://cortex-ingress.odpf.io",
					URN:  "odpf-cortex",
					Name: "odpf-cortex",
					Type: "cortex",
				},
				{
					ID:   2,
					Host: "http://prometheus-ingress.odpf.io",
					URN:  "odpf-prometheus",
					Name: "odpf-prometheus",
					Type: "prometheus",
				},
			},
		},
		{
			Description: "should filter by urn",
			Filter: provider.Filter{
				URN: "odpf-prometheus",
			},
			ExpectedProviders: []*provider.Provider{
				{
					ID:   2,
					Host: "http://prometheus-ingress.odpf.io",
					URN:  "odpf-prometheus",
					Name: "odpf-prometheus",
					Type: "prometheus",
				},
			},
		},
		{
			Description: "should filter by type",
			Filter: provider.Filter{
				Type: "cortex",
			},
			ExpectedProviders: []*provider.Provider{
				{
					ID:   1,
					Host: "http://cortex-ingress.odpf.io",
					URN:  "odpf-cortex",
					Name: "odpf-cortex",
					Type: "cortex",
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
			if !cmp.Equal(got, tc.ExpectedProviders, cmpopts.IgnoreFields(provider.Provider{}, "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedProviders)
			}
		})
	}
}

func (s *ProviderRepositoryTestSuite) TestGet() {
	type testCase struct {
		Description      string
		PassedID         uint64
		ExpectedProvider *provider.Provider
		ErrString        string
	}

	var testCases = []testCase{
		{
			Description: "should get a provider",
			PassedID:    uint64(2),
			ExpectedProvider: &provider.Provider{
				ID:   2,
				Host: "http://prometheus-ingress.odpf.io",
				URN:  "odpf-prometheus",
				Name: "odpf-prometheus",
				Type: "prometheus",
			},
		},
		{
			Description: "should return not found if id not found",
			PassedID:    uint64(1000),
			ErrString:   "provider with id 1000 not found",
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
			if !cmp.Equal(got, tc.ExpectedProvider, cmpopts.IgnoreFields(provider.Provider{}, "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedProvider)
			}
		})
	}
}

func (s *ProviderRepositoryTestSuite) TestCreate() {
	type testCase struct {
		Description      string
		ProviderToCreate *provider.Provider
		ExpectedID       uint64
		ErrString        string
	}

	var testCases = []testCase{
		{
			Description: "should create a provider",
			ProviderToCreate: &provider.Provider{
				Host: "http://new-provider-ingress.odpf.io",
				URN:  "odpf-new-provider",
				Name: "odpf-new-provider",
				Type: "new-provider",
			},
			ExpectedID: uint64(3), // autoincrement in db side
		},
		{
			Description: "should return error duplicate if URN already exist",
			ProviderToCreate: &provider.Provider{
				Host: "http://newhostcortex",
				URN:  "odpf-cortex",
				Name: "odpf-cortex-new",
				Type: "cortex",
			},
			ErrString: "urn already exist",
		},
		{
			Description: "should return error if provider is nil",
			ErrString:   "provider domain is nil",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.Create(s.ctx, tc.ProviderToCreate)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedID) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedID)
			}
		})
	}
}

func (s *ProviderRepositoryTestSuite) TestUpdate() {
	type testCase struct {
		Description      string
		ProviderToUpdate *provider.Provider
		ExpectedID       uint64
		ErrString        string
	}

	var testCases = []testCase{
		{
			Description: "should update existing provider",
			ProviderToUpdate: &provider.Provider{
				ID:   1,
				Host: "http://new-provider-ingress.odpf.io",
				URN:  "odpf-new-provider",
				Name: "odpf-new-provider",
				Type: "new-provider",
			},
			ExpectedID: uint64(1),
		},
		{
			Description: "should return error duplicate if URN already exist",
			ProviderToUpdate: &provider.Provider{
				ID:   2,
				Host: "http://prometheus",
				URN:  "odpf-new-provider",
				Name: "odpf-prometheus",
				Type: "prometheus",
			},
			ErrString: "urn already exist",
		},
		{
			Description: "should return error not found if id not found",
			ProviderToUpdate: &provider.Provider{
				ID:   1000,
				Host: "http://prometheus",
				URN:  "odpf-new-provider",
				Name: "odpf-prometheus",
				Type: "prometheus",
			},
			ErrString: "provider with id 1000 not found",
		},
		{
			Description: "should return error if provider is nil",
			ErrString:   "provider domain is nil",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.Update(s.ctx, tc.ProviderToUpdate)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedID) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedID)
			}
		})
	}
}

func (s *ProviderRepositoryTestSuite) TestDelete() {
	type testCase struct {
		Description string
		IDToDelete  uint64
		ErrString   string
	}

	var testCases = []testCase{
		{
			Description: "should delete a provider",
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

func TestProviderRepository(t *testing.T) {
	suite.Run(t, new(ProviderRepositoryTestSuite))
}
