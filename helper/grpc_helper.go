package helper

import (
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCLogError(log *zap.Logger, codes codes.Code, err error) error {
	log.Error("handler", zap.Error(err))
	return status.Errorf(codes, err.Error())
}
