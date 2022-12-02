package e2e_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/mcuadros/go-defaults"
	"github.com/odpf/siren/config"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/internal/server"
	"github.com/odpf/siren/pkg/telemetry"
	sirenv1beta1 "github.com/odpf/siren/proto/odpf/siren/v1beta1"
	"github.com/stretchr/testify/suite"
)

const testEncryptionKey = "vHhhFjhFYULDHLyxPXDqVmBADbTxnQnd"

type CortexRuleTestSuite struct {
	suite.Suite
	client       sirenv1beta1.SirenServiceClient
	cancelClient func()
	appConfig    *config.Config
	testBench    *CortexTest
}

func (s *CortexRuleTestSuite) SetupTest() {

	apiPort, err := getFreePort()
	s.Require().Nil(err)

	s.appConfig = &config.Config{
		Log: config.Log{
			Level: "debug",
		},
		Telemetry: telemetry.Config{
			Debug: "",
		},
		Service: server.Config{
			Port:          apiPort,
			EncryptionKey: testEncryptionKey,
		},
		Notification: notification.Config{
			MessageHandler: notification.HandlerConfig{
				Enabled: false,
			},
			DLQHandler: notification.HandlerConfig{
				Enabled: false,
			},
		},
	}

	defaults.SetDefaults(s.appConfig)

	s.testBench, err = InitCortexEnvironment(s.appConfig)
	s.Require().NoError(err)

	// TODO host.docker.internal only works for docker-desktop to call a service in host (siren)
	s.appConfig.Providers.Cortex.WebhookBaseAPI = fmt.Sprintf("http://host.docker.internal:%d/v1beta1/alerts/cortex", apiPort)
	s.appConfig.Providers.Cortex.GroupWaitDuration = "1s"
	StartSirenServer(*s.appConfig)

	ctx := context.Background()
	s.client, s.cancelClient, err = CreateClient(ctx, fmt.Sprintf("localhost:%d", apiPort))
	s.Require().NoError(err)

	bootstrapCortexTestData(&s.Suite, ctx, s.client, s.testBench.NginxHost)
}

func (s *CortexRuleTestSuite) TearDownTest() {
	s.cancelClient()
	// Clean tests
	err := s.testBench.CleanUp()
	s.Require().NoError(err)
}

func (s *CortexRuleTestSuite) TestRules() {
	ctx := context.Background()

	s.Run("1. initial state has no rule groups, upload rules and templates should return `testdata/cortex/expected-cortexrule-scenario-1.yaml`", func() {
		err := uploadTemplate(ctx, s.client, "testdata/cortex/template-rule-sample-1.yaml")
		s.Require().NoError(err)
		err = uploadTemplate(ctx, s.client, "testdata/cortex/template-rule-sample-2.yaml")
		s.Require().NoError(err)

		err = uploadRule(ctx, s.client, "testdata/cortex/rule-sample-scenario-1.yaml")
		s.Require().NoError(err)

		bodyBytes, err := fetchCortexRules(s.testBench.NginxHost, "fake")
		s.Require().NoError(err)
		expectedScenarioCortexRule, err := os.ReadFile("testdata/cortex/expected-cortexrule-scenario-1.yaml")
		s.Require().NoError(err)

		s.Assert().Empty(diffYaml(bodyBytes, expectedScenarioCortexRule))
	})

	s.Run("2. initial state `testdata/cortex/expected-cortexrule-scenario-1.yaml`, disabling one rule, should return `testdata/cortex/expected-cortexrule-scenario-2.yaml`", func() {
		err := uploadRule(ctx, s.client, "testdata/cortex/rule-sample-scenario-2.yaml")
		s.Require().NoError(err)

		bodyBytes, err := fetchCortexRules(s.testBench.NginxHost, "fake")
		s.Require().NoError(err)
		expectedScenarioCortexRule, err := os.ReadFile("testdata/cortex/expected-cortexrule-scenario-2.yaml")
		s.Require().NoError(err)

		s.Assert().Empty(diffYaml(bodyBytes, expectedScenarioCortexRule))
	})
}

func TestCortexRuleTestSuite(t *testing.T) {
	suite.Run(t, new(CortexRuleTestSuite))
}
