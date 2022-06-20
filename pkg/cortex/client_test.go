package cortex_test

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/odpf/siren/pkg/cortex"
	"github.com/odpf/siren/pkg/cortex/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateAlertmanagerConfig(t *testing.T) {
	type testCase struct {
		Description string
		Setup       func(*mocks.CortexCaller) *cortex.Client
		AMConfig    cortex.AlertManagerConfig
		Err         error
	}

	var (
		tenantID = "123123"
		amConfig = cortex.AlertManagerConfig{
			Receivers: []cortex.ReceiverConfig{
				{
					Receiver: "config1",
					Type:     "slack",
					Match: map[string]string{
						"foo": "bar"},
					Configuration: map[string]string{
						"token":        "xoxb",
						"channel_name": "test",
					},
				},
				{
					Receiver: "config2",
					Type:     "pagerduty",
					Match: map[string]string{
						"bar": "baz",
					},
					Configuration: map[string]string{
						"service_key": "1234",
					},
				},
				{
					Receiver: "config3",
					Type:     "http",
					Match:    map[string]string{},
					Configuration: map[string]string{
						"url": "http://localhost:3000",
					},
				},
			},
		}
		testCases = []testCase{
			{
				Description: "return error if error parsing config yaml",
				Setup: func(cc *mocks.CortexCaller) *cortex.Client {
					c, err := cortex.NewClient(cortex.Config{},
						cortex.WithCortexClient(cc),
						cortex.WithHelperTemplate("[[$", ""))
					require.NoError(t, err)
					return c
				},
				Err: errors.New("template: alertmanagerConfigTemplate:1: unclosed action"),
			},
			{
				Description: "return error if error loading with promconfig",
				Setup: func(cc *mocks.CortexCaller) *cortex.Client {
					c, err := cortex.NewClient(cortex.Config{}, cortex.WithCortexClient(cc))
					require.NoError(t, err)
					return c
				},
				Err: errors.New("no route provided in config"),
			},
			{
				Description: "return error if error CreateAlertmanagerConfig with cortex client",
				Setup: func(cc *mocks.CortexCaller) *cortex.Client {
					configYaml, err := ioutil.ReadFile("./testdata/config.goyaml")
					require.NoError(t, err)
					helperTemplate, err := ioutil.ReadFile("./testdata/helper.tmpl")
					require.NoError(t, err)
					expectedConfigYaml, err := ioutil.ReadFile("./testdata/expected_config.yaml")
					require.NoError(t, err)

					cc.EXPECT().CreateAlertmanagerConfig(mock.AnythingOfType("*context.valueCtx"), string(expectedConfigYaml), map[string]string{
						"helper.tmpl": string(helperTemplate),
					}).Return(errors.New("some error"))
					c, err := cortex.NewClient(cortex.Config{},
						cortex.WithCortexClient(cc),
						cortex.WithHelperTemplate(string(configYaml), string(helperTemplate)))
					require.NoError(t, err)
					return c
				},
				AMConfig: amConfig,
				Err:      errors.New("some error"),
			},
			{
				Description: "return nil error if succeed",
				Setup: func(cc *mocks.CortexCaller) *cortex.Client {
					configYaml, err := ioutil.ReadFile("./testdata/config.goyaml")
					require.NoError(t, err)
					helperTemplate, err := ioutil.ReadFile("./testdata/helper.tmpl")
					require.NoError(t, err)
					expectedConfigYaml, err := ioutil.ReadFile("./testdata/expected_config.yaml")
					require.NoError(t, err)

					cc.EXPECT().CreateAlertmanagerConfig(mock.AnythingOfType("*context.valueCtx"), string(expectedConfigYaml), map[string]string{
						"helper.tmpl": string(helperTemplate),
					}).Return(nil)
					c, err := cortex.NewClient(cortex.Config{},
						cortex.WithCortexClient(cc),
						cortex.WithHelperTemplate(string(configYaml), string(helperTemplate)))
					require.NoError(t, err)
					return c
				},
				AMConfig: amConfig,
				Err:      nil,
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			mockCortexClient := new(mocks.CortexCaller)

			c := tc.Setup(mockCortexClient)

			err := c.CreateAlertmanagerConfig(tc.AMConfig, tenantID)
			if err != tc.Err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err, tc.Err)
				}
			}
		})
	}

}
