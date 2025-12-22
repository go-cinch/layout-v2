{{ if and .Computed.enable_db_final (eq .Computed.orm_type_final "gorm") (not .Computed.enable_gen_model_final) }}
package model

import (
	"github.com/golang-module/carbon/v2"
)

// {{ .Computed.service_name_capitalized }} represents the {{ .Computed.service_name_final }} table model
// This is a manually-maintained model file (gentool auto-generation is disabled).
// You can manually define your models here, or enable gentool for automatic generation.
//
// Note: This example uses carbon.DateTime. If you prefer time.Time or int64:
//   - For time.Time: import "time", use *time.Time
//   - For Unix milliseconds: remove import, use *int64
type {{ .Computed.service_name_capitalized }} struct {
	ID        uint64           `gorm:"primaryKey;autoIncrement;comment:auto increment id"`
	CreatedAt *carbon.DateTime `gorm:"autoCreateTime:milli;comment:create time"`
	UpdatedAt *carbon.DateTime `gorm:"autoUpdateTime:milli;comment:update time"`
	Name      string           `gorm:"size:50;comment:name"`
}

// TableName overrides the table name
func ({{ .Computed.service_name_capitalized }}) TableName() string {
	return "{{ .Computed.service_name_final }}"
}
{{ end }}
