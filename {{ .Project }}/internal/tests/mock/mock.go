package mock

import (
	"context"
	"fmt"
{{- if or (eq .Computed.db_type_final "mysql") .Computed.enable_redis_final }}
	"net"
{{- end }}
	"net/http"
	"runtime/debug"
	"sync"

{{- if or (eq .Computed.db_type_final "mysql") .Computed.enable_redis_final }}
	"github.com/go-cinch/common/mock"
{{- end }}
	"github.com/go-kratos/kratos/v2/transport"

	"{{ .Computed.module_name_final }}/internal/biz"
	"{{ .Computed.module_name_final }}/internal/conf"
	"{{ .Computed.module_name_final }}/internal/data"
	"{{ .Computed.module_name_final }}/internal/service"
)

{{- if eq .Scaffold.proto_template "full" }}

func {{ .Computed.service_name_capitalized }}Service() (svc *service.{{ .Computed.service_name_capitalized }}Service) {
	useCase := {{ .Computed.service_name_capitalized }}UseCase()
	svc = service.New{{ .Computed.service_name_capitalized }}Service(useCase)
	return
}

func {{ .Computed.service_name_capitalized }}UseCase() (useCase *biz.{{ .Computed.service_name_capitalized }}UseCase) {
{{- if .Computed.enable_biz_tx_final }}
	{{- if .Computed.enable_cache_final }}
	c, dataData, cache := Data()
	repo := {{ .Computed.service_name_capitalized }}Repo()
	transaction := data.NewTransaction(dataData)
	useCase = biz.New{{ .Computed.service_name_capitalized }}UseCase(c, repo, transaction, cache)
	{{- else }}
	c, dataData, _ := Data()
	repo := {{ .Computed.service_name_capitalized }}Repo()
	transaction := data.NewTransaction(dataData)
	useCase = biz.New{{ .Computed.service_name_capitalized }}UseCase(c, repo, transaction)
	{{- end }}
{{- else }}
	{{- if .Computed.enable_cache_final }}
	c, _, cache := Data()
	repo := {{ .Computed.service_name_capitalized }}Repo()
	useCase = biz.New{{ .Computed.service_name_capitalized }}UseCase(c, repo, cache)
	{{- else }}
	c, _, _ := Data()
	repo := {{ .Computed.service_name_capitalized }}Repo()
	useCase = biz.New{{ .Computed.service_name_capitalized }}UseCase(c, repo)
	{{- end }}
{{- end }}
	return
}

func {{ .Computed.service_name_capitalized }}Repo() biz.{{ .Computed.service_name_capitalized }}Repo {
	_, d, _ := Data()
	return data.New{{ .Computed.service_name_capitalized }}Repo(d)
}
{{- end }}

type headerCarrier http.Header

func (hc headerCarrier) Get(key string) string { return http.Header(hc).Get(key) }

func (hc headerCarrier) Set(key string, value string) { http.Header(hc).Set(key, value) }

func (hc headerCarrier) Add(key string, value string) { http.Header(hc).Add(key, value) }

// Keys lists the keys stored in this carrier.
func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice value associated with the passed key.
func (hc headerCarrier) Values(key string) []string {
	return http.Header(hc).Values(key)
}

func newUserHeader(k, v string) *headerCarrier {
	header := &headerCarrier{}
	header.Set(k, v)
	return header
}

type Transport struct {
	kind      transport.Kind
	endpoint  string
	operation string
	reqHeader transport.Header
}

func (tr *Transport) Kind() transport.Kind {
	return tr.kind
}

func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

func (tr *Transport) Operation() string {
	return tr.operation
}

func (tr *Transport) RequestHeader() transport.Header {
	return tr.reqHeader
}

func (*Transport) ReplyHeader() transport.Header {
	return nil
}

func NewContextWithUserId(ctx context.Context, u string) context.Context {
	tr := &Transport{
		reqHeader: newUserHeader("X-Md-Global-Code", u),
	}
	return transport.NewServerContext(ctx, tr)
}

var (
	onceC     *conf.Bootstrap
	onceData  *data.Data
{{- if .Computed.enable_cache_final }}
	onceCache biz.Cache
{{- else }}
	onceCache any
{{- end }}
	once      sync.Once
)

{{- if .Computed.enable_db_final }}

func Data() (c *conf.Bootstrap, dataData *data.Data, cache {{- if .Computed.enable_cache_final }} biz.Cache{{- else }} any{{- end }}) {
	debug.SetGCPercent(-1)
	once.Do(func() {
{{- if eq .Computed.db_type_final "mysql" }}
		onceC = MySQLAndRedis()
{{- else }}
		onceC = PostgreSQLAndRedis()
{{- end }}
		onceData, _, _ = data.NewData(onceC)
{{- if and .Computed.enable_redis_final .Computed.enable_cache_final }}
		universalClient, err := data.NewRedis(onceC)
		if err != nil {
			panic(err)
		}
		onceCache = data.NewCache(onceC, universalClient)
{{- else }}
		onceCache = nil
{{- end }}
	})
	return onceC, onceData, onceCache
}

{{- if eq .Computed.db_type_final "mysql" }}

func MySQLAndRedis() *conf.Bootstrap {
	host1, port1, err := mock.NewMySQL()
	if err != nil {
		panic(err)
	}
{{- if .Computed.enable_redis_final }}
	host2, port2, err := mock.NewRedis()
	if err != nil {
		panic(err)
	}
{{- end }}
	return &conf.Bootstrap{
		Server: &conf.Server{
			MachineId: "123",
		},
		Log: &conf.Log{
			Level:   "debug",
			JSON:    false,
			ShowSQL: true,
		},
		Db: &conf.DB{
			Driver:  "mysql",
			Dsn:     fmt.Sprintf("root:password@tcp(%s)/{{ .Computed.service_name_final }}?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local&timeout=10000ms", net.JoinHostPort(host1, fmt.Sprintf("%d", port1))),
			Migrate: false,
		},
{{- if .Computed.enable_redis_final }}
		Redis: &conf.Redis{
			Dsn: fmt.Sprintf("redis://%s", net.JoinHostPort(host2, fmt.Sprintf("%d", port2))),
		},
{{- end }}
{{- if .Computed.enable_trace_final }}
		Tracer: &conf.Tracer{
			Enable: true,
			Otlp:   &conf.Tracer_Otlp{},
			Stdout: &conf.Tracer_Stdout{},
		},
{{- end }}
	}
}
{{- else }}

func PostgreSQLAndRedis() *conf.Bootstrap {
	// Use real database connection for PostgreSQL
	host1, port1 := "localhost", 5432
{{- if .Computed.enable_redis_final }}
	host2, port2, err := mock.NewRedis()
	if err != nil {
		panic(err)
	}
{{- end }}
	return &conf.Bootstrap{
		Server: &conf.Server{
			MachineId: "123",
		},
		Log: &conf.Log{
			Level:   "debug",
			JSON:    false,
			ShowSQL: true,
		},
		Db: &conf.DB{
			Driver:  "postgres",
			Dsn:     fmt.Sprintf("host=%s user=root password=password dbname={{ .Computed.service_name_final }} port=%d sslmode=disable TimeZone=UTC", host1, port1),
			Migrate: false,
		},
{{- if .Computed.enable_redis_final }}
		Redis: &conf.Redis{
			Dsn: fmt.Sprintf("redis://%s", net.JoinHostPort(host2, fmt.Sprintf("%d", port2))),
		},
{{- end }}
{{- if .Computed.enable_trace_final }}
		Tracer: &conf.Tracer{
			Enable: true,
			Otlp:   &conf.Tracer_Otlp{},
			Stdout: &conf.Tracer_Stdout{},
		},
{{- end }}
	}
}
{{- end }}
{{- else }}

func Data() (c *conf.Bootstrap, dataData *data.Data, cache {{- if .Computed.enable_cache_final }} biz.Cache{{- else }} any{{- end }}) {
	once.Do(func() {
		onceC = &conf.Bootstrap{
			Name: "{{ .Computed.service_name_final }}",
		}
		onceData, _ = data.NewData(onceC)
		onceCache = nil
	})
	return onceC, onceData, onceCache
}
{{- end }}
