package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	_ "github.com/google/gnostic/openapiv3"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	"{{.Computed.common_module_final}}/log"
	"{{.Computed.common_module_final}}/log/caller"
	"{{.Computed.common_module_final}}/plugins/k8s/pod"
	"{{.Computed.common_module_final}}/plugins/kratos/config/env"
	_ "{{.Computed.common_module_final}}/plugins/kratos/encoding/yml"
	_ "{{.Computed.common_module_final}}/proto/params"
	"{{.Computed.common_module_final}}/utils"

	"{{.Computed.module_name_final}}/internal/conf"
)

var (
	Name     = "{{.Computed.service_name_final}}"
	Version  = "1.0.0"
	flagConf string
	id, _    = os.Hostname()
)

func newApp(gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Logger(log.DefaultWrapper.Options().Logger()),
		kratos.Server(gs, hs),
	)
}

func main() {
	flag.StringVar(&flagConf, "c", "../../configs", "config path, eg: -c config.yaml")
	flag.Parse()

	// Initialize log
	logOps := make([]func(*log.Options), 0, 10)
	logOps = append(logOps,
		log.WithJSON(false),
		log.WithLevel(log.InfoLevel),
		log.WithValuer("service.id", id),
		log.WithValuer("service.name", Name),
		log.WithValuer("service.version", Version),
		log.WithValuer("trace.id", tracing.TraceID()),
		log.WithValuer("span.id", tracing.SpanID()),
		log.WithCallerOptions(
			caller.WithSource(false),
			caller.WithLevel(2),
			caller.WithVersion(true),
		),
	)
	log.DefaultWrapper = log.NewWrapper(logOps...)

	c := config.New(
		config.WithSource(file.NewSource(flagConf)),
		config.WithResolver(
			env.NewRevolver(
				env.WithPrefix("SERVICE"),
				env.WithLoaded(func(k string, v interface{}) {
					// mask sensitive fields before logging
					val := fmt.Sprint(v)
					lowerKey := strings.ToLower(k)
					if strings.Contains(lowerKey, "password") ||
						strings.Contains(lowerKey, "secret") ||
						strings.Contains(lowerKey, "token") ||
						strings.Contains(lowerKey, "key") {
						switch {
						case len(val) > 6:
							val = val[:3] + "***" + val[len(val)-3:]
						case len(val) > 0:
							val = "***"
						}
					}
					log.Info("env loaded: %s=%v", k, val)
				}),
			),
		),
	)
	defer c.Close()

	fields := log.Fields{
		"conf": flagConf,
	}
	if err := c.Load(); err != nil {
		log.
			WithError(err).
			WithFields(fields).
			Fatal("load conf failed")
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		log.
			WithError(err).
			WithFields(fields).
			Fatal("scan conf failed")
	}
	bc.Name = Name
	bc.Version = Version
	// override log level after read config
	logOps = append(logOps,
		[]func(*log.Options){
			log.WithLevel(log.NewLevel(bc.Log.Level)),
			log.WithJSON(bc.Log.JSON),
		}...,
	)
	log.DefaultWrapper = log.NewWrapper(logOps...)
	if bc.Server.MachineId == "" {
		// if machine id not set, gen from pod ip
		machineId, err := pod.MachineID()
		if err == nil {
			bc.Server.MachineId = strconv.FormatUint(uint64(machineId), 10)
		} else {
			bc.Server.MachineId = "0"
		}
	}
	// os.Setenv("COPIERX_UTC", "true")

	app, cleanup, err := wireApp(&bc)
	if err != nil {
		str := utils.Struct2JSON(&bc)
		log.
			WithError(err).
			Error("wire app failed")
		// env str maybe very long, log with another line
		log.
			WithFields(fields).
			Fatal(str)
	}
	defer cleanup()

	// start and wait for stop signal
	if err = app.Run(); err != nil {
		log.
			WithError(err).
			WithFields(fields).
			Fatal("run app failed")
	}
}
