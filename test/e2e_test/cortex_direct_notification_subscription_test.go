package e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/goto/salt/db"
	"github.com/goto/siren/config"
	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/internal/server"
	"github.com/goto/siren/plugins"
	cortexv1plugin "github.com/goto/siren/plugins/providers/cortex/v1"
	sirenv1beta1 "github.com/goto/siren/proto/gotocompany/siren/v1beta1"
	"github.com/mcuadros/go-defaults"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/structpb"
)

type CortexDirectNotificationSubscriptionTestSuite struct {
	suite.Suite
	cancelContext context.CancelFunc
	grpcClient    sirenv1beta1.SirenServiceClient
	dbClient      *db.Client
	cancelClient  func()
	appConfig     *config.Config
	testBench     *CortexTest
}

func (s *CortexDirectNotificationSubscriptionTestSuite) SetupTest() {
	apiHTTPPort, err := getFreePort()
	s.Require().Nil(err)
	apiGRPCPort, err := getFreePort()
	s.Require().Nil(err)

	s.appConfig = &config.Config{}

	defaults.SetDefaults(s.appConfig)

	s.appConfig.Log.Level = "error"
	s.appConfig.Service = server.Config{
		Port: apiHTTPPort,
		GRPC: server.GRPCConfig{
			Port: apiGRPCPort,
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

	// setup custom cortex configå
	// TODO host.docker.internal only works for docker-desktop to call a service in host (siren)
	pathProject, _ := os.Getwd()
	rootProject := filepath.Dir(filepath.Dir(pathProject))
	s.appConfig.Providers.PluginPath = filepath.Join(rootProject, "plugin")
	s.appConfig.Providers.Plugins = make(map[string]plugins.PluginConfig, 0)
	s.appConfig.Providers.Plugins["cortex"] = plugins.PluginConfig{
		Handshake: plugins.HandshakeConfig{
			ProtocolVersion:  cortexv1plugin.Handshake.ProtocolVersion,
			MagicCookieKey:   cortexv1plugin.Handshake.MagicCookieKey,
			MagicCookieValue: cortexv1plugin.Handshake.MagicCookieValue,
		},
		ServiceConfig: map[string]interface{}{
			"webhook_base_api": fmt.Sprintf("http://test:%d/v1beta1/alerts/cortex", apiHTTPPort),
			"group_wait":       "1s",
			"group_interval":   "1s",
			"repeat_interval":  "1s",
		},
	}

	// enable worker
	s.appConfig.Notification.MessageHandler.Enabled = true
	s.appConfig.Notification.DLQHandler.Enabled = true

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelContext = cancel

	StartSirenServer(ctx, *s.appConfig)

	s.grpcClient, s.cancelClient, err = CreateClient(ctx, fmt.Sprintf("localhost:%d", apiGRPCPort))
	s.Require().NoError(err)

	_, err = s.grpcClient.CreateProvider(ctx, &sirenv1beta1.CreateProviderRequest{
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

func (s *CortexDirectNotificationSubscriptionTestSuite) TearDownTest() {
	s.cancelClient()

	// Clean tests
	err := s.testBench.CleanUp()
	s.Require().NoError(err)

	s.cancelContext()
}

func (s *CortexDirectNotificationSubscriptionTestSuite) TestSendNotification() {
	ctx := context.Background()

	_, err := s.grpcClient.CreateNamespace(ctx, &sirenv1beta1.CreateNamespaceRequest{
		Name:        "new-gotocompany-1",
		Urn:         "new-gotocompany-1",
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

		configs, err := structpb.NewStruct(map[string]any{
			"url": testServer.URL,
		})
		s.Require().NoError(err)
		_, err = s.grpcClient.CreateReceiver(ctx, &sirenv1beta1.CreateReceiverRequest{
			Name: "gotocompany-http",
			Type: "http",
			Labels: map[string]string{
				"entity": "gotocompany",
				"kind":   "http",
			},
			Configurations: configs,
		})
		s.Require().NoError(err)

		_, err = s.grpcClient.CreateSubscription(ctx, &sirenv1beta1.CreateSubscriptionRequest{
			Urn:       "subscribe-http-three",
			Namespace: 1,
			Receivers: []*sirenv1beta1.ReceiverMetadata{
				{
					Id: 1,
				},
			},
			Match: map[string]string{
				"team":        "gotocompany",
				"service":     "some-service",
				"environment": "integration",
			},
		})
		s.Require().NoError(err)

		payload, err := structpb.NewStruct(map[string]any{
			"data": map[string]any{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		})
		s.Require().NoError(err)

		_, err = s.grpcClient.NotifyReceiver(ctx, &sirenv1beta1.NotifyReceiverRequest{
			Id:      1,
			Payload: payload,
		})
		s.Assert().NoError(err)

		<-waitChan
	})
}

func TestCortexDirectNotificationSubscriptionTestSuite(t *testing.T) {
	suite.Run(t, new(CortexDirectNotificationSubscriptionTestSuite))
}
