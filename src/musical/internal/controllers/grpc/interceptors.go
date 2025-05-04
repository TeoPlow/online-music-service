package grpc

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/TeoPlow/online-music-service/src/musical/internal/logger"
)

func logInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	logger.Logger.Info("gRPC_Request",
		slog.String("Method", info.FullMethod))

	resp, err := handler(ctx, req)

	code := status.Code(err)
	if code == codes.Internal {
		logger.Logger.Error("gRPC_Response",
			slog.String("Error", err.Error()))
	} else {
		logger.Logger.Info("gRPC_Response",
			slog.String("Code", code.String()))
	}

	return resp, err
}
