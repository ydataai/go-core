package grpc

import (
	"context"

	"github.com/ydataai/go-core/pkg/common/logging"
	"google.golang.org/grpc"
)

// LoggingUnaryServerInterceptor returns a new unary server interceptors that logs the payloads of requests.
func LoggingUnaryServerInterceptor(logger logging.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.Infof("gRPC req %v -> %v", info.FullMethod, req)
		resp, err := handler(ctx, req)
		if err == nil {
			logger.Errorf("gRPC err %v -> %v", info.FullMethod, err)
		}
		logger.Infof("gRPC res %v -> %v", info.FullMethod, resp)
		return resp, err
	}
}
