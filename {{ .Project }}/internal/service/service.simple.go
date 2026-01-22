package service

import (
	"context"

	"github.com/google/wire"
{{- if .Computed.enable_trace_final }}
	"go.opentelemetry.io/otel"
{{- end }}

	v1 "{{.Computed.module_name_final}}/api/{{.Computed.service_name_final}}"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(New{{.Computed.service_name_capitalized}}Service)

// {{.Computed.service_name_capitalized}}Service is a {{.Computed.service_name_final}} service.
type {{.Computed.service_name_capitalized}}Service struct {
	v1.Unimplemented{{.Computed.service_name_capitalized}}Server
}

// New{{.Computed.service_name_capitalized}}Service creates a new {{.Computed.service_name_final}} service.
func New{{.Computed.service_name_capitalized}}Service() *{{.Computed.service_name_capitalized}}Service {
	return &{{.Computed.service_name_capitalized}}Service{}
}

// Get{{.Computed.service_name_capitalized}} gets a record by id.
func (s *{{.Computed.service_name_capitalized}}Service) Get{{.Computed.service_name_capitalized}}(ctx context.Context, req *v1.Get{{.Computed.service_name_capitalized}}Request) (*v1.Get{{.Computed.service_name_capitalized}}Reply, error) {
{{- if .Computed.enable_trace_final }}
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "Get{{.Computed.service_name_capitalized}}")
	defer span.End()
{{- end }}
	return &v1.Get{{.Computed.service_name_capitalized}}Reply{
		Id:   req.Id,
		Name: "example",
	}, nil
}
