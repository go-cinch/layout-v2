//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"

	"{{.Computed.module_name_final}}/internal/conf"
{{- if and (eq .Scaffold.proto_template "full") .Computed.enable_db_final }}
	"{{.Computed.module_name_final}}/internal/biz"
	"{{.Computed.module_name_final}}/internal/data"
{{- end }}
	"{{.Computed.module_name_final}}/internal/server"
	"{{.Computed.module_name_final}}/internal/service"
)

// wireApp initializes the Kratos application.
func wireApp(*conf.Bootstrap) (*kratos.App, func(), error) {
	panic(wire.Build(
{{- if and (eq .Scaffold.proto_template "full") .Computed.enable_db_final }}
		data.ProviderSet,
		biz.ProviderSet,
{{- end }}
		server.ProviderSet,
		service.ProviderSet,
		newApp,
	))
}
