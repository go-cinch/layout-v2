{{ if .Computed.enable_db_final }}
package data

import (
	"context"
{{- if eq .Computed.orm_type_final "none" }}
	"database/sql"
{{- end }}
{{- if eq .Computed.orm_type_final "gorm" }}
	"gorm.io/gorm"
{{- end }}

	"{{.Computed.common_module_final}}/log"
{{- if eq .Computed.orm_type_final "gorm" }}
	"{{.Computed.common_module_final}}/copierx"
{{- end }}

	"{{.Computed.module_name_final}}/internal/biz"
{{- if eq .Computed.orm_type_final "gorm" }}
	"{{.Computed.module_name_final}}/internal/data/model"
{{- end }}
)

type {{ .Computed.service_name_final }}Repo struct {
	data *Data
}

func New{{ .Computed.service_name_capitalized }}Repo(data *Data) biz.{{ .Computed.service_name_capitalized }}Repo {
	return &{{ .Computed.service_name_final }}Repo{
		data: data,
	}
}
{{- if eq .Computed.orm_type_final "none" }}

{{ if eq .Computed.db_type_final "postgres" }}
// Get retrieves a record by ID using raw SQL queries.
func (ro {{ .Computed.service_name_final }}Repo) Get(ctx context.Context, id uint64) (item *biz.{{ .Computed.service_name_capitalized }}, err error) {
	item = &biz.{{ .Computed.service_name_capitalized }}{}
	query := "SELECT id, name FROM {{ .Computed.service_name_final }} WHERE id = $1"
	err = ro.data.DB.QueryRowContext(ctx, query, id).Scan(&item.ID, &item.Name)
	if err == sql.ErrNoRows {
		err = biz.ErrRecordNotFound(ctx)
		return
	}
	if err != nil {
		log.WithError(err).Error("get {{ .Computed.service_name_final }} failed")
	}
	return
}
{{ else }}
// Get retrieves a record by ID using raw SQL queries.
func (ro {{ .Computed.service_name_final }}Repo) Get(ctx context.Context, id uint64) (item *biz.{{ .Computed.service_name_capitalized }}, err error) {
	item = &biz.{{ .Computed.service_name_capitalized }}{}
	query := "SELECT id, name FROM {{ .Computed.service_name_final }} WHERE id = ?"
	err = ro.data.DB.QueryRowContext(ctx, query, id).Scan(&item.ID, &item.Name)
	if err == sql.ErrNoRows {
		err = biz.ErrRecordNotFound(ctx)
		return
	}
	if err != nil {
		log.WithError(err).Error("get {{ .Computed.service_name_final }} failed")
	}
	return
}
{{- end }}
{{- else if eq .Computed.orm_type_final "gorm" }}

// Get retrieves a record by ID using GORM Generics.
func (ro {{ .Computed.service_name_final }}Repo) Get(ctx context.Context, id uint64) (item *biz.{{ .Computed.service_name_capitalized }}, err error) {
	item = &biz.{{ .Computed.service_name_capitalized }}{}
	m, err := gorm.G[model.{{ .Computed.service_name_capitalized }}](ro.data.DB(ctx)).Where("id = ?", id).First(ctx)
	if err == gorm.ErrRecordNotFound {
		err = biz.ErrRecordNotFound(ctx)
		return
	}
	if err != nil {
		log.WithError(err).Error("get {{ .Computed.service_name_final }} failed")
		return
	}
	copierx.Copy(&item, m)
	return
}
{{- end }}
{{- end }}
