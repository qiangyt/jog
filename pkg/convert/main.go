package convert

import (
	"os"

	"github.com/gookit/color"
	"github.com/qiangyt/jog/pkg/util"
)

func Main(done chan bool, debug bool, args []string, version string) ConvertContext {
	util.InitDefaultGrokLibraryDir()

	ok, options := NewOptionsWithCommandLine(args)
	if !ok {
		close(done)
		return nil
	}

	ctx := NewConvertContext(options, util.JogHomeDir(true), version)
	defer ctx.Close()

	if !ctx.LoadConfig() {
		close(done)
		return ctx
	}

	go func() {
		defer close(done)

		//TODO: trap CTRL+C signal
		if !debug {
			defer func() {
				if p := recover(); p != nil {
					color.Red.Printf("%v\n\n", p)
					os.Exit(1)
					return
				}
			}()
		}

		if len(options.LogFilePath) == 0 {
			ctx.LogInfo("read JSON log lines from stdin")
			ProcessReader(ctx, os.Stdin, 1)
		} else {
			ctx.LogInfo("processing local JSON log file", "logFilePath", options.LogFilePath)
			ProcessLocalFile(ctx)
		}
	}()

	return ctx
}
