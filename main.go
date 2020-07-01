package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gookit/color"
)

const (
	// AppVersion ...
	AppVersion = "v0.9.0"
)

// PrintVersion ...
func PrintVersion() {
	fmt.Println(AppVersion)
}

// PrintConfigTemplate ...
func PrintConfigTemplate() {
	fmt.Println(ConfigDefaultYAML)
}

// PrintHelp ...
func PrintHelp() {
	color.New(color.Blue, color.OpBold).Println("Convert and view structured (JSON) log")
	PrintVersion()
	fmt.Println()

	color.OpBold.Println("Usage:")
	fmt.Println("  jog  [option...]  <your JSON log file path>")
	fmt.Println("  cat  <your JSON file path>  |  jog  [option...]")
	fmt.Println()

	color.OpBold.Println("Options:")
	fmt.Printf("  -c, --config <config file path>  Specify config YAML file path. The default is .jog.yaml or $HOME/.job.yaml \n")
	fmt.Printf("  -t, --template                   Print a config YAML file template\n")
	fmt.Printf("  -h, --help                       Display this information\n")
	fmt.Printf("  -V, --version                    Display app version information\n")
	fmt.Printf("  -d, --debug                      Print more error detail\n")
	fmt.Println()
}

func main() {

	var configFilePath string
	var logFilePath string
	var debug bool
	// logFilePath = "./example_logs/logstash.log"

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

	logFile := InitLogger()
	defer logFile.Close()

	var cfg Config
	if len(configFilePath) == 0 {
		cfg = ConfigWithDefaultYamlFile()
	} else {
		cfg = ConfigWithYamlFile(configFilePath)
	}

	if len(logFilePath) == 0 {
		log.Println("Read JSON log lines from stdin")
		ProcessLinesWithReader(cfg, os.Stdin)
	} else {
		log.Printf("processing local JSON log file: %s\n", logFilePath)
		ProcessLinesWithLocalFile(cfg, logFilePath)
	}

	fmt.Println()
}
