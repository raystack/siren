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
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/goto/siren/config"
	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/core/template"
	"github.com/goto/siren/internal/server"
	"github.com/goto/siren/plugins"
	cortexv1plugin "github.com/goto/siren/plugins/providers/cortex/v1"
	sirenv1beta1 "github.com/goto/siren/proto/gotocompany/siren/v1beta1"
	testdatatemplate_test "github.com/goto/siren/test/e2e_test/testdata/templates"
	"github.com/mcuadros/go-defaults"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v3"
)

type NotificationTemplateTestSuite struct {
	suite.Suite
	cancelContext      context.CancelFunc
	client             sirenv1beta1.SirenServiceClient
	cancelClient       func()
	appConfig          *config.Config
	testBench          *CortexTest
	closeWorkerChannel chan struct{}
}

func (s *NotificationTemplateTestSuite) SetupTest() {

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
			Enabled: true,
		},
		DLQHandler: notification.HandlerConfig{
			Enabled: false,
		},
	}
	s.appConfig.Telemetry.OpenTelemetry.Enabled = false

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
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelContext = cancel

	StartSirenServer(ctx, *s.appConfig)

	s.closeWorkerChannel = make(chan struct{}, 1)

	time.Sleep(500 * time.Millisecond)
	StartSirenMessageWorker(ctx, *s.appConfig, s.closeWorkerChannel)

	s.client, s.cancelClient, err = CreateClient(ctx, fmt.Sprintf("localhost:%d", apiPort))
	s.Require().NoError(err)

	bootstrapCortexTestData(&s.Suite, ctx, s.client, s.testBench.NginxHost)
}

func (s *NotificationTemplateTestSuite) TearDownTest() {
	s.cancelClient()

	// Clean tests
	err := s.testBench.CleanUp()
	s.Require().NoError(err)

	s.closeWorkerChannel <- struct{}{}

	s.cancelContext()
}

func (s *NotificationTemplateTestSuite) TestSendNotificationWithTemplate() {
	ctx := context.Background()

	controlChan := make(chan struct{}, 1)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		s.Assert().NoError(err)

		type sampleStruct struct {
			Title      string `json:"title"`
			Desription string `json:"description"`
			Category   string `json:"category"`
		}

		expectedNotification := sampleStruct{
			Title: "This is the test notification with template",
			Desription: `Plain flow scalars are picky about the (:) and (#) characters. 
They can be in the string, but (:) cannot appear before a space or newline.
And (#) cannot appear after a space or newline; doing this will cause a syntax error. 
If you need to use these characters you are probably better off using one of the quoted styles instead.
`,
			Category: "httpreceiver",
		}
		var (
			resultStruct sampleStruct
		)
		s.Assert().NoError(json.Unmarshal(bodyBytes, &resultStruct))

		if diff := cmp.Diff(expectedNotification, resultStruct); diff != "" {
			s.T().Errorf("got diff: %v", diff)
		}
		controlChan <- struct{}{}

	}))
	defer testServer.Close()

	// add test server http receiver
	configs, err := structpb.NewStruct(map[string]any{
		"url": testServer.URL,
	})
	s.Require().NoError(err)
	rcv, err := s.client.CreateReceiver(ctx, &sirenv1beta1.CreateReceiverRequest{
		Name:           "notification-http-template",
		Type:           "http",
		Labels:         nil,
		Configurations: configs,
	})
	s.Require().NoError(err)

	sampleTemplateFile, err := template.YamlStringToFile(testdatatemplate_test.SampleMessageTemplate)
	s.Require().NoError(err)

	body, err := yaml.Marshal(sampleTemplateFile.Body)
	s.Require().NoError(err)

	variables := make([]*sirenv1beta1.TemplateVariables, 0)
	for _, variable := range sampleTemplateFile.Variables {
		variables = append(variables, &sirenv1beta1.TemplateVariables{
			Name:        variable.Name,
			Type:        variable.Type,
			Default:     variable.Default,
			Description: variable.Description,
		})
	}

	_, err = s.client.UpsertTemplate(ctx, &sirenv1beta1.UpsertTemplateRequest{
		Name:      sampleTemplateFile.Name,
		Body:      string(body),
		Tags:      sampleTemplateFile.Tags,
		Variables: variables,
	})
	s.Require().NoError(err)

	_, err = s.client.CreateSubscription(ctx, &sirenv1beta1.CreateSubscriptionRequest{
		Urn:       "subscribe-http",
		Namespace: 1,
		Receivers: []*sirenv1beta1.ReceiverMetadata{
			{
				Id: rcv.GetId(),
			},
		},
		Match: map[string]string{
			"category": "httpreceiver",
		},
	})
	s.Require().NoError(err)

	time.Sleep(100 * time.Millisecond)

	_, err = s.client.PostNotification(ctx, &sirenv1beta1.PostNotificationRequest{
		Receivers: []*structpb.Struct{
			{
				Fields: map[string]*structpb.Value{
					"id": structpb.NewStringValue(fmt.Sprintf("%d", rcv.GetId())),
				},
			},
		},
		Data: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"title":      structpb.NewStringValue("This is the test notification with template"),
				"icon_emoji": structpb.NewStringValue(":smile:"),
			},
		},
		Template: "test-message",
		Labels: map[string]string{
			"category": "httpreceiver",
		},
	})
	s.Require().NoError(err)

	<-controlChan

}

func TestNotificationTemplateTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationTemplateTestSuite))
}
