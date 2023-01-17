package grpc

import (
	"context"
	"file_processing_service/pkg/logger"
	"file_processing_service/pkg/ratelimiter"

	"google.golang.org/grpc"
)

func logUnaryInterceptor(l *logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		l.Info("request received",
			"method", info.FullMethod,
			"request", req,
		)
		res, err := handler(ctx, req)
		l.Info("response sent",
			"response", res,
		)
		return res, err
	}
}

func logStreamingInterceptor(l *logger.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		l.Info("request received", "method", info.FullMethod)
		err := handler(srv, ss)
		if err != nil {
			l.Error("error while handling request", "error", err)
		}
		return err
	}
}

func limitUnaryConnections(maxConcurrentRequestCount int) grpc.UnaryServerInterceptor {
	limiter := ratelimiter.NewLimiter(ratelimiter.WithMaxConcurrentConnections(maxConcurrentRequestCount))
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		err := limiter.Acquire()
		defer limiter.Release()
		if err != nil {
			return nil, err
		}
		res, err := handler(ctx, req)
		return res, err
	}
}

func GetUnaryInterceptors(l *logger.Logger, maxUnaryConcurrentRequestCount int) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		logUnaryInterceptor(l),
		limitUnaryConnections(maxUnaryConcurrentRequestCount),
	}
}

func GetStreamingInterceptors(l *logger.Logger) []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		logStreamingInterceptor(l),
	}
}
