package main

import (
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/gookit/goutil/strutil"
	"github.com/qiangyt/jog/config"
)

const (
	// AppVersion ...
	AppVersion = "v0.9.13"
)

// PrintVersion ...
func PrintVersion() {
	fmt.Println(AppVersion)
}

// PrintConfigTemplate ...
func PrintConfigTemplate() {
	fmt.Println(config.DefaultYAML)
}

// PrintHelp ...
func PrintHelp() {
	color.New(color.Blue, color.OpBold).Println("\nConvert and view structured (JSON) log")
	PrintVersion()
	fmt.Println()

	color.New(color.FgBlue, color.OpBold).Println("Usage:")
	color.FgBlue.Println("  jog  [option...]  <your JSON log file path>")
	color.FgBlue.Println("    or")
	color.FgBlue.Println("  cat  <your JSON file path>  |  jog  [option...]")
	fmt.Println()

	color.New(color.FgBlue, color.OpBold).Println("Examples:")
	fmt.Println("  1) follow with last 10 lines:         jog -f app-20200701-1.log")
	fmt.Println("  2) follow with specified lines:       jog -n 100 -f app-20200701-1.log")
	fmt.Println("  3) with specified config file:        jog -c another.jog.yml app-20200701-1.log")
	fmt.Println("  4) view docker-compose log:           docker-compose logs | jog")
	fmt.Println("  5) print the default template:        jog -t")
	fmt.Println("  6) with WARN level foreground color set to RED: jog -cs fields.level.enums.WARN.color=FgRed app-20200701-1.log")
	fmt.Println("  7) view the WARN level config item:   jog -cg fields.level.enums.WARN")
	fmt.Println("  8) disable colorization:              jog -cs colorization=false app-20200701-1.log")
	fmt.Println()

	color.New(color.FgBlue, color.OpBold).Println("Options:")
	fmt.Printf("  -c,  --config <config file path>                            Specify config YAML file path. The default is .jog.yaml or $HOME/.jog.yaml \n")
	fmt.Printf("  -cs, --config-set <config item path>=<config item value>    Set value to specified config item \n")
	fmt.Printf("  -cg, --config-get <config item path>                        Get value to specified config item \n")
	fmt.Printf("  -f,  --follow                                               Follow mode - follow log output\n")
	fmt.Printf("  -n,  --lines <number of tail lines>                         Number of tail lines. Available ONLY for follow mode\n")
	fmt.Printf("  -t,  --template                                             Print a config YAML file template\n")
	fmt.Printf("  -h,  --help                                                 Display this information\n")
	fmt.Printf("  -V,  --version                                              Display app version information\n")
	fmt.Printf("  -d,  --debug                                                Print more error detail\n")
	fmt.Println()
}

// CommandLineT ...
type CommandLineT struct {
	LogFilePath     string
	ConfigFilePath  string
	Debug           bool
	ConfigItemPath  string
	ConfigItemValue string
	FollowMode      bool
	NumberOfLines   int
}

// CommandLine ...
type CommandLine = *CommandLineT

// ParseCommandLine ...
func ParseCommandLine() (bool, CommandLine) {

	r := &CommandLineT{
		Debug:         false,
		FollowMode:    false,
		NumberOfLines: -1,
	}
	var err error
	var hasNumberOfLines = false

	for i := 0; i < len(os.Args); i++ {
		if i == 0 {
			continue
		}

		arg := os.Args[i]

		if arg[0:1] == "-" {
			if arg == "-c" || arg == "--config" {
				if i+1 >= len(os.Args) {
					color.Red.Println("Missing config file path\n")
					PrintHelp()
					return false, nil
				}

				r.ConfigFilePath = os.Args[i+1]
				i++
			} else if arg == "-cs" || arg == "--config-set" {
				if i+1 >= len(os.Args) {
					color.Red.Println("Missing config item expression\n")
					PrintHelp()
					return false, nil
				}

				r.ConfigItemPath, r.ConfigItemValue, err = ParseConfigExpression(os.Args[i+1])
				if err != nil {
					color.Red.Println("%v\n", err)
					PrintHelp()
					return false, nil
				}
				i++
			} else if arg == "-cg" || arg == "--config-get" {
				if i+1 >= len(os.Args) {
					color.Red.Println("Missing config item path\n")
					PrintHelp()
					return false, nil
				}

				r.ConfigItemPath = os.Args[i+1]
				i++
			} else if arg == "-f" || arg == "--follow" {
				r.FollowMode = true
			} else if arg == "-n" || arg == "--lines" {
				if i+1 >= len(os.Args) {
					color.Red.Println("Missing lines argument\n")
					PrintHelp()
					return false, nil
				}

				r.NumberOfLines = strutil.MustInt(os.Args[i+1])
				hasNumberOfLines = true
				i++
			} else if arg == "-t" || arg == "--template" {
				PrintConfigTemplate()
				return false, nil
			} else if arg == "-h" || arg == "--help" {
				PrintHelp()
				return false, nil
			} else if arg == "-V" || arg == "--version" {
				PrintVersion()
				return false, nil
			} else if arg == "-d" || arg == "--debug" {
				r.Debug = true
			} else {
				color.Red.Printf("Unknown option: '%s'\n\n", arg)
				PrintHelp()
				return false, nil
			}
		} else {
			r.LogFilePath = arg
		}
	}

	if !hasNumberOfLines {
		if r.FollowMode {
			r.NumberOfLines = 10
		}
	}

	return true, r
}
