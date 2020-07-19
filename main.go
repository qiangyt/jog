package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/qiangyt/jog/config"
	"github.com/qiangyt/jog/jsonpath"
	"github.com/qiangyt/jog/util"
	"gopkg.in/yaml.v2"
)

const (
	// AppVersion ...
	AppVersion = "v0.9.7"
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
	fmt.Println("  1) view a json log:                                               jog app-20200701-1.log")
	fmt.Println("  2) view a json log with specified config file:                    jog -c another.jog.yml app-20200701-1.log")
	fmt.Println("  3) view docker-compose log:                                       docker-compose logs | jog")
	fmt.Println("  4) print the default template:                                    jog -t")
	fmt.Println("  5) view the json log with WARN level foreground color set to RED: jog -cs fields.level.enums.WARN.color=FgRed app-20200701-1.log")
	fmt.Println("  6) view the WARN level config item:                               jog -cg fields.level.enums.WARN")
	fmt.Println()

	color.New(color.FgBlue, color.OpBold).Println("Options:")
	fmt.Printf("  -c,  --config <config file path>                            Specify config YAML file path. The default is .jog.yaml or $HOME/.job.yaml \n")
	fmt.Printf("  -cs, --config-set <config item path>=<config item value>    Set value to specified config item \n")
	fmt.Printf("  -cg, --config-get <config item path>                        Get value to specified config item \n")
	fmt.Printf("  -t,  --template                                             Print a config YAML file template\n")
	fmt.Printf("  -h,  --help                                                 Display this information\n")
	fmt.Printf("  -V,  --version                                              Display app version information\n")
	fmt.Printf("  -d,  --debug                                                Print more error detail\n")
	fmt.Println()
}

// ParseConfigExpression ...
func ParseConfigExpression(expr string) (string, string, error) {
	arr := strings.Split(expr, "=")
	if len(arr) != 2 {
		return "", "", fmt.Errorf("invalid config item expression: <%s>", expr)
	}
	return arr[0], arr[1], nil
}

// ReadConfig ...
func ReadConfig(configFilePath string) Config {
	if len(configFilePath) == 0 {
		return ConfigWithDefaultYamlFile()
	}
	return ConfigWithYamlFile(configFilePath)
}

func main() {

	var configFilePath string
	var logFilePath string
	var debug bool
	var err error
	var configItemPath, configItemValue string

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
					return
				}

				if i+1 < len(os.Args) {
					configFilePath = os.Args[i+1]
				}
				i++
			} else if arg == "-cs" || arg == "--config-set" {
				if i+1 >= len(os.Args) {
					color.Red.Println("Missing config item expression\n")
					PrintHelp()
					return
				}

				if i+1 < len(os.Args) {
					configItemPath, configItemValue, err = ParseConfigExpression(os.Args[i+1])
					if err != nil {
						color.Red.Println("%v\n", err)
						PrintHelp()
						return
					}
				}
				i++
			} else if arg == "-cg" || arg == "--config-get" {
				if i+1 >= len(os.Args) {
					color.Red.Println("Missing config item path\n")
					PrintHelp()
					return
				}

				if i+1 < len(os.Args) {
					configItemPath = os.Args[i+1]
				}
				i++
			} else if arg == "-t" || arg == "--template" {
				PrintConfigTemplate()
				return
			} else if arg == "-h" || arg == "--help" {
				PrintHelp()
				return
			} else if arg == "-V" || arg == "--version" {
				PrintVersion()
				return
			} else if arg == "-d" || arg == "--debug" {
				debug = true
			} else {
				color.Red.Printf("Unknown option: '%s'\n\n", arg)
				PrintHelp()
				return
			}
		} else {
			logFilePath = arg
		}
	}

	if !debug {
		defer func() {
			if p := recover(); p != nil {
				color.Red.Printf("%v\n\n", p)
				os.Exit(1)
				return
			}
		}()
	}

	logFile := util.InitLogger()
	defer logFile.Close()

	cfg := ReadConfig(configFilePath)

	if len(configItemPath) > 0 {
		m := cfg.ToMap()
		if len(configItemValue) > 0 {
			if err := jsonpath.Set(m, configItemPath, configItemValue); err != nil {
				panic(errors.Wrap(err, ""))
			}
			if err := cfg.FromMap(m); err != nil {
				panic(errors.Wrap(err, ""))
			}
		} else {
			item, err := jsonpath.Get(m, configItemPath)
			if err != nil {
				panic(errors.Wrap(err, ""))
			}
			out, err := yaml.Marshal(item)
			if err != nil {
				panic(errors.Wrap(err, ""))
			}
			fmt.Print(string(out))
			return
		}
	}

	if len(logFilePath) == 0 {
		log.Println("read JSON log lines from stdin")
		ProcessReader(cfg, os.Stdin)
	} else {
		log.Printf("processing local JSON log file: %s\n", logFilePath)
		ProcessLocalFile(cfg, logFilePath)
	}

	fmt.Println()
}
