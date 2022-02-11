package server

import (
	"net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	v1 "github.com/qiangyt/jog/api/go/helloworld/v1"
	"github.com/qiangyt/jog/pkg/server/conf"
	"github.com/qiangyt/jog/pkg/server/service"
	_ "github.com/qiangyt/jog/web/statik"
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

	// web socket
	wsRouter := mux.NewRouter()
	wsRouter.HandleFunc("/web/ws", WsHandler)
	srv.HandlePrefix("/", wsRouter)

	// web static
	webFS, err := statikFs.NewWithNamespace("web")
	if err != nil {
		panic(errors.Wrapf(err, "failed to load web fs"))
	}
	srv.HandlePrefix("/web/", http.StripPrefix("/web/", http.FileServer(webFS)))

	// resource
	resFS, err := statikFs.NewWithNamespace("res")
	if err != nil {
		panic(errors.Wrapf(err, "failed to load res fs"))
	}
	srv.HandlePrefix("/res/", http.StripPrefix("/res/", http.FileServer(resFS)))

	// default route
	srv.Handle("/", http.RedirectHandler("/web/", 301))

	return srv
}
