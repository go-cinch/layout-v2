{{ if and .Computed.enable_db_final (eq .Scaffold.proto_template "full") }}
package data

import (
	"context"
{{- if eq .Computed.orm_type_final "none" }}
	"database/sql"
	"fmt"
	"strings"
{{- end }}
{{- if eq .Computed.orm_type_final "gorm" }}
	"gorm.io/gorm"
{{- end }}

	"github.com/go-cinch/common/log"
{{- if eq .Computed.orm_type_final "gorm" }}
	"github.com/go-cinch/common/copierx"
	"github.com/go-cinch/common/utils"
{{- end }}
	"github.com/google/wire"

	"{{.Computed.module_name_final}}/internal/biz"
{{- if eq .Computed.orm_type_final "gorm" }}
	"{{.Computed.module_name_final}}/internal/data/model"
{{- end }}
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
{{- if and .Computed.enable_redis_final .Computed.enable_cache_final }}
	NewRedis,
	NewCache,
{{- end }}
{{- if .Computed.enable_biz_tx_final }}
	NewTransaction,
{{- end }}
	New{{ .Computed.service_name_capitalized }}Repo,
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
func (ro {{ .Computed.service_name_final }}Repo) Create(ctx context.Context, item *biz.Create{{ .Computed.service_name_capitalized }}) (err error) {
	// Check if name exists
	var count int
	checkSQL := "SELECT COUNT(*) FROM {{ .Computed.service_name_final }} WHERE name = $1"
	err = ro.data.DB.QueryRowContext(ctx, checkSQL, item.Name).Scan(&count)
	if err != nil {
		log.WithError(err).Error("check name exists failed")
		return
	}
	if count > 0 {
		err = biz.ErrDuplicateField(ctx, "name", item.Name)
		return
	}

	insertSQL := "INSERT INTO {{ .Computed.service_name_final }} (id, name, created_at, updated_at) VALUES ($1, $2, NOW(), NOW())"
	if item.ID == 0 {
		item.ID = ro.data.ID(ctx)
	}
	_, err = ro.data.DB.ExecContext(ctx, insertSQL, item.ID, item.Name)
	if err != nil {
		log.WithError(err).Error("create {{ .Computed.service_name_final }} failed")
	}
	return
}

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
func (ro {{ .Computed.service_name_final }}Repo) Create(ctx context.Context, item *biz.Create{{ .Computed.service_name_capitalized }}) (err error) {
	// Check if name exists
	var count int
	checkSQL := "SELECT COUNT(*) FROM {{ .Computed.service_name_final }} WHERE name = ?"
	err = ro.data.DB.QueryRowContext(ctx, checkSQL, item.Name).Scan(&count)
	if err != nil {
		log.WithError(err).Error("check name exists failed")
		return
	}
	if count > 0 {
		err = biz.ErrDuplicateField(ctx, "name", item.Name)
		return
	}

	insertSQL := "INSERT INTO {{ .Computed.service_name_final }} (id, name, created_at, updated_at) VALUES (?, ?, NOW(), NOW())"
	if item.ID == 0 {
		item.ID = ro.data.ID(ctx)
	}
	_, err = ro.data.DB.ExecContext(ctx, insertSQL, item.ID, item.Name)
	if err != nil {
		log.WithError(err).Error("create {{ .Computed.service_name_final }} failed")
	}
	return
}

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
{{ end }}

func (ro {{ .Computed.service_name_final }}Repo) Find(ctx context.Context, condition *biz.Find{{ .Computed.service_name_capitalized }}) (rp []biz.{{ .Computed.service_name_capitalized }}) {
	rp = make([]biz.{{ .Computed.service_name_capitalized }}, 0)

	// Get total count for pagination first
	if !condition.Page.Disable {
		countQuery := "SELECT COUNT(*) FROM {{ .Computed.service_name_final }} WHERE 1=1"
		countArgs := make([]interface{}, 0)
		if condition.Name != nil {
			countQuery += " AND name LIKE ?"
			countArgs = append(countArgs, "%"+*condition.Name+"%")
		}
		var total int64
		_ = ro.data.DB.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
		condition.Page.Total = total

		// Early return if no results
		if total == 0 {
			return
		}
	}

	query := "SELECT id, name FROM {{ .Computed.service_name_final }} WHERE 1=1"
	args := make([]interface{}, 0)

	if condition.Name != nil {
		query += " AND name LIKE ?"
		args = append(args, "%"+*condition.Name+"%")
	}

	query += " ORDER BY id DESC"

	// Apply pagination using page/v2
	if !condition.Page.Disable {
		limit, offset := condition.Page.Limit()
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	rows, err := ro.data.DB.QueryContext(ctx, query, args...)
	if err != nil {
		log.WithError(err).Error("find {{ .Computed.service_name_final }} failed")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item biz.{{ .Computed.service_name_capitalized }}
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			log.WithError(err).Warn("scan {{ .Computed.service_name_final }} row failed")
			continue
		}
		rp = append(rp, item)
	}

	return
}

func (ro {{ .Computed.service_name_final }}Repo) Update(ctx context.Context, item *biz.Update{{ .Computed.service_name_capitalized }}) (err error) {
	// Check if record exists
	var exists int
	checkSQL := "SELECT COUNT(*) FROM {{ .Computed.service_name_final }} WHERE id = ?"
	err = ro.data.DB.QueryRowContext(ctx, checkSQL, item.ID).Scan(&exists)
	if err != nil {
		log.WithError(err).Error("check {{ .Computed.service_name_final }} exists failed")
		return
	}
	if exists == 0 {
		err = biz.ErrRecordNotFound(ctx)
		return
	}

	// Check name uniqueness if updating name
	if item.Name != nil {
		var count int
		nameCheckSQL := "SELECT COUNT(*) FROM {{ .Computed.service_name_final }} WHERE name = ? AND id != ?"
		err = ro.data.DB.QueryRowContext(ctx, nameCheckSQL, *item.Name, item.ID).Scan(&count)
		if err != nil {
			log.WithError(err).Error("check name uniqueness failed")
			return
		}
		if count > 0 {
			err = biz.ErrDuplicateField(ctx, "name", *item.Name)
			return
		}
	}

	updateSQL := "UPDATE {{ .Computed.service_name_final }} SET updated_at = NOW()"
	args := make([]interface{}, 0)

	if item.Name != nil {
		updateSQL += ", name = ?"
		args = append(args, *item.Name)
	}

	updateSQL += " WHERE id = ?"
	args = append(args, item.ID)

	_, err = ro.data.DB.ExecContext(ctx, updateSQL, args...)
	if err != nil {
		log.WithError(err).Error("update {{ .Computed.service_name_final }} failed")
	}
	return
}

func (ro {{ .Computed.service_name_final }}Repo) Delete(ctx context.Context, ids ...uint64) (err error) {
	if len(ids) == 0 {
		return
	}

	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1]

	deleteSQL := fmt.Sprintf("DELETE FROM {{ .Computed.service_name_final }} WHERE id IN (%s)", placeholders)
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	_, err = ro.data.DB.ExecContext(ctx, deleteSQL, args...)
	if err != nil {
		log.WithError(err).Error("delete {{ .Computed.service_name_final }} failed")
	}
	return
}
{{- end }}
{{- if eq .Computed.orm_type_final "gorm" }}

// Create creates a new {{ .Computed.service_name_final }} record using GORM generics API.
func (ro {{ .Computed.service_name_final }}Repo) Create(ctx context.Context, item *biz.Create{{ .Computed.service_name_capitalized }}) (err error) {
	db := gorm.G[model.{{ .Computed.service_name_capitalized }}](ro.data.DB(ctx))

	// Check if name exists
	count, err := db.Where("name = ?", item.Name).Count(ctx, "*")
	if err != nil {
		log.WithError(err).Error("check name exists failed")
		return
	}
	if count > 0 {
		err = biz.ErrDuplicateField(ctx, "name", item.Name)
		return
	}

	if item.ID == 0 {
		item.ID = ro.data.ID(ctx)
	}

	m := model.{{ .Computed.service_name_capitalized }}{
		ID:   item.ID,
		Name: &item.Name,
	}
	err = db.Create(ctx, &m)
	if err != nil {
		log.WithError(err).Error("create {{ .Computed.service_name_final }} failed")
	}
	return
}

func (ro {{ .Computed.service_name_final }}Repo) Get(ctx context.Context, id uint64) (item *biz.{{ .Computed.service_name_capitalized }}, err error) {
	db := gorm.G[model.{{ .Computed.service_name_capitalized }}](ro.data.DB(ctx))
	item = &biz.{{ .Computed.service_name_capitalized }}{}

	m, err := db.Where("id = ?", id).First(ctx)
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

func (ro {{ .Computed.service_name_final }}Repo) Find(ctx context.Context, condition *biz.Find{{ .Computed.service_name_capitalized }}) (rp []biz.{{ .Computed.service_name_capitalized }}) {
	rp = make([]biz.{{ .Computed.service_name_capitalized }}, 0)
	db := gorm.G[model.{{ .Computed.service_name_capitalized }}](ro.data.DB(ctx))

	// Apply filters
	if condition.Name != nil {
		db.Where("name LIKE ?", "%"+*condition.Name+"%")
	}

	// Count total before pagination
	if !condition.Page.Disable {
		count, err := db.Count(ctx, "*")
		if err != nil {
			log.WithError(err).Error("count {{ .Computed.service_name_final }} failed")
			return
		}
		condition.Page.Total = count

		// Early return if no results
		if count == 0 {
			return
		}
	}

	// Apply ordering
	db.Order("id DESC")

	// Apply pagination
	if !condition.Page.Disable {
		limit, offset := condition.Page.Limit()
		db.Limit(int(limit)).Offset(int(offset))
	}

	// Execute query
	list, err := db.Find(ctx)
	if err != nil {
		log.WithError(err).Error("find {{ .Computed.service_name_final }} failed")
		return
	}

	copierx.Copy(&rp, list)
	return
}

func (ro {{ .Computed.service_name_final }}Repo) Update(ctx context.Context, item *biz.Update{{ .Computed.service_name_capitalized }}) (err error) {
	db := gorm.G[model.{{ .Computed.service_name_capitalized }}](ro.data.DB(ctx))

	// Get existing record
	m, err := db.Where("id = ?", item.ID).First(ctx)
	if err == gorm.ErrRecordNotFound {
		err = biz.ErrRecordNotFound(ctx)
		return
	}
	if err != nil {
		log.WithError(err).Error("get {{ .Computed.service_name_final }} failed")
		return
	}

	change := make(map[string]interface{})
	utils.CompareDiff(m, item, &change)
	if len(change) == 0 {
		err = biz.ErrDataNotChange(ctx)
		return
	}

	// Check name uniqueness if name is being updated
	if item.Name != nil && (m.Name == nil || *item.Name != *m.Name) {
		count, err := db.Where("name = ? AND id != ?", *item.Name, item.ID).Count(ctx, "*")
		if err != nil {
			log.WithError(err).Error("check name uniqueness failed")
			return err
		}
		if count > 0 {
			err = biz.ErrDuplicateField(ctx, "name", *item.Name)
			return err
		}
	}

	// Update with changes map
	// Note: Use native DB.Updates for map updates, gorm.G.Updates expects struct type
	err = ro.data.DB(ctx).Where("id = ?", item.ID).Updates(change).Error
	if err != nil {
		log.WithError(err).Error("update {{ .Computed.service_name_final }} failed")
	}
	return
}

func (ro {{ .Computed.service_name_final }}Repo) Delete(ctx context.Context, ids ...uint64) (err error) {
	db := gorm.G[model.{{ .Computed.service_name_capitalized }}](ro.data.DB(ctx))

	_, err = db.Where("id IN ?", ids).Delete(ctx)
	if err != nil {
		log.WithError(err).Error("delete {{ .Computed.service_name_final }} failed")
	}
	return
}
{{- end }}
{{- end }}
