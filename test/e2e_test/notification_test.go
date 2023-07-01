package e2e_test

import (
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
	"github.com/raystack/siren/config"
	"github.com/raystack/siren/core/notification"
	"github.com/raystack/siren/internal/server"
	sirenv1beta1 "github.com/raystack/siren/proto/raystack/siren/v1beta1"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/structpb"
)

type NotificationTestSuite struct {
	suite.Suite
	// testServer   *httptest.Server
	client             sirenv1beta1.SirenServiceClient
	cancelClient       func()
	appConfig          *config.Config
	testBench          *CortexTest
	closeWorkerChannel chan struct{}
}

func (s *NotificationTestSuite) SetupTest() {

	apiPort, err := getFreePort()
	s.Require().Nil(err)

	grpcPort, err := getFreePort()
	s.Require().Nil(err)

	s.appConfig = &config.Config{}

	defaults.SetDefaults(s.appConfig)

	s.appConfig.Log.Level = "error"
	s.appConfig.Service = server.Config{
		GRPC: server.GRPCConfig{
			Port: grpcPort,
		},
		Port:          apiPort,
		EncryptionKey: testEncryptionKey,
	}
	s.appConfig.Notification = notification.Config{
		MessageHandler: notification.HandlerConfig{
			Enabled: true,
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
	s.appConfig.Providers.Cortex.WebhookBaseAPI = fmt.Sprintf("http://host.docker.internal:%d/v1beta1/alerts/cortex", apiPort)
	s.appConfig.Providers.Cortex.GroupWaitDuration = "1s"

	StartSirenServer(*s.appConfig)

	s.closeWorkerChannel = make(chan struct{}, 1)

	time.Sleep(500 * time.Millisecond)
	StartSirenMessageWorker(*s.appConfig, s.closeWorkerChannel)

	ctx := context.Background()
	s.client, s.cancelClient, err = CreateClient(ctx, fmt.Sprintf("localhost:%d", grpcPort))
	s.Require().NoError(err)

	bootstrapCortexTestData(&s.Suite, ctx, s.client, s.testBench.NginxHost)
}

func (s *NotificationTestSuite) TearDownTest() {
	s.cancelClient()
	// Clean tests
	err := s.testBench.CleanUp()
	s.Require().NoError(err)

	s.closeWorkerChannel <- struct{}{}
}

func (s *NotificationTestSuite) TestSendNotification() {
	ctx := context.Background()

	controlChan := make(chan struct{}, 1)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		s.Assert().NoError(err)

		type sampleStruct struct {
			ID               string `json:"id"`
			IconEmoji        string `json:"icon_emoji"`
			NotificationType string `json:"notification_type"`
			ReceiverID       string `json:"receiver_id"`
			Text             string `json:"text"`
		}

		expectedNotification := `{"icon_emoji":":smile:","notification_type":"receiver","receiver_id":"2","text":"test send notification"}`

		var (
			resultStruct   sampleStruct
			expectedStruct sampleStruct
		)
		s.Assert().NoError(json.Unmarshal(bodyBytes, &resultStruct))
		s.Assert().NoError(json.Unmarshal([]byte(expectedNotification), &expectedStruct))

		if diff := cmp.Diff(expectedStruct, resultStruct, cmpopts.IgnoreFields(sampleStruct{}, "ID")); diff != "" {
			s.T().Errorf("got diff: %v", diff)
		}
		controlChan <- struct{}{}

	}))
	defer testServer.Close()

	// add test server http receiver
	configs, err := structpb.NewStruct(map[string]interface{}{
		"url": testServer.URL,
	})
	s.Require().NoError(err)
	rcv, err := s.client.CreateReceiver(ctx, &sirenv1beta1.CreateReceiverRequest{
		Name:           "notification-http",
		Type:           "http",
		Labels:         nil,
		Configurations: configs,
	})
	s.Require().NoError(err)

	time.Sleep(100 * time.Millisecond)

	_, err = s.client.NotifyReceiver(ctx, &sirenv1beta1.NotifyReceiverRequest{
		Id: rcv.GetId(),
		Payload: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"data": structpb.NewStructValue(&structpb.Struct{
					Fields: map[string]*structpb.Value{
						"text":       structpb.NewStringValue("test send notification"),
						"icon_emoji": structpb.NewStringValue(":smile:"),
					},
				}),
			},
		},
	})
	s.Require().NoError(err)

	<-controlChan

}

func TestNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationTestSuite))
}
