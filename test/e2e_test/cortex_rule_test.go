package e2e_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/goto/siren/config"
	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/internal/server"
	"github.com/goto/siren/plugins"
	cortexv1plugin "github.com/goto/siren/plugins/providers/cortex/v1"
	sirenv1beta1 "github.com/goto/siren/proto/gotocompany/siren/v1beta1"
	"github.com/mcuadros/go-defaults"
	"github.com/stretchr/testify/suite"
)

const testEncryptionKey = "vHhhFjhFYULDHLyxPXDqVmBADbTxnQnd"

type CortexRuleTestSuite struct {
	suite.Suite
	cancelContext context.CancelFunc
	client        sirenv1beta1.SirenServiceClient
	cancelClient  func()
	appConfig     *config.Config
	testBench     *CortexTest
}

func (s *CortexRuleTestSuite) SetupTest() {

	apiPort, err := getFreePort()
	s.Require().Nil(err)

	s.appConfig = &config.Config{}

	defaults.SetDefaults(s.appConfig)

	s.appConfig.Log.Level = "error"
	s.appConfig.Service = server.Config{
		GRPC: server.GRPCConfig{
			Port: apiPort,
		},
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
	s.appConfig.Providers.Plugins = make(map[string]plugins.PluginConfig, 0)
	pathProject, _ := os.Getwd()
	rootProject := filepath.Dir(filepath.Dir(pathProject))
	s.appConfig.Providers.PluginPath = filepath.Join(rootProject, "plugin")
	s.appConfig.Providers.Plugins["cortex"] = plugins.PluginConfig{
		Handshake: plugins.HandshakeConfig{
			ProtocolVersion:  cortexv1plugin.Handshake.ProtocolVersion,
			MagicCookieKey:   cortexv1plugin.Handshake.MagicCookieKey,
			MagicCookieValue: cortexv1plugin.Handshake.MagicCookieValue,
		},
		ServiceConfig: map[string]interface{}{
			"webhook_base_api": fmt.Sprintf("http://host.docker.internal:%d/v1beta1/alerts/cortex", apiPort),
			"group_wait":       "1s",
			"repeat_interval":  "1s",
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelContext = cancel

	StartSirenServer(ctx, *s.appConfig)

	s.client, s.cancelClient, err = CreateClient(ctx, fmt.Sprintf("localhost:%d", apiPort))
	s.Require().NoError(err)

	bootstrapCortexTestData(&s.Suite, ctx, s.client, s.testBench.NginxHost)
}

func (s *CortexRuleTestSuite) TearDownTest() {
	s.cancelClient()

	// Clean tests
	err := s.testBench.CleanUp()
	s.Require().NoError(err)

	s.cancelContext()
}

func (s *CortexRuleTestSuite) TestRules() {
	ctx := context.Background()

	s.Run("initial state has no rule groups, upload rules and templates should return `testdata/cortex/expected-cortexrule-scenario-1.yaml`", func() {
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
