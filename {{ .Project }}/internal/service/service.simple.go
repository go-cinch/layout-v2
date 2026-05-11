package service

import (
	"context"

	"github.com/google/wire"
{{- if .Computed.enable_trace_final }}
	"go.opentelemetry.io/otel"
{{- end }}

	v1 "{{.Computed.module_name_final}}/api/{{.Computed.service_name_kebab}}"
{{- if .Computed.enable_db_final }}
	"{{.Computed.module_name_final}}/internal/biz"
{{- end }}
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(New{{.Computed.service_name_capitalized}}Service)

// {{.Computed.service_name_capitalized}}Service is a {{.Computed.service_name_final}} service.
type {{.Computed.service_name_capitalized}}Service struct {
	v1.Unimplemented{{.Computed.service_name_capitalized}}Server
{{- if .Computed.enable_db_final }}
	uc *biz.{{.Computed.service_name_capitalized}}UseCase
{{- end }}
}

// New{{.Computed.service_name_capitalized}}Service creates a new {{.Computed.service_name_final}} service.
func New{{.Computed.service_name_capitalized}}Service({{ if .Computed.enable_db_final }}uc *biz.{{.Computed.service_name_capitalized}}UseCase{{ end }}) *{{.Computed.service_name_capitalized}}Service {
	return &{{.Computed.service_name_capitalized}}Service{
{{- if .Computed.enable_db_final }}
		uc: uc,
{{- end }}
	}
}

// Get{{.Computed.service_name_capitalized}} gets a record by id.
func (s *{{.Computed.service_name_capitalized}}Service) Get{{.Computed.service_name_capitalized}}(ctx context.Context, req *v1.Get{{.Computed.service_name_capitalized}}Request) (*v1.Get{{.Computed.service_name_capitalized}}Reply, error) {
{{- if .Computed.enable_trace_final }}
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "Get{{.Computed.service_name_capitalized}}")
	defer span.End()
{{- end }}
{{- if .Computed.enable_db_final }}
	res, err := s.uc.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.Get{{.Computed.service_name_capitalized}}Reply{
		Id:   res.ID,
		Name: res.Name,
	}, nil
{{- else }}
	return &v1.Get{{.Computed.service_name_capitalized}}Reply{
		Id:   req.Id,
		Name: "example",
	}, nil
{{- end }}
}
