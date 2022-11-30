package e2e_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/mcuadros/go-defaults"
	"github.com/odpf/siren/config"
	"github.com/odpf/siren/internal/server"
	sirenv1beta1 "github.com/odpf/siren/proto/odpf/siren/v1beta1"
	"github.com/stretchr/testify/suite"
)

type CortexNamespaceTestSuite struct {
	suite.Suite
	client       sirenv1beta1.SirenServiceClient
	cancelClient func()
	appConfig    *config.Config
	testBench    *CortexTest
}

func (s *CortexNamespaceTestSuite) SetupTest() {
	apiPort, err := getFreePort()
	s.Require().Nil(err)

	s.appConfig = &config.Config{
		Log: config.Log{
			Level: "debug",
		},
		Service: server.Config{
			Port:          apiPort,
			EncryptionKey: testEncryptionKey,
		},
	}

	defaults.SetDefaults(s.appConfig)

	s.testBench, err = InitCortexEnvironment(s.appConfig)
	s.Require().NoError(err)

	// TODO host.docker.internal only works for docker-desktop to call a service in host (siren)
	s.appConfig.Providers.Cortex.WebhookBaseAPI = "http://host.docker.internal:8080/v1beta1/alerts/cortex"
	StartSirenServer(*s.appConfig)

	ctx := context.Background()
	s.client, s.cancelClient, err = CreateClient(ctx, fmt.Sprintf("localhost:%d", apiPort))
	s.Require().NoError(err)

	_, err = s.client.CreateProvider(ctx, &sirenv1beta1.CreateProviderRequest{
		Host: fmt.Sprintf("http://%s", s.testBench.NginxHost),
		Urn:  "cortex-test",
		Name: "cortex-test",
		Type: "cortex",
	})
	s.Require().NoError(err)
}

func (s *CortexNamespaceTestSuite) TearDownTest() {
	s.cancelClient()
	// Clean tests
	err := s.testBench.CleanUp()
	s.Require().NoError(err)
}

func (s *CortexNamespaceTestSuite) TestNamespace() {
	ctx := context.Background()

	s.Run("Initial state alert config not set, add a namespace will set config for the provider tenant", func() {
		_, err := s.client.CreateNamespace(ctx, &sirenv1beta1.CreateNamespaceRequest{
			Name:        "new-odpf-1",
			Urn:         "new-odpf-1",
			Provider:    1,
			Credentials: nil,
			Labels: map[string]string{
				"key1": "value1",
			},
		})
		s.Require().NoError(err)

		bodyBytes, err := fetchCortexAlertmanagerConfig(s.testBench.NginxHost, "new-odpf-1")
		s.Require().NoError(err)

		expectedScenarioCortexAM, err := os.ReadFile("testdata/cortex/expected-cortexamconfig-scenario.yaml")
		s.Require().NoError(err)

		s.Assert().Empty(diffYaml(bodyBytes, expectedScenarioCortexAM))
	})
}

func TestCortexNamespaceTestSuite(t *testing.T) {
	suite.Run(t, new(CortexNamespaceTestSuite))
}
