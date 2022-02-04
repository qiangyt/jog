package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/mitchellh/go-homedir"
	"github.com/qiangyt/jog/static"
	"github.com/qiangyt/jog/util"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

func LoadConfigFile(configFilePath string) *Bootstrap {
	kratosConfig := config.New(
		config.WithSource(
			file.NewSource(configFilePath),
		),
	)
	defer kratosConfig.Close()

	if err := kratosConfig.Load(); err != nil {
		panic(err)
	}

	var bc Bootstrap
	if err := kratosConfig.Scan(&bc); err != nil {
		panic(err)
	}

	normalizeLog(&bc)
	normalizeServer(&bc)

	return &bc
}

func normalizeServerHttp(server *Server) {
	http := server.Http
	if http == nil {
		server.Http = &Server_HTTP{}
		http = server.Http
	}

	if len(http.GetAddr()) == 0 {
		http.Addr = "0.0.0.0:8585"
	}
	if http.GetTimeout() == nil {
		http.Timeout = &durationpb.Duration{Seconds: 6}
	}
}

func normalizeServerGrpc(server *Server) {
	grpc := server.Grpc
	if grpc == nil {
		server.Grpc = &Server_GRPC{}
		grpc = server.Grpc
	}

	if len(grpc.GetAddr()) == 0 {
		grpc.Addr = "0.0.0.0:9595"
	}
	if grpc.GetTimeout() == nil {
		grpc.Timeout = &durationpb.Duration{Seconds: 6}
	}
}

func normalizeServer(bc *Bootstrap) {
	server := bc.Server
	if server == nil {
		bc.Server = &Server{}
		server = bc.Server
	}

	if len(server.GetId()) == 0 {
		server.Id, _ = os.Hostname()
	}

	normalizeServerHttp(server)
	normalizeServerGrpc(server)
}

func normalizeLog(bc *Bootstrap) {
	log := bc.Log
	if log == nil {
		bc.Log = &Log{}
		log = bc.Log
	}

	if len(log.Target) == 0 {
		log.Target = "stdio"
	}
	if log.Target == "file" {
		if len(log.FilePath) == 0 {
			log.FilePath = filepath.Join(util.ExeDirectory(), "jog.server.log")
		}
	}
}

func lookForConfigFile(logger *log.Helper, dir string) string {
	logger.Infof("looking for config files in directory %s", dir)
	r := filepath.Join(dir, "jog.server.yaml")
	if util.FileExists(r) {
		return r
	}
	r = filepath.Join(dir, "jog.server.yml")
	if util.FileExists(r) {
		return r
	}
	return ""
}

// DetermineConfigFilePath return (file path)
func determineConfigFilePath(logger *log.Helper) string {
	exeDir := util.ExeDirectory()
	r := lookForConfigFile(logger, exeDir)
	if len(r) != 0 {
		return r
	}

	homeDir, err := homedir.Dir()
	if err != nil {
		logger.Errorf("failed to get home dir. Error: %v", err)
	} else {
		r = lookForConfigFile(logger, homeDir)
	}
	if len(r) != 0 {
		return r
	}

	r = filepath.Join(exeDir, "jog.server.yaml")
	util.WriteFileIfNotFound(r, []byte(static.DefaultServer_yml))
	return r
}

func ParseCommandLine(args []string) string {
	var configFilePath string

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg[0:1] != "-" {
			util.PrintErrorHint("Invalid argument: '%s'", arg)
		} else {
			if arg == "-c" || arg == "--config" {
				if i+1 >= len(args) {
					util.PrintErrorHint("Missing config file path")
					return ""
				}

				configFilePath = args[i+1]
				i++
			} else if arg == "-t" || arg == "--template" {
				fmt.Println(static.DefaultServer_yml)
				return ""
			} else {
				util.PrintErrorHint("Unknown option: '%s'", arg)
				return ""
			}
		}
	}

	if len(configFilePath) == 0 {
		tmpLogger := log.NewHelper(log.With(log.NewStdLogger(os.Stdout),
			"ts", log.DefaultTimestamp,
			"caller", log.DefaultCaller,
		))

		configFilePath = determineConfigFilePath(tmpLogger)
	}
	return configFilePath
}
