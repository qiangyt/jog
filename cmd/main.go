package main

import (
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/qiangyt/jog/common"
	"github.com/qiangyt/jog/convert"
	"github.com/qiangyt/jog/server"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Version is the version of the compiled software.
	Version string
)

func main() {
	ok, globalOptions := common.GlobalOptionsWithCommandLine(Version)
	if !ok {
		return
	}

	if !globalOptions.Debug() {
		defer func() {
			if p := recover(); p != nil {
				color.Red.Printf("%v\n\n", p)
				os.Exit(1)
				return
			}
		}()
	}

	if globalOptions.RunMode() == common.RunMode_Client {
		convertDone := make(chan bool)
		convertCtx := convert.Main(convertDone, globalOptions)
		if !convertCtx.Options().OpenWebGUI {
			<-convertDone
			fmt.Println()
			return
		}
	}

	server.Main(Version, globalOptions)
}
