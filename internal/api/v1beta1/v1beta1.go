package v1beta1

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/odpf/salt/log"
	"github.com/odpf/siren/internal/api"
	"github.com/odpf/siren/pkg/errors"
	sirenv1beta1 "github.com/odpf/siren/proto/odpf/siren/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	newrelic *newrelic.Application
	logger   log.Logger
	sirenv1beta1.UnimplementedSirenServiceServer
	templateService     api.TemplateService
	ruleService         api.RuleService
	alertService        api.AlertService
	providerService     api.ProviderService
	namespaceService    api.NamespaceService
	receiverService     api.ReceiverService
	subscriptionService api.SubscriptionService
}

func NewGRPCServer(
	nr *newrelic.Application,
	logger log.Logger,
	apiDeps *api.Deps) *GRPCServer {
	return &GRPCServer{
		newrelic:            nr,
		logger:              logger,
		templateService:     apiDeps.TemplateService,
		ruleService:         apiDeps.RuleService,
		alertService:        apiDeps.AlertService,
		providerService:     apiDeps.ProviderService,
		namespaceService:    apiDeps.NamespaceService,
		receiverService:     apiDeps.ReceiverService,
		subscriptionService: apiDeps.SubscriptionService,
	}
}

func (s *GRPCServer) generateRPCErr(e error) error {
	var err = errors.E(e)

	var code codes.Code
	switch {
	case errors.Is(err, errors.ErrNotFound):
		code = codes.NotFound

	case errors.Is(err, errors.ErrConflict):
		code = codes.AlreadyExists

	case errors.Is(err, errors.ErrInvalid):
		code = codes.InvalidArgument

	default:
		// TODO This will create 2 logs, grpc log and
		// the error detail (Message & Cause) log
		// there might be a better approach to solve this
		code = codes.Internal
		s.logger.Error(errors.Verbose(err).Error())
	}

	return status.Error(code, err.Error())
}
