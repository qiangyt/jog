package main

import (
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/qiangyt/jog/convert"
	"github.com/qiangyt/jog/server"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Version is the version of the compiled software.
	Version string
)

func main() {
	ok, options := GlobalOptionsWithCommandLine(Version)
	if !ok {
		return
	}

	if !options.Debug() {
		defer func() {
			if p := recover(); p != nil {
				color.Red.Printf("%v\n\n", p)
				os.Exit(1)
				return
			}
		}()
	}

	if options.RunMode() == RunMode_Client {
		convertDone := make(chan bool)
		convertCtx := convert.Main(convertDone, options.Debug(), options.SubArgs(), Version)
		if !convertCtx.Options().OpenWebGUI {
			<-convertDone
			fmt.Println()
			return
		}
	}

	server.Main(Version, options.SubArgs())
}
