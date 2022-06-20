package v1beta1

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/odpf/salt/log"
	sirenv1beta1 "go.buf.build/odpf/gw/odpf/proton/odpf/siren/v1beta1"
)

type GRPCServer struct {
	newrelic *newrelic.Application
	logger   log.Logger
	sirenv1beta1.UnimplementedSirenServiceServer
	templateService     TemplateService
	ruleService         RuleService
	alertService        AlertService
	providerService     ProviderService
	namespaceService    NamespaceService
	receiverService     ReceiverService
	subscriptionService SubscriptionService
}

func NewGRPCServer(
	nr *newrelic.Application, logger log.Logger,
	templateService TemplateService,
	ruleService RuleService,
	alertService AlertService,
	providerService ProviderService,
	namespaceService NamespaceService,
	receiverService ReceiverService,
	subscriptionService SubscriptionService) *GRPCServer {
	return &GRPCServer{
		newrelic:            nr,
		logger:              logger,
		templateService:     templateService,
		ruleService:         ruleService,
		alertService:        alertService,
		providerService:     providerService,
		namespaceService:    namespaceService,
		receiverService:     receiverService,
		subscriptionService: subscriptionService,
	}
}

func (s *GRPCServer) Ping(ctx context.Context, in *sirenv1beta1.PingRequest) (*sirenv1beta1.PingResponse, error) {
	return &sirenv1beta1.PingResponse{Message: "Pong"}, nil
}
