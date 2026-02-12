-- +migrate Up
{{- if eq .Computed.db_type_final "postgres" }}
CREATE TABLE t_{{ .Computed.service_name_snake }}
(
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP(3) NULL,
    updated_at TIMESTAMP(3) NULL,
    name       VARCHAR(50)
);

COMMENT ON COLUMN t_{{ .Computed.service_name_snake }}.id IS 'auto increment id';
COMMENT ON COLUMN t_{{ .Computed.service_name_snake }}.created_at IS 'create time';
COMMENT ON COLUMN t_{{ .Computed.service_name_snake }}.updated_at IS 'update time';
COMMENT ON COLUMN t_{{ .Computed.service_name_snake }}.name IS 'name';
{{- else }}
CREATE TABLE `t_{{ .Computed.service_name_snake }}`
(
    `id`         BIGINT UNSIGNED AUTO_INCREMENT COMMENT 'auto increment id' PRIMARY KEY,
    `created_at` DATETIME(3) NULL COMMENT 'create time',
    `updated_at` DATETIME(3) NULL COMMENT 'update time',
    `name`       VARCHAR(50) COMMENT 'name'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;
{{- end }}

-- +migrate Down
{{- if eq .Computed.db_type_final "postgres" }}
DROP TABLE t_{{ .Computed.service_name_snake }};
{{- else }}
DROP TABLE `t_{{ .Computed.service_name_snake }}`;
{{- end }}
