package v1beta1

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/odpf/salt/log"

	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	sirenv1beta1 "go.buf.build/odpf/gw/odpf/proton/odpf/siren/v1beta1"
)

func TestGRPCServer_ListRules(t *testing.T) {
	dummyPayload := &sirenv1beta1.ListRulesRequest{
		Name:              "foo",
		Namespace:         "test",
		GroupName:         "foo",
		Template:          "foo",
		ProviderNamespace: 1,
	}

	t.Run("should return stored rules", func(t *testing.T) {
		ctx := context.Background()
		mockedRuleService := &mocks.RuleService{}
		dummyResult := []domain.Rule{
			{
				Name:      "foo",
				Enabled:   true,
				GroupName: "foo",
				Namespace: "test",
				Template:  "foo",
				Variables: []domain.RuleVariable{
					{
						Name:        "foo",
						Type:        "int",
						Value:       "bar",
						Description: "",
					},
				},
				ProviderNamespace: 1,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
		}

		dummyGRPCServer := GRPCServer{
			container: &Container{
				RulesService: mockedRuleService,
			},
			logger: log.NewNoop(),
		}
		mockedRuleService.
			On("Get", ctx, dummyPayload.Name, dummyPayload.Namespace, dummyPayload.GroupName,
				dummyPayload.Template, dummyPayload.ProviderNamespace).
			Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.ListRules(ctx, dummyPayload)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res.GetRules()))
		assert.Equal(t, "foo", res.GetRules()[0].GetName())
		assert.Equal(t, "test", res.GetRules()[0].GetNamespace())
		assert.Equal(t, true, res.GetRules()[0].GetEnabled())
		assert.Equal(t, 1, len(res.GetRules()[0].GetVariables()))
		mockedRuleService.AssertExpectations(t)
	})

	t.Run("should return error code 13 if getting rules failed", func(t *testing.T) {
		ctx := context.Background()
		mockedRuleService := &mocks.RuleService{}

		dummyGRPCServer := GRPCServer{
			container: &Container{
				RulesService: mockedRuleService,
			},
			logger: log.NewNoop(),
		}
		mockedRuleService.
			On("Get", ctx, dummyPayload.Name, dummyPayload.Namespace, dummyPayload.GroupName,
				dummyPayload.Template, dummyPayload.ProviderNamespace).
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.ListRules(ctx, dummyPayload)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})
}

func TestGRPCServer_UpdateRules(t *testing.T) {
	dummyPayload := &domain.Rule{
		Enabled:   true,
		GroupName: "foo",
		Namespace: "test",
		Template:  "foo",
		Variables: []domain.RuleVariable{
			{
				Name:        "foo",
				Type:        "int",
				Value:       "bar",
				Description: "",
			},
		},
		ProviderNamespace: 1,
	}
	dummyReq := &sirenv1beta1.UpdateRuleRequest{
		Enabled:   true,
		GroupName: "foo",
		Namespace: "test",
		Template:  "foo",
		Variables: []*sirenv1beta1.Variables{
			{
				Name:        "foo",
				Type:        "int",
				Value:       "bar",
				Description: "",
			},
		},
		ProviderNamespace: 1,
	}

	t.Run("should update rule", func(t *testing.T) {
		ctx := context.Background()
		mockedRuleService := &mocks.RuleService{}
		dummyGRPCServer := GRPCServer{
			container: &Container{
				RulesService: mockedRuleService,
			},
			logger: log.NewNoop(),
		}
		dummyResult := &domain.Rule{}
		*dummyResult = *dummyPayload
		dummyResult.Enabled = false
		dummyResult.Name = "foo"

		mockedRuleService.
			On("Upsert", ctx, dummyPayload).
			Run(func(args mock.Arguments) {
				r := args.Get(1).(*domain.Rule)
				*r = *dummyResult
			}).
			Return(nil).Once()
		res, err := dummyGRPCServer.UpdateRule(ctx, dummyReq)
		assert.Nil(t, err)

		assert.Equal(t, "foo", res.GetRule().GetName())
		assert.Equal(t, false, res.GetRule().GetEnabled())
		assert.Equal(t, "test", res.GetRule().GetNamespace())
		assert.Equal(t, 1, len(res.GetRule().GetVariables()))
		mockedRuleService.AssertExpectations(t)
	})

	t.Run("should return error code 13 if getting rules failed", func(t *testing.T) {
		ctx := context.Background()
		mockedRuleService := &mocks.RuleService{}

		dummyGRPCServer := GRPCServer{
			container: &Container{
				RulesService: mockedRuleService,
			},
			logger: log.NewNoop(),
		}
		mockedRuleService.
			On("Upsert", ctx, dummyPayload).
			Return(errors.New("random error")).Once()
		res, err := dummyGRPCServer.UpdateRule(ctx, dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})
}
