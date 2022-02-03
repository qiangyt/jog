package convert

import (
	"os"

	"github.com/gookit/color"
	"github.com/qiangyt/jog/common"
	"github.com/qiangyt/jog/util"
)

func Main(done chan bool, globalOptions common.GlobalOptions) ConvertContext {
	util.InitDefaultGrokLibraryDir()

	ok, options := NewOptionsWithCommandLine(globalOptions.SubArgs())
	if !ok {
		close(done)
		return nil
	}

	ctx := NewConvertContext(options, util.JogHomeDir(true), globalOptions.Version())
	defer ctx.Close()

	if !ctx.LoadConfig() {
		close(done)
		return ctx
	}

	go func() {
		defer close(done)

		//TODO: trap CTRL+C signal
		if !globalOptions.Debug() {
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
