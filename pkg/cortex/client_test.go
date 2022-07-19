package cortex_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"
	"github.com/odpf/siren/pkg/cortex"
	"github.com/odpf/siren/pkg/cortex/mocks"
	"github.com/odpf/siren/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClient_New(t *testing.T) {
	t.Run("should initiate cortex client if not passed from option", func(t *testing.T) {
		c, err := cortex.NewClient(cortex.Config{})
		if err != nil {
			t.Fatalf("got error %v, expected was nil", err)
		}
		if c == nil {
			t.Fatalf("got client %v, expected was not nil", c)
		}
	})

	t.Run("should return error when cortex client client creation return error", func(t *testing.T) {
		c, err := cortex.NewClient(cortex.Config{
			Address: ":::",
		})
		expectedErrorString := "parse \":::\": missing protocol scheme"
		if err.Error() != expectedErrorString {
			t.Fatalf("got error %v, expected was %v", err, expectedErrorString)
		}
		if c != nil {
			t.Fatalf("got client %v, expected was nil", c)
		}
	})
}

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

func TestClient_CreateRuleGroup(t *testing.T) {

	t.Run("should return error if cortex client return error", func(t *testing.T) {
		cortexCallerMock := &mocks.CortexCaller{}

		c, err := cortex.NewClient(cortex.Config{}, cortex.WithCortexClient(cortexCallerMock))
		require.Nil(t, err)
		require.NotNil(t, c)

		cortexCallerMock.EXPECT().CreateRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("rwrulefmt.RuleGroup")).Return(errors.New("some error"))

		err = c.CreateRuleGroup(context.TODO(), "namespace", rwrulefmt.RuleGroup{})
		assert.NotNil(t, err)

		cortexCallerMock.AssertExpectations(t)
	})

	t.Run("should return nil error if cortex client return nil error", func(t *testing.T) {
		cortexCallerMock := &mocks.CortexCaller{}

		c, err := cortex.NewClient(cortex.Config{}, cortex.WithCortexClient(cortexCallerMock))
		require.Nil(t, err)
		require.NotNil(t, c)

		cortexCallerMock.EXPECT().CreateRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("rwrulefmt.RuleGroup")).Return(nil)

		err = c.CreateRuleGroup(context.TODO(), "namespace", rwrulefmt.RuleGroup{})
		assert.Nil(t, err)

		cortexCallerMock.AssertExpectations(t)
	})
}

func TestClient_DeleteRuleGroup(t *testing.T) {

	t.Run("should return error if cortex client return error", func(t *testing.T) {
		cortexCallerMock := &mocks.CortexCaller{}

		c, err := cortex.NewClient(cortex.Config{}, cortex.WithCortexClient(cortexCallerMock))
		require.Nil(t, err)
		require.NotNil(t, c)

		cortexCallerMock.EXPECT().DeleteRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New("some error"))

		err = c.DeleteRuleGroup(context.TODO(), "namespace", "groupname")
		assert.NotNil(t, err)

		cortexCallerMock.AssertExpectations(t)
	})

	t.Run("should return nil error if cortex client return nil error", func(t *testing.T) {
		cortexCallerMock := &mocks.CortexCaller{}

		c, err := cortex.NewClient(cortex.Config{}, cortex.WithCortexClient(cortexCallerMock))
		require.Nil(t, err)
		require.NotNil(t, c)

		cortexCallerMock.EXPECT().DeleteRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

		err = c.DeleteRuleGroup(context.TODO(), "namespace", "groupname")
		assert.Nil(t, err)

		cortexCallerMock.AssertExpectations(t)
	})
}

func TestClient_GetRuleGroup(t *testing.T) {

	t.Run("should return error if cortex client return error", func(t *testing.T) {
		cortexCallerMock := &mocks.CortexCaller{}

		c, err := cortex.NewClient(cortex.Config{}, cortex.WithCortexClient(cortexCallerMock))
		require.Nil(t, err)
		require.NotNil(t, c)

		cortexCallerMock.EXPECT().GetRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("some error"))

		rg, err := c.GetRuleGroup(context.TODO(), "namespace", "groupname")
		assert.NotNil(t, err)
		assert.Nil(t, rg)

		cortexCallerMock.AssertExpectations(t)
	})

	t.Run("should return nil error if cortex client return nil error", func(t *testing.T) {
		cortexCallerMock := &mocks.CortexCaller{}

		c, err := cortex.NewClient(cortex.Config{}, cortex.WithCortexClient(cortexCallerMock))
		require.Nil(t, err)
		require.NotNil(t, c)

		cortexCallerMock.EXPECT().GetRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&rwrulefmt.RuleGroup{}, nil)

		rg, err := c.GetRuleGroup(context.TODO(), "namespace", "groupname")
		assert.Nil(t, err)
		assert.NotNil(t, rg)

		cortexCallerMock.AssertExpectations(t)
	})
}

func TestClient_ListRules(t *testing.T) {

	t.Run("should return error if cortex client return error", func(t *testing.T) {
		cortexCallerMock := &mocks.CortexCaller{}

		c, err := cortex.NewClient(cortex.Config{}, cortex.WithCortexClient(cortexCallerMock))
		require.Nil(t, err)
		require.NotNil(t, c)

		cortexCallerMock.EXPECT().ListRules(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(nil, errors.New("some error"))

		rg, err := c.ListRules(context.TODO(), "namespace")
		assert.NotNil(t, err)
		assert.Nil(t, rg)

		cortexCallerMock.AssertExpectations(t)
	})

	t.Run("should return nil error if cortex client return nil error", func(t *testing.T) {
		cortexCallerMock := &mocks.CortexCaller{}

		c, err := cortex.NewClient(cortex.Config{}, cortex.WithCortexClient(cortexCallerMock))
		require.Nil(t, err)
		require.NotNil(t, c)

		cortexCallerMock.EXPECT().ListRules(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(map[string][]rwrulefmt.RuleGroup{}, nil)

		rg, err := c.ListRules(context.TODO(), "namespace")
		assert.Nil(t, err)
		assert.NotNil(t, rg)

		cortexCallerMock.AssertExpectations(t)
	})
}
