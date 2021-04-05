package uploader

import (
	"errors"
	"github.com/odpf/siren/client"
	"github.com/odpf/siren/pkg/uploader/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestService_Upload(t *testing.T) {
	templateBody := `apiVersion: v2
type: template
name: test
body:
  - alert: test
    expr: avg by (host) (cpu_usage_user{cpu=\"cpu-total\"}) > 20
    for: "[[.for]]"
    labels:
      severity: WARNING
    annotations:
      description: test
variables:
  - name: for
    type: string
    default: 10m
tags:
  - systems
`
	ruleBody := `apiVersion: v2
type: rule
namespace: test
entity: real
rules:
  CPUHigh:
    - template: CPU
      status: enabled
      variables:
        - name: for
          value: 1m
        - name: team
          value: test
    - template: Lag
      status: enabled
      variables:
        - name: for
          value: 1m
        - name: team
          value: test
`
	badYaml := `abcd`

	t.Run("should call siren's cortex API with proper payload for template upsert and updates associated rules", func(t *testing.T) {
		rulesAPIMock := &mocks.RulesAPICaller{}
		templatesAPIMock := &mocks.TemplatesAPICaller{}
		mockedClient := &SirenClient{
			RulesAPI:     rulesAPIMock,
			TemplatesAPI: templatesAPIMock,
		}
		dummyService := &Service{
			SirenClient: mockedClient,
		}
		templatesAPIMock.On("CreateTemplateRequest", mock.Anything, mock.Anything).Return(client.Template{
			Name: "test",
		}, nil, nil)
		rulesAPIMock.On("ListRulesRequest", mock.Anything, mock.Anything).
			Return([]client.Rule{{
				Id:   31,
				Name: "test-rule",
				Variables: []client.RuleVariable{{
					Name: "varName", Description: "test-description", Type_: "int", Value: "30",
				},
				}}}, nil, nil)
		rulesAPIMock.On("CreateRuleRequest", mock.Anything, mock.Anything).Return(client.Rule{
			Name: "test-rule",
		}, nil, nil)
		oldFileReader := fileReader
		defer func() { fileReader = oldFileReader }()
		fileReader = func(_ string) ([]byte, error) {
			return []byte(templateBody), nil
		}
		tmp, err := dummyService.Upload("test.txt")
		temp := tmp.(*client.Template)
		assert.Equal(t, temp.Name, "test")
		assert.Nil(t, err)
	})

	t.Run("should handle error in siren's cortex API call in template upsert", func(t *testing.T) {
		rulesAPIMock := &mocks.RulesAPICaller{}
		templatesAPIMock := &mocks.TemplatesAPICaller{}
		mockedClient := &SirenClient{
			RulesAPI:     rulesAPIMock,
			TemplatesAPI: templatesAPIMock,
		}
		dummyService := &Service{
			SirenClient: mockedClient,
		}
		templatesAPIMock.On("CreateTemplateRequest", mock.Anything, mock.Anything).Return(client.Template{}, nil, errors.New("random error"))
		oldFileReader := fileReader
		defer func() { fileReader = oldFileReader }()
		fileReader = func(_ string) ([]byte, error) {
			return []byte(templateBody), nil
		}
		tmp, err := dummyService.Upload("test.txt")
		assert.EqualError(t, err, "random error")
		assert.Nil(t, tmp)
	})

	t.Run("should handle errors in getting associated rules after template upsert", func(t *testing.T) {
		rulesAPIMock := &mocks.RulesAPICaller{}
		templatesAPIMock := &mocks.TemplatesAPICaller{}
		mockedClient := &SirenClient{
			RulesAPI:     rulesAPIMock,
			TemplatesAPI: templatesAPIMock,
		}
		dummyService := &Service{
			SirenClient: mockedClient,
		}
		templatesAPIMock.On("CreateTemplateRequest", mock.Anything, mock.Anything).Return(client.Template{
			Name: "test",
		}, nil, nil)
		rulesAPIMock.On("ListRulesRequest", mock.Anything, mock.Anything).
			Return([]client.Rule{{
				Id:   31,
				Name: "test-rule",
				Variables: []client.RuleVariable{{
					Name: "varName", Description: "test-description", Type_: "int", Value: "30",
				},
				}}}, nil, errors.New("random error"))
		oldFileReader := fileReader
		defer func() { fileReader = oldFileReader }()
		fileReader = func(_ string) ([]byte, error) {
			return []byte(templateBody), nil
		}
		tmp, err := dummyService.Upload("test.txt")
		temp := tmp.(*client.Template)
		assert.Equal(t, temp.Name, "test")
		assert.EqualError(t, err, "random error")
	})

	t.Run("should handle errors in updating associated rules after template upsert", func(t *testing.T) {
		rulesAPIMock := &mocks.RulesAPICaller{}
		templatesAPIMock := &mocks.TemplatesAPICaller{}
		mockedClient := &SirenClient{
			RulesAPI:     rulesAPIMock,
			TemplatesAPI: templatesAPIMock,
		}
		dummyService := &Service{
			SirenClient: mockedClient,
		}
		templatesAPIMock.On("CreateTemplateRequest", mock.Anything, mock.Anything).Return(client.Template{
			Name: "test",
		}, nil, nil)
		rulesAPIMock.On("ListRulesRequest", mock.Anything, mock.Anything).
			Return([]client.Rule{{
				Id:   31,
				Name: "test-rule",
				Variables: []client.RuleVariable{{
					Name: "varName", Description: "test-description", Type_: "int", Value: "30",
				},
				}}}, nil, nil)
		rulesAPIMock.On("CreateRuleRequest", mock.Anything, mock.Anything).Return(client.Rule{}, nil, errors.New("random error"))
		oldFileReader := fileReader
		defer func() { fileReader = oldFileReader }()
		fileReader = func(_ string) ([]byte, error) {
			return []byte(templateBody), nil
		}
		tmp, err := dummyService.Upload("test.txt")
		temp := tmp.(*client.Template)
		assert.Equal(t, temp.Name, "test")
		assert.EqualError(t, err, "random error")
	})

	t.Run("should call siren's cortex API with proper payload for rule upsert", func(t *testing.T) {
		rulesAPIMock := &mocks.RulesAPICaller{}
		templatesAPIMock := &mocks.TemplatesAPICaller{}
		mockedClient := &SirenClient{
			RulesAPI:     rulesAPIMock,
			TemplatesAPI: templatesAPIMock,
		}
		dummyService := &Service{
			SirenClient: mockedClient,
		}
		rulesAPIMock.On("CreateRuleRequest", mock.Anything, mock.Anything).Return(client.Rule{
			Name: "foo",
		}, nil, nil)
		oldFileReader := fileReader
		defer func() { fileReader = oldFileReader }()
		fileReader = func(_ string) ([]byte, error) {
			return []byte(ruleBody), nil
		}
		result, err := dummyService.Upload("test.txt")
		rules := result.([]*client.Rule)
		assert.Equal(t, 2, len(rules))
		assert.Equal(t, "foo", rules[0].Name)
		rulesAPIMock.AssertNumberOfCalls(t, "CreateRuleRequest", 2)
		assert.Nil(t, err)
	})

	t.Run("should handle error in siren's cortex API call in template rule upsert", func(t *testing.T) {
		rulesAPIMock := &mocks.RulesAPICaller{}
		templatesAPIMock := &mocks.TemplatesAPICaller{}
		mockedClient := &SirenClient{
			RulesAPI:     rulesAPIMock,
			TemplatesAPI: templatesAPIMock,
		}
		dummyService := &Service{
			SirenClient: mockedClient,
		}
		rulesAPIMock.On("CreateRuleRequest", mock.Anything, mock.Anything).
			Return(client.Rule{}, nil, errors.New("random error"))
		oldFileReader := fileReader
		defer func() { fileReader = oldFileReader }()
		fileReader = func(_ string) ([]byte, error) {
			return []byte(ruleBody), nil
		}
		result, err := dummyService.Upload("test.txt")
		rules := result.([]*client.Rule)
		assert.Equal(t, 0, len(rules))
		assert.EqualError(t, err, "random error")
	})

	t.Run("should handle file read errors", func(t *testing.T) {
		rulesAPIMock := &mocks.RulesAPICaller{}
		templatesAPIMock := &mocks.TemplatesAPICaller{}
		mockedClient := &SirenClient{
			RulesAPI:     rulesAPIMock,
			TemplatesAPI: templatesAPIMock,
		}
		dummyService := &Service{
			SirenClient: mockedClient,
		}
		oldFileReader := fileReader
		defer func() { fileReader = oldFileReader }()
		fileReader = func(_ string) ([]byte, error) {
			return nil, errors.New("random error")
		}
		tmp, err := dummyService.Upload("test.txt")
		assert.EqualError(t, err, "random error")
		assert.Nil(t, tmp)
	})

	t.Run("should handle errors with bad yaml files", func(t *testing.T) {
		rulesAPIMock := &mocks.RulesAPICaller{}
		templatesAPIMock := &mocks.TemplatesAPICaller{}
		mockedClient := &SirenClient{
			RulesAPI:     rulesAPIMock,
			TemplatesAPI: templatesAPIMock,
		}
		dummyService := &Service{
			SirenClient: mockedClient,
		}
		oldFileReader := fileReader
		defer func() { fileReader = oldFileReader }()
		fileReader = func(_ string) ([]byte, error) {
			return []byte(badYaml), nil
		}
		tmp, err := dummyService.Upload("test.txt")
		assert.EqualError(t, err, "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `abcd` into uploader.yamlObject")
		assert.Nil(t, tmp)
	})

	t.Run("should handle unknown types", func(t *testing.T) {
		mockedClient := &SirenClient{
			RulesAPI: &mocks.RulesAPICaller{},
			TemplatesAPI: &mocks.TemplatesAPICaller{
			}}
		dummyService := &Service{
			SirenClient: mockedClient,
		}
		oldFileReader := fileReader
		defer func() { fileReader = oldFileReader }()
		fileReader = func(_ string) ([]byte, error) {
			body := `type: abcd`
			return []byte(body), nil
		}
		res, err := dummyService.Upload("test.txt")
		assert.Nil(t, res)
		assert.EqualError(t, err, "unknown type given")
	})
}
