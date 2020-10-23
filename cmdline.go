package main

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/qiangyt/jog/static"
)

const (
	// AppVersion ...
	AppVersion = "v0.9.17"
)

// PrintVersion ...
func PrintVersion() {
	fmt.Println(AppVersion)
}

// PrintConfigTemplate ...
func PrintConfigTemplate() {
	fmt.Println(static.DefaultConfiguration_yml)
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
	fmt.Println("  6) only shows WARN & ERROR level:     jog -l warn -l error app-20200701-1.log")
	fmt.Println("  7) with WARN level foreground color set to RED: jog -cs fields.level.enums.WARN.color=FgRed app-20200701-1.log")
	fmt.Println("  8) view the WARN level config item:   jog -cg fields.level.enums.WARN")
	fmt.Println("  9) disable colorization:              jog -cs colorization=false app-20200701-1.log")
	fmt.Println()

	color.New(color.FgBlue, color.OpBold).Println("Options:")
	fmt.Printf("  -c,  --config <config file path>                            Specify config YAML file path. The default is .jog.yaml or $HOME/.jog.yaml \n")
	fmt.Printf("  -cs, --config-set <config item path>=<config item value>    Set value to specified config item \n")
	fmt.Printf("  -cg, --config-get <config item path>                        Get value to specified config item \n")
	fmt.Printf("  -d,  --debug                                                Print more error detail\n")
	fmt.Printf("  -f,  --follow                                               Follow mode - follow log output\n")
	fmt.Printf("  -h,  --help                                                 Display this information\n")
	fmt.Printf("  -l,  --level <level value>                                  Filter by log level. For ex. --level warn \n")
	fmt.Printf("  -n,  --lines <number of tail lines>                         Number of tail lines. 10 by default, for follow mode\n")
	fmt.Printf("  -t,  --template                                             Print a config YAML file template\n")
	fmt.Printf("  -V,  --version                                              Display app version information\n")
	fmt.Println()
}
