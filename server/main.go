package server

import (
	"flag"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/qiangyt/jog/server/conf"
)

type ServerId string
type ServerName string
type ServerVersion string

type ServerT struct {
	app     *kratos.App
	id      string
	name    string
	version string
}

type Server = *ServerT

func newServer(id ServerId, name ServerName, version ServerVersion, logger log.Logger, hs *http.Server, gs *grpc.Server) Server {
	app := kratos.New(
		kratos.ID(string(id)),
		kratos.Name(string(name)),
		kratos.Version(string(version)),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
	)

	return &ServerT{app: app, id: string(id), name: string(name), version: string(version)}
}

func (i Server) Id() string {
	return i.id
}

func (i Server) Name() string {
	return i.name
}

func (i Server) Version() string {
	return i.version
}

func (i Server) Run() error {
	return i.Run()
}

func Main(id string, name string, version string, flagconf string) {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", name,
		"service.version", version,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	server, cleanup, err := initServer(ServerId(id), ServerName(name), ServerVersion(version), bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := server.Run(); err != nil {
		panic(err)
	}
}
