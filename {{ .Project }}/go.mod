module {{.Computed.module_name_final}}

go 1.25

require (
	github.com/go-cinch/common/log v1.2.0
	github.com/go-cinch/common/proto/params v1.0.1
	{{- if .Computed.enable_reason_final }}
	github.com/go-kratos/kratos/v2 v2.8.3
	{{- else }}
	github.com/go-kratos/kratos/v2 v2.8.3
	{{- end }}
	{{- if .Computed.enable_i18n_final }}
	github.com/go-cinch/common/middleware/i18n v1.0.2
	{{- end }}
	{{- if .Computed.enable_db_final }}
	{{- if eq .Computed.orm_type_final "gorm" }}
	github.com/go-cinch/common/id v1.0.6
	github.com/go-cinch/common/plugins/gorm/tenant/v2 v2.0.2
	github.com/go-cinch/common/middleware/tenant/v2 v2.0.1
	github.com/go-cinch/common/copierx v1.0.4
	github.com/go-cinch/common/constant v1.0.5
	github.com/go-cinch/common/utils v1.0.5
	gorm.io/gorm v1.31.1
	{{- if eq .Computed.db_type_final "mysql" }}
	gorm.io/driver/mysql v1.5.7
	{{- else if eq .Computed.db_type_final "postgres" }}
	gorm.io/driver/postgres v1.5.11
	{{- end }}
	{{- else if eq .Computed.orm_type_final "none" }}
	{{- if eq .Computed.db_type_final "mysql" }}
	github.com/go-sql-driver/mysql v1.8.1
	{{- else if eq .Computed.db_type_final "postgres" }}
	github.com/lib/pq v1.10.9
	{{- end }}
	{{- end }}
	{{- end }}
	github.com/google/gnostic v0.7.0
	github.com/google/wire v0.6.0
	google.golang.org/genproto/googleapis/api v0.0.0-20250115164207-1a7da9e5054f
	google.golang.org/grpc v1.69.4
	google.golang.org/protobuf v1.36.5
	{{- if .Computed.enable_trace_final }}
	go.opentelemetry.io/otel v1.34.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.34.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.34.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.34.0
	go.opentelemetry.io/otel/sdk v1.34.0
	{{- end }}
	{{- if .Computed.enable_ws_final }}
	github.com/gorilla/websocket v1.5.3
	{{- end }}
)
