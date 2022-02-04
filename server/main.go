package server

import (
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/qiangyt/jog/server/conf"
	"github.com/qiangyt/jog/util"
)

func Main(version string, args []string) {
	configFileUrl := conf.ParseCommandLine(args)
	bc := conf.LoadConfigFile(configFileUrl)

	var logFile util.LogFile
	var logger log.Logger

	if true { //conf.Log_Target() == conf.Log_stdio {
		logger = log.With(log.NewStdLogger(os.Stdout),
			"ts", log.DefaultTimestamp,
			"caller", log.DefaultCaller,
			"trace_id", tracing.TraceID(),
			"span_id", tracing.SpanID(),
		)
	} else {
		logFile = util.NewLogFile(bc.Log.GetFilePath())
		defer logFile.Close()

		logger = log.With(log.NewStdLogger(logFile.File()),
			"ts", log.DefaultTimestamp,
			"caller", log.DefaultCaller,
			"trace_id", tracing.TraceID(),
			"span_id", tracing.SpanID(),
		)
	}

	server, cleanup, err := initServer(logger, ServerVersion(version), bc, bc.Server, bc.Data)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := server.Run(); err != nil {
		panic(err)
	}
}

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer)

type ServerVersion string

func NewServer(logger log.Logger, version ServerVersion, config *conf.Bootstrap, hs *http.Server, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(config.Server.GetId()),
		kratos.Name("jog-server"),
		kratos.Version(string(version)),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(hs, gs),
	)
}
