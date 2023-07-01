package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/mcuadros/go-defaults"
	"github.com/raystack/salt/db"
	"github.com/raystack/siren/config"
	"github.com/raystack/siren/core/alert"
	"github.com/raystack/siren/core/log"
	"github.com/raystack/siren/core/notification"
	"github.com/raystack/siren/core/silence"
	"github.com/raystack/siren/internal/server"
	"github.com/raystack/siren/internal/store/model"
	sirenv1beta1 "github.com/raystack/siren/proto/raystack/siren/v1beta1"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/structpb"
)

type CortexAlertingTestSuite struct {
	suite.Suite
	client       sirenv1beta1.SirenServiceClient
	dbClient     *db.Client
	cancelClient func()
	appConfig    *config.Config
	testBench    *CortexTest
}

func (s *CortexAlertingTestSuite) SetupTest() {
	apiPort, err := getFreePort()
	s.Require().Nil(err)

	s.appConfig = &config.Config{}

	defaults.SetDefaults(s.appConfig)

	s.appConfig.Log.Level = "error"
	s.appConfig.Service = server.Config{
		Port: apiPort,
		GRPC: server.GRPCConfig{
			Port: 8081,
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

	// setup custom cortex config
	// TODO host.docker.internal only works for docker-desktop to call a service in host (siren)
	s.appConfig.Providers.Cortex.WebhookBaseAPI = fmt.Sprintf("http://test:%d/v1beta1/alerts/cortex", apiPort)
	s.appConfig.Providers.Cortex.GroupWaitDuration = "1s"
	s.appConfig.Providers.Cortex.GroupIntervalDuration = "1s"
	s.appConfig.Providers.Cortex.RepeatIntervalDuration = "1s"

	// enable worker
	s.appConfig.Notification.MessageHandler.Enabled = true
	s.appConfig.Notification.DLQHandler.Enabled = true

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

	s.dbClient, err = db.New(s.testBench.PGConfig)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *CortexAlertingTestSuite) TearDownTest() {
	s.cancelClient()

	// Clean tests
	err := s.testBench.CleanUp()
	s.Require().NoError(err)
}

func (s *CortexAlertingTestSuite) TestAlerting() {
	ctx := context.Background()
	triggerAlertBody := `
	[
		{
			"state": "firing",
			"value": 1,
			"labels": {
				"severity": "WARNING",
				"team": "raystack",
				"service": "some-service",
				"environment": "integration"
			},
			"annotations": {
				"resource": "test_alert",
				"metric_name": "test_alert",
				"metric_value": "1",
				"template": "alert_test"
			}
		}
	]`

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

	s.Run("Triggering cortex alert with matching subscription labels should trigger notification", func() {
		configs, err := structpb.NewStruct(map[string]interface{}{
			"url": "http://some-url",
		})
		s.Require().NoError(err)
		_, err = s.client.CreateReceiver(ctx, &sirenv1beta1.CreateReceiverRequest{
			Name: "raystack-http",
			Type: "http",
			Labels: map[string]string{
				"entity": "raystack",
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
				"team":        "raystack",
				"service":     "some-service",
				"environment": "integration",
			},
		})
		s.Require().NoError(err)

		for {
			bodyBytes, err := triggerCortexAlert(s.testBench.NginxHost, "new-raystack-1", triggerAlertBody)
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
	var (
		ctx              = context.Background()
		triggerAlertBody = `
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
						"team": "raystack"
					},
					"annotations": {
						"metric_name": "test_alert",
						"metric_value": "1",
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
				"team": "raystack"
			},
			"commonAnnotations": {
				"metric_name": "test_alert",
				"metric_value": "1",
				"resource": "test_alert",
				"template": "alert_test"
			},
			"externalURL": "/api/prom/alertmanager",
			"version": "4",
			"groupKey": "{}/{environment=\"integration\",team=\"raystack\"}:{}",
			"truncatedAlerts": 0
		}`
	)

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

	s.Run("incoming alert in alerts hook API with matching subscription labels should trigger notification", func() {
		waitChan := make(chan struct{}, 1)

		// add receiver raystack-http
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			s.Assert().NoError(err)

			type sampleStruct struct {
				ID               string `json:"id"`
				Alertname        string `json:"alertname"`
				Environment      string `json:"environment"`
				GeneratorURL     string `json:"generator_url"`
				Key1             string `json:"key1"`
				Key2             string `json:"key2"`
				MetricName       string `json:"metric_name"`
				MetricValue      string `json:"metric_value"`
				NotificationType string `json:"notification_type"`
				NumAlertsFiring  int    `json:"num_alerts_firing"`
				Resource         string `json:"resource"`
				Service          string `json:"service"`
				Severity         string `json:"severity"`
				Status           string `json:"status"`
				Firing           string `json:"firing"`
				Summary          string `json:"summary"`
				Team             string `json:"team"`
				Template         string `json:"template"`
			}

			expectedBody := `{"alertname":"some alert name","environment":"integration","generator_url":"","id":"0998ab88-3f5d-4d4a-a66f-40960b105f37","key1":"value1","key2":"value2","metric_name":"test_alert","metric_value":"1","notification_type":"subscriber","num_alerts_firing":1,"resource":"test_alert","service":"some-service","severity":"WARNING","status":"firing","summary":"this is test alert","team":"raystack","template":"alert_test"}`

			var (
				expectedStruct sampleStruct
				resultStruct   sampleStruct
			)

			s.Require().NoError(json.Unmarshal([]byte(expectedBody), &expectedStruct))
			s.Require().NoError(json.Unmarshal([]byte(body), &resultStruct))

			if diff := cmp.Diff(expectedStruct, resultStruct, cmpopts.IgnoreFields(sampleStruct{}, "ID")); diff != "" {
				s.T().Errorf("got diff: %v", diff)
			}
			close(waitChan)
		}))
		s.Require().Nil(err)
		defer testServer.Close()

		configs, err := structpb.NewStruct(map[string]interface{}{
			"url": testServer.URL,
		})
		s.Require().NoError(err)
		_, err = s.client.CreateReceiver(ctx, &sirenv1beta1.CreateReceiverRequest{
			Name: "raystack-http",
			Type: "http",
			Labels: map[string]string{
				"entity": "raystack",
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
				"team":        "raystack",
				"service":     "some-service",
				"environment": "integration",
			},
		})
		s.Require().NoError(err)

		res, err := http.DefaultClient.Post(fmt.Sprintf("http://localhost:%d/v1beta1/alerts/cortex/1/1", s.appConfig.Service.Port), "application/json", bytes.NewBufferString(triggerAlertBody))
		s.Require().NoError(err)

		bodyJSon, _ := io.ReadAll(res.Body)
		fmt.Println(string(bodyJSon))

		_, err = io.Copy(io.Discard, res.Body)
		s.Require().NoError(err)
		defer res.Body.Close()

		<-waitChan
	})

	s.Run("triggering cortex alert with matching subscription labels and silence by labels should not trigger notification", func() {
		targetExpression, err := structpb.NewStruct(map[string]interface{}{
			"team":        "raystack",
			"service":     "some-service",
			"environment": "integration",
		})
		s.Require().NoError(err)

		_, err = s.client.CreateSilence(ctx, &sirenv1beta1.CreateSilenceRequest{
			NamespaceId:      1,
			Type:             silence.TypeMatchers,
			TargetExpression: targetExpression,
		})
		s.Require().NoError(err)

		res, err := http.DefaultClient.Post(fmt.Sprintf("http://localhost:%d/v1beta1/alerts/cortex/1/1", s.appConfig.Service.Port), "application/json", bytes.NewBufferString(triggerAlertBody))
		s.Require().NoError(err)

		bodyJSon, _ := io.ReadAll(res.Body)
		fmt.Println(string(bodyJSon))

		_, err = io.Copy(io.Discard, res.Body)
		s.Require().NoError(err)
		defer res.Body.Close()

		time.Sleep(5 * time.Second)

		rows, err := s.dbClient.QueryxContext(context.Background(), `select * from notification_log`)
		s.Require().NoError(err)

		var notificationLogs []log.Notification
		for rows.Next() {
			var nlModel model.NotificationLog
			s.Require().NoError(rows.StructScan(&nlModel))
			notificationLogs = append(notificationLogs, nlModel.ToDomain())
		}

		// check alert ids of notification logs
		if diff := cmp.Diff(notificationLogs,
			[]log.Notification{
				{
					NamespaceID:    1,
					ReceiverID:     1,
					AlertIDs:       []int64{1},
					SubscriptionID: 1,
				},
				{
					NamespaceID:    1,
					SubscriptionID: 1,
					AlertIDs:       []int64{2},
				},
			},
			cmpopts.IgnoreFields(log.Notification{}, "ID", "NotificationID", "SilenceIDs", "CreatedAt")); diff != "" {
			s.T().Fatalf("found diff %v", diff)
		}

		var silenceExist bool
		for _, nl := range notificationLogs {
			if len(nl.SilenceIDs) != 0 {
				silenceExist = true
			}
		}
		s.Assert().True(silenceExist)

		rows, err = s.dbClient.QueryxContext(context.Background(), `select * from alerts`)
		s.Require().NoError(err)

		var alerts []alert.Alert
		for rows.Next() {
			var alrtModel model.Alert
			s.Require().NoError(rows.StructScan(&alrtModel))
			alerts = append(alerts, *alrtModel.ToDomain())
		}

		if diff := cmp.Diff(alerts,
			[]alert.Alert{
				{
					ID:           1,
					ProviderID:   1,
					NamespaceID:  1,
					ResourceName: "test_alert",
					MetricName:   "test_alert",
					MetricValue:  "1",
					Severity:     "WARNING",
					Rule:         "alert_test",
				},
				{
					ID:            2,
					ProviderID:    1,
					NamespaceID:   1,
					ResourceName:  "test_alert",
					MetricName:    "test_alert",
					MetricValue:   "1",
					Severity:      "WARNING",
					Rule:          "alert_test",
					SilenceStatus: alert.SilenceStatusTotal,
				},
			},
			cmpopts.IgnoreFields(alert.Alert{}, "ID", "TriggeredAt", "CreatedAt", "UpdatedAt")); diff != "" {
			s.T().Fatalf("found diff %v", diff)
		}

	})
}

func (s *CortexAlertingTestSuite) TestSendNotification() {
	ctx := context.Background()

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

	s.Run("triggering alert with matching subscription labels should trigger notification", func() {
		waitChan := make(chan struct{}, 1)

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			s.Assert().NoError(err)

			type sampleStruct struct {
				ID               string `json:"id"`
				Key1             string `json:"key1"`
				Key2             string `json:"key2"`
				Key3             string `json:"key3"`
				NotificationType string `json:"notification_type"`
				ReceiverID       string `json:"receiver_id"`
			}

			expectedBody := `{"key1":"value1","key2":"value2","key3":"value3","notification_type":"receiver","receiver_id":"1"}`
			var (
				expectedStruct sampleStruct
				resultStruct   sampleStruct
			)

			s.Require().NoError(json.Unmarshal([]byte(expectedBody), &expectedStruct))
			s.Require().NoError(json.Unmarshal([]byte(body), &resultStruct))

			if diff := cmp.Diff(expectedStruct, resultStruct, cmpopts.IgnoreFields(sampleStruct{}, "ID")); diff != "" {
				s.T().Errorf("got diff: %v", diff)
			}

			close(waitChan)
		}))
		s.Require().Nil(err)
		defer testServer.Close()

		configs, err := structpb.NewStruct(map[string]interface{}{
			"url": testServer.URL,
		})
		s.Require().NoError(err)
		_, err = s.client.CreateReceiver(ctx, &sirenv1beta1.CreateReceiverRequest{
			Name: "raystack-http",
			Type: "http",
			Labels: map[string]string{
				"entity": "raystack",
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
				"team":        "raystack",
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
