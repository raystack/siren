package e2e_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mcuadros/go-defaults"
	"github.com/odpf/siren/config"
	"github.com/odpf/siren/internal/server"
	sirenv1beta1 "github.com/odpf/siren/proto/odpf/siren/v1beta1"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/structpb"
)

type CortexAlertingTestSuite struct {
	suite.Suite
	client       sirenv1beta1.SirenServiceClient
	cancelClient func()
	appConfig    *config.Config
	testBench    *CortexTest
}

func (s *CortexAlertingTestSuite) SetupTest() {
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

	// setup custom cortex config
	s.appConfig.Providers.Cortex.WebhookBaseAPI = fmt.Sprintf("http://host.docker.internal:%d/v1beta1/alerts/cortex", apiPort)
	s.appConfig.Providers.Cortex.GroupWaitDuration = "1s"

	// enable worker
	s.appConfig.Notification.MessageHandler.Enabled = true
	s.appConfig.Notification.DLQHandler.Enabled = true

	StartSiren(*s.appConfig)

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

func (s *CortexAlertingTestSuite) TearDownTest() {
	s.cancelClient()
	// Clean tests
	err := s.testBench.CleanUp()
	s.Require().NoError(err)
}

func (s *CortexAlertingTestSuite) TestSendingNotification() {
	ctx := context.Background()

	s.Run("Triggering alert with matching subscription labels should trigger notification", func() {
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

		// add receiver odpf-http
		triggerAlertBody := `
		[
			{
				"state": "firing",
				"value": 1,
				"labels": {
					"severity": "WARNING",
					"team": "odpf",
					"service": "some-service",
					"environment": "integration"
				},
				"annotations": {
					"resource": "test_alert",
					"metricName": "test_alert",
					"metricValue": "1",
					"template": "alert_test"
				}
			}
		]`

		waitChan := make(chan struct{}, 1)

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			s.Assert().NoError(err)

			expectedBody := `{"environment":"integration","generatorUrl":"","groupKey":"{}:{severity=\"WARNING\"}","id":"cortex-39255dc96d0f642c","metricName":"test_alert","metricValue":"1","numAlertsFiring":1,"resource":"test_alert","routing_method":"subscribers","service":"some-service","severity":"WARNING","status":"firing","team":"odpf","template":"alert_test"}`
			s.Assert().Equal(expectedBody, string(body))
			close(waitChan)
		}))
		s.Require().Nil(err)
		defer testServer.Close()

		configs, err := structpb.NewStruct(map[string]interface{}{
			"url": testServer.URL,
		})
		s.Require().NoError(err)
		_, err = s.client.CreateReceiver(ctx, &sirenv1beta1.CreateReceiverRequest{
			Name: "odpf-http",
			Type: "http",
			Labels: map[string]string{
				"entity": "odpf",
				"kind":   "http",
			},
			Configurations: configs,
		})
		s.Require().NoError(err)

		_, err = s.client.CreateSubscription(ctx, &sirenv1beta1.CreateSubscriptionRequest{
			Urn:       "subscribe-http-one",
			Namespace: 1,
			Receivers: []*sirenv1beta1.ReceiverMetadata{
				{
					Id: 1,
				},
			},
			Match: map[string]string{
				"team":        "odpf",
				"service":     "some-service",
				"environment": "integration",
			},
		})
		s.Require().NoError(err)

		for {
			bodyBytes, err := triggerCortexAlert(s.testBench.NginxHost, "new-odpf-1", triggerAlertBody)
			s.Assert().NoError(err)
			if err != nil {
				break
			}

			if string(bodyBytes) != "the Alertmanager is not configured\n" {
				break
			}
		}

		<-waitChan
	})
}

func TestCortexAlertingTestSuite(t *testing.T) {
	suite.Run(t, new(CortexAlertingTestSuite))
}
