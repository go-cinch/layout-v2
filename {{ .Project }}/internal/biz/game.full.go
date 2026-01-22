{{ if .Computed.enable_db_final }}
package biz

import (
	"context"
{{- if .Computed.enable_cache_final }}
	"strconv"
	"strings"
{{- end }}

{{- if .Computed.enable_cache_final }}
	"{{.Computed.common_module_final}}/constant"
	"{{.Computed.common_module_final}}/copierx"
{{- end }}
	"{{.Computed.common_module_final}}/page/v2"
{{- if .Computed.enable_cache_final }}
	"{{.Computed.common_module_final}}/utils"
	"github.com/pkg/errors"
{{- end }}

	"{{.Computed.module_name_final}}/internal/conf"
)

type Create{{ .Computed.service_name_capitalized }} struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type {{ .Computed.service_name_capitalized }} struct {
	ID   uint64 `json:"id,string"`
	Name string `json:"name"`
}

type Find{{ .Computed.service_name_capitalized }} struct {
	Page page.Page `json:"page"`
	Name *string   `json:"name"`
}

type Find{{ .Computed.service_name_capitalized }}Cache struct {
	Page page.Page `json:"page"`
	List []{{ .Computed.service_name_capitalized }}    `json:"list"`
}

type Update{{ .Computed.service_name_capitalized }} struct {
	ID   uint64  `json:"id,string"`
	Name *string `json:"name,omitempty"`
}

type {{ .Computed.service_name_capitalized }}UseCase struct {
	c     *conf.Bootstrap
	repo  {{ .Computed.service_name_capitalized }}Repo
{{- if .Computed.enable_biz_tx_final }}
	tx    Transaction
{{- end }}
{{- if .Computed.enable_cache_final }}
	cache Cache
{{- end }}
}

func New{{ .Computed.service_name_capitalized }}UseCase(c *conf.Bootstrap, repo {{ .Computed.service_name_capitalized }}Repo{{ if .Computed.enable_biz_tx_final }}, tx Transaction{{ end }}{{ if .Computed.enable_cache_final }}, cache Cache{{ end }}) *{{ .Computed.service_name_capitalized }}UseCase {
	return &{{ .Computed.service_name_capitalized }}UseCase{
		c:    c,
		repo: repo,
{{- if .Computed.enable_biz_tx_final }}
		tx:   tx,
{{- end }}
{{- if .Computed.enable_cache_final }}
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "{{ .Computed.service_name_final }}",
		}, "_")),
{{- end }}
	}
}

func (uc *{{ .Computed.service_name_capitalized }}UseCase) Create(ctx context.Context, item *Create{{ .Computed.service_name_capitalized }}) error {
{{- if .Computed.enable_biz_tx_final }}
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		{{- if .Computed.enable_cache_final }}
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
		{{- else }}
		return uc.repo.Create(ctx, item)
		{{- end }}
	})
{{- else }}
	{{- if .Computed.enable_cache_final }}
	return uc.cache.Flush(ctx, func(ctx context.Context) error {
		return uc.repo.Create(ctx, item)
	})
	{{- else }}
	return uc.repo.Create(ctx, item)
	{{- end }}
{{- end }}
}

func (uc *{{ .Computed.service_name_capitalized }}UseCase) Get(ctx context.Context, id uint64) (rp *{{ .Computed.service_name_capitalized }}, err error) {
{{- if .Computed.enable_cache_final }}
	rp = &{{ .Computed.service_name_capitalized }}{}
	action := strings.Join([]string{"get", strconv.FormatUint(id, 10)}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.get(ctx, action, id)
	})
	if err != nil {
		return
	}
	utils.JSON2Struct(&rp, str)
	if rp.ID == constant.UI0 {
		err = ErrRecordNotFound(ctx)
		return
	}
	return
{{- else }}
	return uc.repo.Get(ctx, id)
{{- end }}
}

{{- if .Computed.enable_cache_final }}
func (uc *{{ .Computed.service_name_capitalized }}UseCase) get(ctx context.Context, action string, id uint64) (res string, err error) {
	// read data from db and write to cache
	rp := &{{ .Computed.service_name_capitalized }}{}
	item, err := uc.repo.Get(ctx, id)
	notFound := errors.Is(err, ErrRecordNotFound(ctx))
	if err != nil && !notFound {
		return
	}
	copierx.Copy(&rp, item)
	res = utils.Struct2JSON(rp)
	uc.cache.Set(ctx, action, res, notFound)
	return
}

{{- end }}

func (uc *{{ .Computed.service_name_capitalized }}UseCase) Find(ctx context.Context, condition *Find{{ .Computed.service_name_capitalized }}) (rp []{{ .Computed.service_name_capitalized }}, err error) {
{{- if .Computed.enable_cache_final }}
	// use md5 string as cache replay json str, key is short
	action := strings.Join([]string{"find", utils.StructMd5(condition)}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return
	}
	var cache Find{{ .Computed.service_name_capitalized }}Cache
	utils.JSON2Struct(&cache, str)
	condition.Page = cache.Page
	rp = cache.List
	return
{{- else }}
	return uc.repo.Find(ctx, condition), nil
{{- end }}
}

{{- if .Computed.enable_cache_final }}
func (uc *{{ .Computed.service_name_capitalized }}UseCase) find(ctx context.Context, action string, condition *Find{{ .Computed.service_name_capitalized }}) (res string, err error) {
	// read data from db and write to cache
	list := uc.repo.Find(ctx, condition)
	var cache Find{{ .Computed.service_name_capitalized }}Cache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2JSON(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return
}

{{- end }}

func (uc *{{ .Computed.service_name_capitalized }}UseCase) Update(ctx context.Context, item *Update{{ .Computed.service_name_capitalized }}) error {
{{- if .Computed.enable_biz_tx_final }}
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		{{- if .Computed.enable_cache_final }}
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Update(ctx, item)
			return
		})
		{{- else }}
		return uc.repo.Update(ctx, item)
		{{- end }}
	})
{{- else }}
	{{- if .Computed.enable_cache_final }}
	return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
		err = uc.repo.Update(ctx, item)
		return
	})
	{{- else }}
	return uc.repo.Update(ctx, item)
	{{- end }}
{{- end }}
}

func (uc *{{ .Computed.service_name_capitalized }}UseCase) Delete(ctx context.Context, ids ...uint64) error {
{{- if .Computed.enable_biz_tx_final }}
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		{{- if .Computed.enable_cache_final }}
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			err = uc.repo.Delete(ctx, ids...)
			return
		})
		{{- else }}
		return uc.repo.Delete(ctx, ids...)
		{{- end }}
	})
{{- else }}
	{{- if .Computed.enable_cache_final }}
	return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
		err = uc.repo.Delete(ctx, ids...)
		return
	})
	{{- else }}
	return uc.repo.Delete(ctx, ids...)
	{{- end }}
{{- end }}
}
{{ end }}
