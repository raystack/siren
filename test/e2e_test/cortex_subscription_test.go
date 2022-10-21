package e2e_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/odpf/siren/config"
	"github.com/odpf/siren/internal/server"
	sirenv1beta1 "github.com/odpf/siren/proto/odpf/siren/v1beta1"
	"github.com/stretchr/testify/suite"
)

type CortexSubscriptionTestSuite struct {
	suite.Suite
	client       sirenv1beta1.SirenServiceClient
	cancelClient func()
	appConfig    *config.Config
	testBench    *CortexTest
}

func (s *CortexSubscriptionTestSuite) SetupTest() {

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

	s.testBench, err = InitCortexEnvironment(s.appConfig)
	s.Require().NoError(err)

	// override address to use alertmanager only
	s.appConfig.Cortex.Address = fmt.Sprintf("http://%s", s.testBench.CortexAMHost)
	StartSiren(*s.appConfig)

	ctx := context.Background()
	s.client, s.cancelClient, err = CreateClient(ctx, fmt.Sprintf("localhost:%d", apiPort))
	s.Require().NoError(err)

	bootstrapCortexTestData(&s.Suite, ctx, s.client, s.appConfig.Cortex.Address)
}

func (s *CortexSubscriptionTestSuite) TearDownTest() {
	s.cancelClient()
	// Clean tests
	err := s.testBench.CleanUp()
	s.Require().NoError(err)
}

func (s *CortexSubscriptionTestSuite) TestSubscriptions() {
	ctx := context.Background()

	s.Run("1. initial state has no alert subscriptions, add a subscription to odpf-http should return `testdata/cortex/expected-cortexamconfig-scenario-1.yaml`", func() {
		_, err := s.client.CreateSubscription(ctx, &sirenv1beta1.CreateSubscriptionRequest{
			Urn:       "subscribe-http-one",
			Namespace: 1,
			Receivers: []*sirenv1beta1.ReceiverMetadata{
				{
					Id: 1,
				},
			},
			Match: map[string]string{
				"team":        "odpf-platform",
				"environment": "integration",
			},
		})
		s.Require().NoError(err)

		bodyBytes, err := fetchCortexAlertmanagerConfig(s.testBench.CortexAMHost)
		s.Require().NoError(err)

		expectedScenarioCortexAM, err := os.ReadFile("testdata/cortex/expected-cortexamconfig-scenario-1.yaml")
		s.Require().NoError(err)

		s.Assert().Empty(diffYaml(bodyBytes, expectedScenarioCortexAM))
	})

	s.Run("2. initial state `testdata/cortex/expected-cortexamconfig-scenario-1.yaml`, updating subscription, should return `testdata/cortex/expected-cortexamconfig-scenario-1-updated.yaml`", func() {
		_, err := s.client.UpdateSubscription(ctx, &sirenv1beta1.UpdateSubscriptionRequest{
			Id:        1,
			Urn:       "subscribe-http-one-updated",
			Namespace: 1,
			Receivers: []*sirenv1beta1.ReceiverMetadata{
				{
					Id: 1,
				},
			},
			Match: map[string]string{
				"team":        "odpf-platform-updated",
				"environment": "integration-updated",
			},
		})
		s.Require().NoError(err)

		bodyBytes, err := fetchCortexAlertmanagerConfig(s.testBench.CortexAMHost)
		s.Require().NoError(err)

		expectedScenarioCortexAM, err := os.ReadFile("testdata/cortex/expected-cortexamconfig-scenario-1-updated.yaml")
		s.Require().NoError(err)

		s.Assert().Empty(diffYaml(bodyBytes, expectedScenarioCortexAM))
	})

	s.Run("3. `testdata/cortex/expected-cortexamconfig-scenario-1-updated.yaml`, updating subscription, should return `testdata/cortex/expected-cortexamconfig-scenario-2.yaml`", func() {
		_, err := s.client.DeleteSubscription(ctx, &sirenv1beta1.DeleteSubscriptionRequest{
			Id: 1,
		})
		s.Require().NoError(err)

		bodyBytes, err := fetchCortexAlertmanagerConfig(s.testBench.CortexAMHost)
		s.Require().NoError(err)

		expectedScenarioCortexAM, err := os.ReadFile("testdata/cortex/expected-cortexamconfig-scenario-2.yaml")
		s.Require().NoError(err)

		s.Assert().Empty(diffYaml(bodyBytes, expectedScenarioCortexAM))
	})
}

func TestCortexSubscriptionTestSuite(t *testing.T) {
	suite.Run(t, new(CortexSubscriptionTestSuite))
}
