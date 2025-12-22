package biz

import (
	"context"

	"github.com/google/wire"
{{- if and .Computed.enable_redis_final .Computed.enable_cache_final }}
	"github.com/redis/go-redis/v9"
{{- end }}
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	New{{ .Computed.service_name_capitalized }}UseCase,
)

{{- if .Computed.enable_biz_tx_final }}

type Transaction interface {
	Tx(ctx context.Context, handler func(context.Context) error) error
}
{{- end }}

{{- if and .Computed.enable_redis_final .Computed.enable_cache_final }}

type Cache interface {
	// Cache is get redis instance
	Cache() redis.UniversalClient
	// WithPrefix will add cache key prefix
	WithPrefix(prefix string) Cache
	// WithRefresh get data from db skip cache and refresh cache
	WithRefresh() Cache
	// Get is get cache data by key from redis, do write handler if cache is empty
	Get(ctx context.Context, action string, write func(context.Context) (string, error)) (string, error)
	// Set is set data to redis
	Set(ctx context.Context, action, data string, short bool)
	// Del delete key
	Del(ctx context.Context, action string)
	// SetWithExpiration is set data to redis with custom expiration
	SetWithExpiration(ctx context.Context, action, data string, seconds int64)
	// Flush is clean association cache if handler err=nil
	Flush(ctx context.Context, handler func(context.Context) error) error
	// FlushByPrefix clean cache by prefix, without prefix equals flush all by default cache prefix
	FlushByPrefix(ctx context.Context, prefix ...string) (err error)
}
{{- end }}
