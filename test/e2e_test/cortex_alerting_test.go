package e2e_test

import (
	"bytes"
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
	// TODO host.docker.internal only works for docker-desktop to call a service in host (siren)
	s.appConfig.Providers.Cortex.WebhookBaseAPI = fmt.Sprintf("http://test:%d/v1beta1/alerts/cortex", apiPort)
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

func (s *CortexAlertingTestSuite) TestAlerting() {
	ctx := context.Background()

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

	s.Run("Triggering cortex alert with matching subscription labels should trigger notification", func() {
		configs, err := structpb.NewStruct(map[string]interface{}{
			//
			"url": "http://some-url",
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
			Urn:       "subscribe-http",
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

	})
}

func (s *CortexAlertingTestSuite) TestIncomingHookAPI() {
	ctx := context.Background()

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

	s.Run("Incoming alert in alerts hook API with matching subscription labels should trigger notification", func() {
		waitChan := make(chan struct{}, 1)

		// add receiver odpf-http
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			s.Assert().NoError(err)

			expectedBody := `{"alertname":"some alert name","environment":"integration","generatorUrl":"","groupKey":"{}/{environment=\"integration\",team=\"odpf\"}:{}","id":"cortex-684c979dcb5ffb96","key1":"value1","key2":"value2","metricName":"test_alert","metricValue":"1","numAlertsFiring":1,"resource":"test_alert","routing_method":"subscribers","service":"some-service","severity":"WARNING","status":"firing","summary":"this is test alert","team":"odpf","template":"alert_test"}`
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
			Urn:       "subscribe-http",
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

		triggerAlertBody := `
		{
			"receiver": "http_subscribe-http-receiver-notification_receiverId_2_idx_0",
			"status": "firing",
			"alerts": [
				{
					"status": "firing",
					"labels": {
						"key1": "value1",
						"key2": "value2",
						"severity": "WARNING",
						"alertname": "some alert name",
						"summary": "this is test alert",
						"service": "some-service",
						"environment": "integration",
						"team": "odpf"
					},
					"annotations": {
						"metricName": "test_alert",
						"metricValue": "1",
						"resource": "test_alert",
						"template": "alert_test",
						"summary": "this is test alert"
					},
					"startsAt": "2022-10-06T03:39:19.2964655Z",
					"endsAt": "0001-01-01T00:00:00Z",
					"generatorURL": "",
					"fingerprint": "684c979dcb5ffb96"
				}
			],
			"groupLabels": {},
			"commonLabels": {
				"environment": "integration",
				"team": "odpf"
			},
			"commonAnnotations": {
				"metricName": "test_alert",
				"metricValue": "1",
				"resource": "test_alert",
				"template": "alert_test"
			},
			"externalURL": "/api/prom/alertmanager",
			"version": "4",
			"groupKey": "{}/{environment=\"integration\",team=\"odpf\"}:{}",
			"truncatedAlerts": 0
		}`

		res, err := http.DefaultClient.Post(fmt.Sprintf("http://localhost:%d/v1beta1/alerts/cortex/1", s.appConfig.Service.Port), "application/json", bytes.NewBufferString(triggerAlertBody))
		s.Require().NoError(err)

		bodyJSon, _ := io.ReadAll(res.Body)
		fmt.Println(string(bodyJSon))

		_, err = io.Copy(io.Discard, res.Body)
		s.Require().NoError(err)
		defer res.Body.Close()

		<-waitChan
	})
}

func (s *CortexAlertingTestSuite) TestSendNotification() {
	ctx := context.Background()

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

	s.Run("3. Triggering alert with matching subscription labels should trigger notification", func() {
		waitChan := make(chan struct{}, 1)

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			s.Assert().NoError(err)

			expectedBody := `{"key1":"value1","key2":"value2","key3":"value3","routing_method":"receiver"}`
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
			Urn:       "subscribe-http-three",
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

		payload, err := structpb.NewStruct(map[string]interface{}{
			"data": map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		})
		s.Require().NoError(err)

		_, err = s.client.NotifyReceiver(ctx, &sirenv1beta1.NotifyReceiverRequest{
			Id:      1,
			Payload: payload,
		})
		s.Assert().NoError(err)

		<-waitChan
	})
}

func TestCortexAlertingTestSuite(t *testing.T) {
	suite.Run(t, new(CortexAlertingTestSuite))
}
