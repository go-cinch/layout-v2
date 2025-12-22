package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-cinch/common/page/v2"
	"github.com/go-cinch/common/proto/params"
	"github.com/google/wire"
	"google.golang.org/protobuf/types/known/emptypb"
{{- if .Computed.enable_trace_final }}
	"go.opentelemetry.io/otel"
{{- end }}

	v1 "{{.Computed.module_name_final}}/api/{{.Computed.service_name_final}}"
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
func (s *{{.Computed.service_name_capitalized}}Service) Create{{.Computed.service_name_capitalized}}(ctx context.Context, req *v1.Create{{.Computed.service_name_capitalized}}Request) (*emptypb.Empty, error) {
	{{- if .Computed.enable_trace_final }}
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Create{{.Computed.service_name_capitalized}}")
	defer span.End()
	{{- end }}
	item := &biz.Create{{.Computed.service_name_capitalized}}{}
	item.Name = req.Name
	err := s.uc.Create(ctx, item)
	return &emptypb.Empty{}, err
}

// Get{{.Computed.service_name_capitalized}} gets a record by id.
func (s *{{.Computed.service_name_capitalized}}Service) Get{{.Computed.service_name_capitalized}}(ctx context.Context, req *v1.Get{{.Computed.service_name_capitalized}}Request) (*v1.Get{{.Computed.service_name_capitalized}}Reply, error) {
	{{- if .Computed.enable_trace_final }}
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Get{{.Computed.service_name_capitalized}}")
	defer span.End()
	{{- end }}
	item, err := s.uc.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	reply := &v1.Get{{.Computed.service_name_capitalized}}Reply{
		Id:   item.ID,
		Name: item.Name,
	}
	return reply, nil
}

// Find{{.Computed.service_name_capitalized}} finds records by page.
func (s *{{.Computed.service_name_capitalized}}Service) Find{{.Computed.service_name_capitalized}}(ctx context.Context, req *v1.Find{{.Computed.service_name_capitalized}}Request) (*v1.Find{{.Computed.service_name_capitalized}}Reply, error) {
	{{- if .Computed.enable_trace_final }}
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Find{{.Computed.service_name_capitalized}}")
	defer span.End()
	{{- end }}
	condition := &biz.Find{{.Computed.service_name_capitalized}}{
		Page: page.Page{
			Num:     req.Page.Num,
			Size:    req.Page.Size,
			Disable: req.Page.Disable,
		},
		Name: req.Name,
	}
	items, err := s.uc.Find(ctx, condition)
	if err != nil {
		return nil, err
	}
	reply := &v1.Find{{.Computed.service_name_capitalized}}Reply{
		Page: &params.Page{
			Num:     condition.Page.Num,
			Size:    condition.Page.Size,
			Total:   condition.Page.Total,
			Disable: condition.Page.Disable,
		},
		List: make([]*v1.{{.Computed.service_name_capitalized}}Reply, 0, len(items)),
	}
	for _, item := range items {
		reply.List = append(reply.List, &v1.{{.Computed.service_name_capitalized}}Reply{
			Id:   item.ID,
			Name: item.Name,
		})
	}
	return reply, nil
}

// Update{{.Computed.service_name_capitalized}} updates a record by id.
func (s *{{.Computed.service_name_capitalized}}Service) Update{{.Computed.service_name_capitalized}}(ctx context.Context, req *v1.Update{{.Computed.service_name_capitalized}}Request) (*emptypb.Empty, error) {
	{{- if .Computed.enable_trace_final }}
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Update{{.Computed.service_name_capitalized}}")
	defer span.End()
	{{- end }}
	item := &biz.Update{{.Computed.service_name_capitalized}}{
		ID:   req.Id,
		Name: req.Name,
	}
	err := s.uc.Update(ctx, item)
	return &emptypb.Empty{}, err
}

// Delete{{.Computed.service_name_capitalized}} deletes records by ids.
func (s *{{.Computed.service_name_capitalized}}Service) Delete{{.Computed.service_name_capitalized}}(ctx context.Context, req *params.IdsRequest) (*emptypb.Empty, error) {
	{{- if .Computed.enable_trace_final }}
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "Delete{{.Computed.service_name_capitalized}}")
	defer span.End()
	{{- end }}
	ids := make([]uint64, 0)
	arr := strings.Split(req.Ids, ",")
	for _, item := range arr {
		id, err := strconv.ParseUint(strings.TrimSpace(item), 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	err := s.uc.Delete(ctx, ids...)
	return &emptypb.Empty{}, err
}
