package common

import (
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/qiangyt/jog/convert/config"
)

type RunMode int

const (
	RunMode_Client = iota
	RunMode_Server
	RunMode_Default = RunMode_Client
)

type GlobalOptionsT struct {
	Debug   bool
	RunMode RunMode
	SubArgs []string
	Version string
}

type GlobalOptions = *GlobalOptionsT

func GlobalOptionsWithCommandLine(version string) (bool, GlobalOptions) {

	r := &GlobalOptionsT{
		Debug:   false,
		RunMode: RunMode_Default,
		SubArgs: []string{},
		Version: version,
	}

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		if i == 1 && arg == "server" {
			r.RunMode = RunMode_Server
		} else if arg[0:1] != "-" {
			r.SubArgs = append(r.SubArgs, arg)
		} else {
			if arg == "-h" || arg == "--help" {
				r.PrintHelp()
				return false, nil
			} else if arg == "-V" || arg == "--version" {
				r.PrintVersion()
				return false, nil
			} else if arg == "-d" || arg == "--debug" {
				r.Debug = true
			} else {
				r.SubArgs = append(r.SubArgs, arg)
			}
		}
	}

	return true, r
}

// PrintVersion ...
func (i GlobalOptions) PrintVersion() {
	fmt.Println(i.Version)
}

func (i GlobalOptions) PrintErrorHint(format string, a ...interface{}) {
	color.Red.Printf(format+". Please check above example\n", a...)
}

// PrintHelp ...
func (i GlobalOptions) PrintHelp() {
	color.New(color.Blue, color.OpBold).Println("\nJog: convert/view/share log")
	i.PrintVersion()
	fmt.Println()

	color.New(color.FgBlue, color.OpBold).Println("Global options:")
	fmt.Printf("  -d,  --debug                Print more error detail\n")
	fmt.Printf("  -h,  --help                 Display this information\n")
	color.New(color.FgGreen, color.OpBold).Printf("  server -h,  server --help   Display server mode help information\n")
	fmt.Printf("  -V,  --version              Display app version information\n")
	fmt.Println()

	if i.RunMode == RunMode_Server {
		i.PrintServerHelp()
	} else {
		i.PrintConvertHelp()
	}
}

func (i GlobalOptions) PrintServerHelp() {
	fmt.Println("server help: TODO")
	fmt.Println()

	color.New(color.FgBlue, color.OpBold).Println("Server mode usage:")
	color.FgBlue.Println("  jog  server [option...]")
	fmt.Println()

	color.New(color.FgBlue, color.OpBold).Println("Servr mode options:")
	fmt.Printf("  -c,  --config <server config file path>                     Specify server config YAML file path. The default is .jog.server.yaml \n")
	fmt.Printf("  -cs, --config-set <config item path>=<config item value>    Set value to specified config item \n")
	fmt.Printf("  -cg, --config-get <config item path>                        Get value to specified config item \n")
	fmt.Printf("  -t,  --template                                             Print a server config YAML file template\n")
	fmt.Println()
}

func (i GlobalOptions) PrintConvertHelp() {
	defaultGrokLibraryDirs := strings.Join(config.DefaultGrokLibraryDirs(false), ", ")

	color.New(color.FgBlue, color.OpBold).Println("Convert/view usage:")
	color.FgBlue.Println("  jog  [option...]  <your JSON log file path>")
	color.FgBlue.Println("    or")
	color.FgBlue.Println("  cat  <your JSON file path>  |  jog  [option...]")
	fmt.Println()

	color.New(color.FgBlue, color.OpBold).Println("Convert/view examples:")
	fmt.Println("   1) follow with last 10 lines:         jog -f app-20200701-1.log")
	fmt.Println("   2) follow with specified lines:       jog -n 100 -f app-20200701-1.log")
	fmt.Println("   3) with specified config file:        jog -c another.jog.yml app-20200701-1.log")
	fmt.Println("   4) view docker-compose log:           docker-compose logs | jog")
	fmt.Println("   5) print the default template:        jog -t")
	fmt.Println("   6) only shows WARN & ERROR level:     jog -l warn -l error app-20200701-1.log")
	fmt.Println("   7) shows with timestamp range:        jog --after 2020-7-1 --before 2020-7-3 app-20200701-1.log")
	fmt.Println("   8) natural timestamp range:           jog --after \"1 week\" --before \"2 days\" app-20200701-1.log")
	fmt.Println("   9) output raw JSON and apply time range filter:      jog --after \"1 week\" --before \"2 days\" app-20200701-1.log --json")
	fmt.Println("   10) disable colorization:             jog -cs colorization=false app-20200701-1.log")
	fmt.Println("   11) view apache log, non-JSON log     jog -g COMMONAPACHELOG example_logs/grok_apache.log")
	fmt.Println()

	color.New(color.FgBlue, color.OpBold).Println("Convert/view options:")
	fmt.Printf("  -a,  --after <timestamp>                                    'after' time filter. Auto-detect the timestamp format; can be natural datetime \n")
	fmt.Printf("  -b,  --before <timestamp>                                   'before' time filter. Auto-detect the timestamp format; can be natural datetime \n")
	fmt.Printf("  -c,  --config <convertion config file path>                 Specify convertion config YAML file path. The default is .jog.yaml or $HOME/.jog.yaml \n")
	fmt.Printf("  -cs, --config-set <config item path>=<config item value>    Set value to specified config item \n")
	fmt.Printf("  -cg, --config-get <config item path>                        Get value to specified config item \n")
	fmt.Printf("  -f,  --follow                                               Follow mode - follow log output\n")
	fmt.Printf("  -g,  --grok <grok pattern name>                             For non-json log line. The default patterns are saved in [%s]\n", defaultGrokLibraryDirs)
	fmt.Printf("  -j,  --json                                                 Output the raw JSON but then able to apply filters\n")
	fmt.Printf("  -l,  --level <level value>                                  Filter by log level. For ex. --level warn \n")
	fmt.Printf("  -n,  --lines <number of tail lines>                         Number of tail lines. 10 by default, for follow mode\n")
	fmt.Printf("       --reset-grok-library-dir                               Save default GROK patterns to [%s]\n", defaultGrokLibraryDirs)
	fmt.Printf("  -t,  --template                                             Print a convertion config YAML file template\n")
	fmt.Println()
}
