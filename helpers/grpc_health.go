package helpers

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/health/grpc_health_v1"
	"dubbo-go-pixiu-benchmark/logger"
)

var (
	ErrServiceUnhealthy = errors.New("service is unhealthy")
	l                   = logger.GetLogger()
)

func HealthCheck(addr string, connTimeout time.Duration, rpcTimeout time.Duration) func() error {
	return func() error {
		conn, err := conn(addr, connTimeout)
		if err != nil {
			return err
		}
		defer conn.Close()
		var resp *grpc_health_v1.HealthCheckResponse
		if err := Request(context.Background(), rpcTimeout, func(rpcCtx context.Context) (err error) {
			resp, err = grpc_health_v1.NewHealthClient(conn).Check(rpcCtx,
				&grpc_health_v1.HealthCheckRequest{
					Service: ""})
			return err
		}); err != nil {
			return err
		}
		if resp.GetStatus() != grpc_health_v1.HealthCheckResponse_SERVING {
			l.Warn().Str("responded_status", resp.GetStatus().String()).Msg("service unhealthy")
			return ErrServiceUnhealthy
		}
		l.Info().Stringer("status", resp.GetStatus()).Msg("connected")
		return nil
	}
}
