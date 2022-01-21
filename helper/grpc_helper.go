package helper

import (
	"github.com/odpf/salt/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCLogError(log log.Logger, codes codes.Code, err error) error {
	log.Error("failed to handle alert", "error", err)
	return status.Errorf(codes, err.Error())
}
