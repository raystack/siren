package v1

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	pb "github.com/odpf/siren/api/proto/odpf/siren"
	"github.com/odpf/siren/service"
)

type GRPCServer struct {
	container *service.Container
	newrelic  *newrelic.Application
	pb.UnimplementedSirenServiceServer
}

func NewGRPCServer(container *service.Container, nr *newrelic.Application) *GRPCServer {
	return &GRPCServer{
		container: container,
		newrelic:  nr,
	}
}

func (s *GRPCServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "Pong"}, nil
}
