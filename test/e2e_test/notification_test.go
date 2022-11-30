package e2e_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mcuadros/go-defaults"
	"github.com/odpf/siren/config"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/internal/server"
	"github.com/odpf/siren/plugins/queues"
	sirenv1beta1 "github.com/odpf/siren/proto/odpf/siren/v1beta1"
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

	s.appConfig = &config.Config{
		Log: config.Log{
			Level: "debug",
		},
		Service: server.Config{
			Port:          apiPort,
			EncryptionKey: testEncryptionKey,
		},
		Notification: notification.Config{
			Queue: queues.Config{
				Kind: queues.KindPostgres,
			},
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

	s.closeWorkerChannel = make(chan struct{}, 1)

	time.Sleep(500 * time.Millisecond)
	StartSirenMessageWorker(*s.appConfig, s.closeWorkerChannel)

	ctx := context.Background()
	s.client, s.cancelClient, err = CreateClient(ctx, fmt.Sprintf("localhost:%d", apiPort))
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
	expectedNotification := `{"icon_emoji":":smile:","routing_method":"receiver","text":"test send notification"}`
	payload := `
	{
		"payload": {
			"data": {
				"text": "test send notification",
				"icon_emoji": ":smile:"
			}	
		}
	}`
	ctx := context.Background()

	controlChan := make(chan struct{}, 1)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		s.Assert().NoError(err)
		s.Assert().Equal(expectedNotification, string(bodyBytes))

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

	notifyAPI := fmt.Sprintf("http://localhost:%d/v1beta1/receivers/%d/send", s.appConfig.Service.Port, rcv.GetId())

	bodyBytes := []byte(payload)
	_, err = http.Post(notifyAPI, "application/json", bytes.NewReader(bodyBytes))
	s.Require().NoError(err)

	<-controlChan

}

func TestNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationTestSuite))
}
