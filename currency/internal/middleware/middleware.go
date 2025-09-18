package middleware

import (
	"context"
	"github.com/BernsteinMondy/currency-service/currency/internal/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

type Middleware struct {
	metrics *metrics.Metrics
	logger  *zap.Logger
}

func Init(metrics *metrics.Metrics, logger *zap.Logger) *Middleware {
	return &Middleware{metrics: metrics, logger: logger}
}

func (m *Middleware) LoggingMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	m.logger.Info("Received request", zap.String("method", info.FullMethod))
	resp, err := handler(ctx, req)
	if err != nil {
		m.logger.Error("Failed to handle request", zap.String("method", info.FullMethod), zap.Error(err))
	}

	return resp, err
}

func (m *Middleware) RequestDurationMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	m.metrics.Counter.WithLabelValues(info.FullMethod).Inc()
	start := time.Now()

	resp, err := handler(ctx, req)

	m.metrics.Latency.WithLabelValues(info.FullMethod).Observe(time.Since(start).Seconds())
	return resp, err
}
