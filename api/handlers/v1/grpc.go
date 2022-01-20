package v1

import (
	"context"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/odpf/salt/log"
	sirenv1beta1 "github.com/odpf/siren/api/proto/odpf/siren/v1beta1"
	"github.com/odpf/siren/service"
)

type GRPCServer struct {
	container *service.Container
	newrelic  *newrelic.Application
	logger    log.Logger
	sirenv1beta1.UnimplementedSirenServiceServer
}

func NewGRPCServer(container *service.Container, nr *newrelic.Application, logger log.Logger) *GRPCServer {
	return &GRPCServer{
		container: container,
		newrelic:  nr,
		logger:    logger,
	}
}

func (s *GRPCServer) Ping(ctx context.Context, in *sirenv1beta1.PingRequest) (*sirenv1beta1.PingResponse, error) {
	return &sirenv1beta1.PingResponse{Message: "Pong"}, nil
}
