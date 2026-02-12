package service

import (
	"context"

	"{{.Computed.common_module_final}}/copierx"
	"{{.Computed.common_module_final}}/page/v2"
	"{{.Computed.common_module_final}}/proto/params"
	"{{.Computed.common_module_final}}/utils"
	"github.com/google/wire"
	"google.golang.org/protobuf/types/known/emptypb"
{{- if .Computed.enable_trace_final }}
	"go.opentelemetry.io/otel"
{{- end }}

	v1 "{{.Computed.module_name_final}}/api/{{.Computed.service_name_snake}}"
	"{{.Computed.module_name_final}}/internal/biz"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(New{{.Computed.service_name_capitalized}}Service)

// {{.Computed.service_name_capitalized}}Service is a {{.Computed.service_name_final}} service.
type {{.Computed.service_name_capitalized}}Service struct {
	v1.Unimplemented{{.Computed.service_name_capitalized}}Server
	uc *biz.{{.Computed.service_name_capitalized}}UseCase
}

// New{{.Computed.service_name_capitalized}}Service creates a new {{.Computed.service_name_final}} service.
func New{{.Computed.service_name_capitalized}}Service(uc *biz.{{.Computed.service_name_capitalized}}UseCase) *{{.Computed.service_name_capitalized}}Service {
	return &{{.Computed.service_name_capitalized}}Service{
		uc: uc,
	}
}

// Create{{.Computed.service_name_capitalized}} creates a new record.
func (s *{{.Computed.service_name_capitalized}}Service) Create{{.Computed.service_name_capitalized}}(ctx context.Context, req *v1.Create{{.Computed.service_name_capitalized}}Request) (rp *emptypb.Empty, err error) {
	{{- if .Computed.enable_trace_final }}
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "Create{{.Computed.service_name_capitalized}}")
	defer span.End()
	{{- end }}
	rp = &emptypb.Empty{}
	r := &biz.Create{{.Computed.service_name_capitalized}}{}
	copierx.Copy(&r, req)
	err = s.uc.Create(ctx, r)
	return
}

// Get{{.Computed.service_name_capitalized}} gets a record by id.
func (s *{{.Computed.service_name_capitalized}}Service) Get{{.Computed.service_name_capitalized}}(ctx context.Context, req *v1.Get{{.Computed.service_name_capitalized}}Request) (rp *v1.Get{{.Computed.service_name_capitalized}}Reply, err error) {
	{{- if .Computed.enable_trace_final }}
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "Get{{.Computed.service_name_capitalized}}")
	defer span.End()
	{{- end }}
	rp = &v1.Get{{.Computed.service_name_capitalized}}Reply{}
	res, err := s.uc.Get(ctx, req.Id)
	if err != nil {
		return
	}
	copierx.Copy(&rp, res)
	return
}

// Find{{.Computed.service_name_capitalized}} finds records by page.
func (s *{{.Computed.service_name_capitalized}}Service) Find{{.Computed.service_name_capitalized}}(ctx context.Context, req *v1.Find{{.Computed.service_name_capitalized}}Request) (rp *v1.Find{{.Computed.service_name_capitalized}}Reply, err error) {
	{{- if .Computed.enable_trace_final }}
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "Find{{.Computed.service_name_capitalized}}")
	defer span.End()
	{{- end }}
	rp = &v1.Find{{.Computed.service_name_capitalized}}Reply{}
	rp.Page = &params.Page{}
	r := &biz.Find{{.Computed.service_name_capitalized}}{}
	r.Page = page.Page{}
	copierx.Copy(&r, req)
	copierx.Copy(&r.Page, req.Page)
	res, err := s.uc.Find(ctx, r)
	if err != nil {
		return
	}
	copierx.Copy(&rp.Page, r.Page)
	copierx.Copy(&rp.List, res)
	return
}

// Update{{.Computed.service_name_capitalized}} updates a record by id.
func (s *{{.Computed.service_name_capitalized}}Service) Update{{.Computed.service_name_capitalized}}(ctx context.Context, req *v1.Update{{.Computed.service_name_capitalized}}Request) (rp *emptypb.Empty, err error) {
	{{- if .Computed.enable_trace_final }}
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "Update{{.Computed.service_name_capitalized}}")
	defer span.End()
	{{- end }}
	rp = &emptypb.Empty{}
	r := &biz.Update{{.Computed.service_name_capitalized}}{}
	copierx.Copy(&r, req)
	err = s.uc.Update(ctx, r)
	return
}

// Delete{{.Computed.service_name_capitalized}} deletes records by ids.
func (s *{{.Computed.service_name_capitalized}}Service) Delete{{.Computed.service_name_capitalized}}(ctx context.Context, req *params.IdsRequest) (rp *emptypb.Empty, err error) {
	{{- if .Computed.enable_trace_final }}
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "Delete{{.Computed.service_name_capitalized}}")
	defer span.End()
	{{- end }}
	rp = &emptypb.Empty{}
	err = s.uc.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...)
	return
}
