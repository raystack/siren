package e2e_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mcuadros/go-defaults"
	"github.com/raystack/siren/config"
	"github.com/raystack/siren/core/notification"
	"github.com/raystack/siren/internal/server"
	sirenv1beta1 "github.com/raystack/siren/proto/raystack/siren/v1beta1"
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

	s.appConfig = &config.Config{}

	defaults.SetDefaults(s.appConfig)

	s.appConfig.Log.Level = "error"
	s.appConfig.Service = server.Config{
		Port:          apiPort,
		EncryptionKey: testEncryptionKey,
	}
	s.appConfig.Notification = notification.Config{
		MessageHandler: notification.HandlerConfig{
			Enabled: false,
		},
		DLQHandler: notification.HandlerConfig{
			Enabled: false,
		},
	}
	s.appConfig.Telemetry.Debug = ""
	s.appConfig.Telemetry.EnableNewrelic = false
	s.appConfig.Telemetry.EnableOtelAgent = false

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

	s.Run("initial state alert config not set, add a namespace will set config for the provider tenant", func() {
		_, err := s.client.CreateNamespace(ctx, &sirenv1beta1.CreateNamespaceRequest{
			Name:        "new-raystack-1",
			Urn:         "new-raystack-1",
			Provider:    1,
			Credentials: nil,
			Labels: map[string]string{
				"key1": "value1",
			},
		})
		s.Require().NoError(err)

		bodyBytes, err := fetchCortexAlertmanagerConfig(s.testBench.NginxHost, "new-raystack-1")
		s.Require().NoError(err)

		expectedScenarioCortexAM, err := os.ReadFile("testdata/cortex/expected-cortexamconfig-scenario.yaml")
		s.Require().NoError(err)

		s.Assert().Empty(cmp.Diff(bodyBytes, expectedScenarioCortexAM))
	})
}

func TestCortexNamespaceTestSuite(t *testing.T) {
	suite.Run(t, new(CortexNamespaceTestSuite))
}
