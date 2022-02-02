package server

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer)

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

func NewServer(id ServerId, name ServerName, version ServerVersion, logger log.Logger, hs *http.Server, gs *grpc.Server) Server {
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
