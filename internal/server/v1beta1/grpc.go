package v1beta1

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/odpf/salt/log"
	sirenv1beta1 "go.buf.build/odpf/gw/odpf/proton/odpf/siren/v1beta1"
)

type GRPCServer struct {
	container *Container
	newrelic  *newrelic.Application
	logger    log.Logger
	sirenv1beta1.UnimplementedSirenServiceServer
}

func NewGRPCServer(container *Container, nr *newrelic.Application, logger log.Logger) *GRPCServer {
	return &GRPCServer{
		container: container,
		newrelic:  nr,
		logger:    logger,
	}
}

func (s *GRPCServer) Ping(ctx context.Context, in *sirenv1beta1.PingRequest) (*sirenv1beta1.PingResponse, error) {
	return &sirenv1beta1.PingResponse{Message: "Pong"}, nil
}
