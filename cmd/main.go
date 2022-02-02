package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/qiangyt/jog/common"
	"github.com/qiangyt/jog/convert"
	"github.com/qiangyt/jog/server"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../server", "config path, eg: -conf config.yaml")
}

func main() {
	ok, globalOptions := common.GlobalOptionsWithCommandLine(Version)
	if !ok {
		return
	}

	if !globalOptions.Debug {
		defer func() {
			if p := recover(); p != nil {
				color.Red.Printf("%v\n\n", p)
				os.Exit(1)
				return
			}
		}()
	}

	if globalOptions.RunMode == common.RunMode_Server {
		ok, _ := server.NewOptionsWithCommandLine(globalOptions.SubArgs)
		if !ok {
			return
		}
	} else {
		convert.Main(globalOptions.SubArgs, Version)
	}

	fmt.Println()
}
