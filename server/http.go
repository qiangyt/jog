package server

import (
	"net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/pkg/errors"
	v1 "github.com/qiangyt/jog/api/helloworld/v1"
	"github.com/qiangyt/jog/server/conf"
	"github.com/qiangyt/jog/server/service"
	_ "github.com/qiangyt/jog/statik"
	statikFs "github.com/rakyll/statik/fs"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *kratosHttp.Server {
	var opts = []kratosHttp.ServerOption{
		kratosHttp.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, kratosHttp.Network(c.Http.Network))
	}

	opts = append(opts, kratosHttp.Address(c.Http.Addr))
	opts = append(opts, kratosHttp.Timeout(c.Http.Timeout.AsDuration()))

	srv := kratosHttp.NewServer(opts...)

	v1.RegisterGreeterHTTPServer(srv, greeter)

	statikFS, err := statikFs.New()
	if err != nil {
		panic(errors.Wrapf(err, "failed to load statik fs"))
	}

	srv.Handle("/", http.RedirectHandler("/web/", 301))
	srv.HandlePrefix("/web/", http.StripPrefix("/web/", http.FileServer(statikFS)))

	return srv
}
