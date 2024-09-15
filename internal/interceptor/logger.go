package interceptor

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/mikhailsoldatkin/auth/internal/logger"
)

// LoggingInterceptor is a gRPC unary server interceptor for logging requests.
func LoggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	now := time.Now()

	res, err := handler(ctx, req)

	if err != nil {
		logger.Error(
			err.Error(),
			zap.String("method", info.FullMethod),
			zap.Any("request", req),
		)
	} else {
		logger.Info(
			"success",
			zap.String("method", info.FullMethod),
			zap.Any("request", req),
			zap.Any("result", res),
			zap.Duration("duration", time.Since(now)),
		)
	}

	return res, err
}
