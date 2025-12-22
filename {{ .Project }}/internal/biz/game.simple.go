{{ if .Computed.enable_db_final }}
package biz

import (
	"context"

	"{{.Computed.module_name_final}}/internal/conf"
)

type {{ .Computed.service_name_capitalized }} struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type {{ .Computed.service_name_capitalized }}Repo interface {
	Get(ctx context.Context, id uint64) (*{{ .Computed.service_name_capitalized }}, error)
}

type {{ .Computed.service_name_capitalized }}UseCase struct {
	c    *conf.Bootstrap
	repo {{ .Computed.service_name_capitalized }}Repo
}

func New{{ .Computed.service_name_capitalized }}UseCase(c *conf.Bootstrap, repo {{ .Computed.service_name_capitalized }}Repo) *{{ .Computed.service_name_capitalized }}UseCase {
	return &{{ .Computed.service_name_capitalized }}UseCase{
		c:    c,
		repo: repo,
	}
}

func (uc *{{ .Computed.service_name_capitalized }}UseCase) Get(ctx context.Context, id uint64) (*{{ .Computed.service_name_capitalized }}, error) {
	return uc.repo.Get(ctx, id)
}
{{ end }}
