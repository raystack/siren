package postgres_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/internal/store/postgres"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
)

type TemplateRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	client     *postgres.Client
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.TemplateRepository
}

func (s *TemplateRepositoryTestSuite) SetupSuite() {
	var err error

	logger := log.NewZap()
	s.client, s.pool, s.resource, err = newTestClient(logger)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	s.repository = postgres.NewTemplateRepository(s.client)
}

func (s *TemplateRepositoryTestSuite) SetupTest() {
	var err error
	_, err = bootstrapTemplate(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *TemplateRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *TemplateRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *TemplateRepositoryTestSuite) cleanup() error {
	queries := []string{
		"TRUNCATE TABLE templates RESTART IDENTITY CASCADE",
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *TemplateRepositoryTestSuite) TestList() {
	type testCase struct {
		Description       string
		Filter            template.Filter
		ExpectedTemplates []template.Template
		ErrString         string
	}

	var testCases = []testCase{
		{
			Description: "should get all templates",
			ExpectedTemplates: []template.Template{
				{
					ID:   1,
					Name: "zookeeper-pending-syncs",
					Body: "- alert: zookeeper pending syncs warning\n  expr: avg by (host, environment) (zookeeper_pending_syncs) > [[.warning]]\n  for: '[[.for]]'\n  labels:\n    alertname: zookeeper pending syncs on host {{ $labels.host }} is greater than \"[[.warning]]\"\n    environment: '{{ $labels.environment }}'\n    severity: WARNING\n    team: '[[.team]]'\n  annotations:\n    metric_name: zookeeper_pending_syncs\n    metric_value: '{{ $value }}'\n    resource: '{{ $labels.host }}'\n    summary: zookeeper pending sync on host {{ $labels.host }} is {{ $value }}\n    template: zookeeper-pending-syncs\n- alert: zookeeper pending sync critical\n  expr: avg by (host, environment) (zookeeper_pending_syncs) > [[.critical]]\n  for: '[[.for]]'\n  labels:\n    alertname: zookeeper pending syncs on host {{ $labels.host }} is greater than \"[[.critical]]\"\n    environment: '{{ $labels.environment }}'\n    severity: CRITICAL\n    team: '[[.team]]'\n  annotations:\n    metric_name: zookeeper_pending_syncs\n    metric_value: '{{ $value }}'\n    resource: '{{ $labels.host }}'\n    summary: zookeeper oustanding requests on host {{ $labels.host }} is {{ $value }}\n    template: zookeeper-pending-syncs\n",
					Tags: []string{
						"zookeeper",
					},
					Variables: []template.Variable{
						{
							Name:        "for",
							Type:        "string",
							Default:     "5m",
							Description: "For eg 5m, 2h; Golang duration format",
						},
						{
							Name:    "warning",
							Type:    "int",
							Default: "10",
						},
						{
							Name:    "critical",
							Type:    "int",
							Default: "100",
						},
						{
							Name:        "team",
							Type:        "string",
							Default:     "odpf",
							Description: "For eg team name which the alert should go to",
						},
					},
				},
				{
					ID:   2,
					Name: "kafka-under-replicated-partitions",
					Body: "- alert: kafka under replicated partitions warning\n  expr: sum by (host, environment) (v2_jolokia_kafka_server_ReplicaManager_UnderReplicatedPartitionsValue) > [[.warning]]\n  for: '[[.for]]'\n  labels:\n    alertname: number of under replicated partitions on host {{ $labels.host }} is {{ $value }}\n    environment: '{{ $labels.environment }}'\n    severity: WARNING\n    team: '[[.team]]'\n  annotations:\n    metric_name: kafka_under_replicated_partitions\n    metric_value: '{{ $value }}'\n    resource: '{{ $labels.host }}'\n    summary: under replicated partitions on host {{ $labels.host }} is {{ $value }} is greather than [[.warning]]\n    template: kafka-under-replicated-partitions\n",
					Tags: []string{
						"kafka",
					},
					Variables: []template.Variable{
						{
							Name:        "for",
							Type:        "string",
							Default:     "10m",
							Description: "For eg 5m, 2h; Golang duration format",
						},
						{
							Name:    "warning",
							Type:    "int",
							Default: "0",
						},
						{
							Name:        "team",
							Type:        "string",
							Default:     "odpf",
							Description: "For eg team name which the alert should go to",
						},
					},
				},
			},
		},
		{
			Description: "should get filtered templates",
			Filter: template.Filter{
				Tag: "zookeeper",
			},
			ExpectedTemplates: []template.Template{
				{
					ID:   1,
					Name: "zookeeper-pending-syncs",
					Body: "- alert: zookeeper pending syncs warning\n  expr: avg by (host, environment) (zookeeper_pending_syncs) > [[.warning]]\n  for: '[[.for]]'\n  labels:\n    alertname: zookeeper pending syncs on host {{ $labels.host }} is greater than \"[[.warning]]\"\n    environment: '{{ $labels.environment }}'\n    severity: WARNING\n    team: '[[.team]]'\n  annotations:\n    metric_name: zookeeper_pending_syncs\n    metric_value: '{{ $value }}'\n    resource: '{{ $labels.host }}'\n    summary: zookeeper pending sync on host {{ $labels.host }} is {{ $value }}\n    template: zookeeper-pending-syncs\n- alert: zookeeper pending sync critical\n  expr: avg by (host, environment) (zookeeper_pending_syncs) > [[.critical]]\n  for: '[[.for]]'\n  labels:\n    alertname: zookeeper pending syncs on host {{ $labels.host }} is greater than \"[[.critical]]\"\n    environment: '{{ $labels.environment }}'\n    severity: CRITICAL\n    team: '[[.team]]'\n  annotations:\n    metric_name: zookeeper_pending_syncs\n    metric_value: '{{ $value }}'\n    resource: '{{ $labels.host }}'\n    summary: zookeeper oustanding requests on host {{ $labels.host }} is {{ $value }}\n    template: zookeeper-pending-syncs\n",
					Tags: []string{
						"zookeeper",
					},
					Variables: []template.Variable{
						{
							Name:        "for",
							Type:        "string",
							Default:     "5m",
							Description: "For eg 5m, 2h; Golang duration format",
						},
						{
							Name:    "warning",
							Type:    "int",
							Default: "10",
						},
						{
							Name:    "critical",
							Type:    "int",
							Default: "100",
						},
						{
							Name:        "team",
							Type:        "string",
							Default:     "odpf",
							Description: "For eg team name which the alert should go to",
						},
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
			if !cmp.Equal(got, tc.ExpectedTemplates, cmpopts.IgnoreFields(template.Template{}, "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedTemplates)
			}
		})
	}
}

func (s *TemplateRepositoryTestSuite) TestGetByName() {
	type testCase struct {
		Description      string
		Name             string
		ExpectedTemplate *template.Template
		ErrString        string
	}

	var testCases = []testCase{
		{
			Description: "should get by name",
			Name:        "zookeeper-pending-syncs",
			ExpectedTemplate: &template.Template{
				ID:   1,
				Name: "zookeeper-pending-syncs",
				Body: "- alert: zookeeper pending syncs warning\n  expr: avg by (host, environment) (zookeeper_pending_syncs) > [[.warning]]\n  for: '[[.for]]'\n  labels:\n    alertname: zookeeper pending syncs on host {{ $labels.host }} is greater than \"[[.warning]]\"\n    environment: '{{ $labels.environment }}'\n    severity: WARNING\n    team: '[[.team]]'\n  annotations:\n    metric_name: zookeeper_pending_syncs\n    metric_value: '{{ $value }}'\n    resource: '{{ $labels.host }}'\n    summary: zookeeper pending sync on host {{ $labels.host }} is {{ $value }}\n    template: zookeeper-pending-syncs\n- alert: zookeeper pending sync critical\n  expr: avg by (host, environment) (zookeeper_pending_syncs) > [[.critical]]\n  for: '[[.for]]'\n  labels:\n    alertname: zookeeper pending syncs on host {{ $labels.host }} is greater than \"[[.critical]]\"\n    environment: '{{ $labels.environment }}'\n    severity: CRITICAL\n    team: '[[.team]]'\n  annotations:\n    metric_name: zookeeper_pending_syncs\n    metric_value: '{{ $value }}'\n    resource: '{{ $labels.host }}'\n    summary: zookeeper oustanding requests on host {{ $labels.host }} is {{ $value }}\n    template: zookeeper-pending-syncs\n",
				Tags: []string{
					"zookeeper",
				},
				Variables: []template.Variable{
					{
						Name:        "for",
						Type:        "string",
						Default:     "5m",
						Description: "For eg 5m, 2h; Golang duration format",
					},
					{
						Name:    "warning",
						Type:    "int",
						Default: "10",
					},
					{
						Name:    "critical",
						Type:    "int",
						Default: "100",
					},
					{
						Name:        "team",
						Type:        "string",
						Default:     "odpf",
						Description: "For eg team name which the alert should go to",
					},
				},
			},
		},
		{
			Description: "should return not found if name does not exist",
			Name:        "random",
			ErrString:   "template with name \"random\" not found",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.GetByName(s.ctx, tc.Name)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedTemplate, cmpopts.IgnoreFields(template.Template{}, "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedTemplate)
			}
		})
	}
}

func (s *TemplateRepositoryTestSuite) TestUpsert() {
	type testCase struct {
		Description      string
		TemplateToUpsert *template.Template
		ExpectedID       uint64
		ErrString        string
	}

	var testCases = []testCase{
		{
			Description: "should create non existent template",
			TemplateToUpsert: &template.Template{
				Name: "new-template",
				Body: "template body",
				Tags: []string{
					"unknown",
				},
				Variables: []template.Variable{},
			},
			ExpectedID: uint64(3),
		},
		{
			Description: "should update the existing template",
			TemplateToUpsert: &template.Template{
				ID:        1,
				Name:      "zookeeper-pending-syncs",
				Body:      "new body",
				Tags:      []string{},
				Variables: []template.Variable{},
			},
			ExpectedID: uint64(1),
		},
		{
			Description: "should return conflict error if try to update same name with different id",
			TemplateToUpsert: &template.Template{
				ID:        234,
				Name:      "zookeeper-pending-syncs",
				Body:      "new body",
				Tags:      []string{},
				Variables: []template.Variable{},
			},
			ErrString: "name already exist",
		},
		{
			Description: "should return error if template is nil",
			ErrString:   "template domain is nil",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			err := s.repository.Upsert(s.ctx, tc.TemplateToUpsert)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
		})
	}
}

func (s *TemplateRepositoryTestSuite) TestDelete() {
	type testCase struct {
		Description  string
		NameToDelete string
		ErrString    string
	}

	var testCases = []testCase{
		{
			Description:  "should delete a template",
			NameToDelete: "zookeeper-pending-syncs",
		},
		{
			Description:  "should return nil if name does not exist",
			NameToDelete: "random",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			err := s.repository.Delete(s.ctx, tc.NameToDelete)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
		})
	}
}

func TestTemplateRepository(t *testing.T) {
	suite.Run(t, new(TemplateRepositoryTestSuite))
}
