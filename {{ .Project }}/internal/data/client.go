package data

import (
	"context"
	"strings"
	"time"

	"github.com/go-cinch/common/log"
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
{{- if .Computed.enable_trace_final }}
	"github.com/go-kratos/kratos/v2/middleware/tracing"
{{- end }}
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/pkg/errors"
	g "google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// NewClient creates a new gRPC client with middleware and health check support.
// T is the client type (e.g., auth.AuthClient)
// name is the service name for logging
// endpoint is the gRPC server address (e.g., "localhost:9000")
// health enables health check before returning the client
// timeout is the client timeout duration
// newClient is the constructor function for the specific client type
func NewClient[T any](name, endpoint string, health bool, timeout time.Duration, newClient func(cc g.ClientConnInterface) T) (client T, err error) {
	ops := []grpc.ClientOption{
		grpc.WithEndpoint(endpoint),
		grpc.WithMiddleware(
{{- if .Computed.enable_trace_final }}
			tracing.Client(),
{{- end }}
			metadata.Client(),
			circuitbreaker.Client(),
			recovery.Recovery(),
		),
		grpc.WithOptions(g.WithDisableHealthCheck()),
		grpc.WithTimeout(timeout),
	}
	conn, err := grpc.DialInsecure(context.Background(), ops...)
	if err != nil {
		err = errors.WithMessage(err, strings.Join([]string{"initialize", name, "client failed"}, " "))
		return
	}
	if health {
		healthClient := healthpb.NewHealthClient(conn)
		_, err = healthClient.Check(context.Background(), &healthpb.HealthCheckRequest{})
		if err != nil {
			err = errors.WithMessage(err, strings.Join([]string{name, "healthcheck failed"}, " "))
			return
		}
	}
	client = newClient(conn)
	log.
		WithField("endpoint", endpoint).
		Info(strings.Join([]string{"initialize", name, "client success"}, " "))
	return
}

// Example usage:
// To create a client for another service, add its proto definition and use:
//
// import "{{ .Computed.module_name_final }}/internal/conf"
//
// func NewAuthClient(c *conf.Bootstrap) (auth.AuthClient, error) {
// 	return NewClient[auth.AuthClient](
// 		"auth",
// 		c.Client.Auth,
// 		c.Client.Health,
// 		c.Client.Timeout.AsDuration(),
// 		auth.NewAuthClient,
// 	)
// }
