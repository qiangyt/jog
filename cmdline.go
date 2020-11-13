package main

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
	"github.com/qiangyt/jog/config"
	"github.com/qiangyt/jog/static"
)

// PrintVersion ...
func PrintVersion() {
	fmt.Println(static.AppVersion)
}

// PrintConfigTemplate ...
func PrintConfigTemplate() {
	fmt.Println(config.BuildDefaultConfigurationYAML())
}

// PrintHelp ...
func PrintHelp() {
	defaultGrokLibraryDirs := strings.Join(config.DefaultGrokLibraryDirs(false), ", ")

	color.New(color.Blue, color.OpBold).Println("\nConvert and view structured (JSON) log")
	PrintVersion()
	fmt.Println()

	color.New(color.FgBlue, color.OpBold).Println("Usage:")
	color.FgBlue.Println("  jog  [option...]  <your JSON log file path>")
	color.FgBlue.Println("    or")
	color.FgBlue.Println("  cat  <your JSON file path>  |  jog  [option...]")
	fmt.Println()

	color.New(color.FgBlue, color.OpBold).Println("Examples:")
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

	color.New(color.FgBlue, color.OpBold).Println("Options:")
	fmt.Printf("  -a,  --after <timestamp>                                    'after' time filter. Auto-detect the timestamp format; can be natural datetime \n")
	fmt.Printf("  -b,  --before <timestamp>                                   'before' time filter. Auto-detect the timestamp format; can be natural datetime \n")
	fmt.Printf("  -c,  --config <config file path>                            Specify config YAML file path. The default is .jog.yaml or $HOME/.jog.yaml \n")
	fmt.Printf("  -cs, --config-set <config item path>=<config item value>    Set value to specified config item \n")
	fmt.Printf("  -cg, --config-get <config item path>                        Get value to specified config item \n")
	fmt.Printf("  -d,  --debug                                                Print more error detail\n")
	fmt.Printf("  -f,  --follow                                               Follow mode - follow log output\n")
	fmt.Printf("  -g,  --grok <grok pattern name>                             For non-json log line. The default patterns are saved in [%s]\n", defaultGrokLibraryDirs)
	fmt.Printf("  -h,  --help                                                 Display this information\n")
	fmt.Printf("  -j,  --json                                                 Output the raw JSON but then able to apply filters\n")
	fmt.Printf("  -l,  --level <level value>                                  Filter by log level. For ex. --level warn \n")
	fmt.Printf("  -n,  --lines <number of tail lines>                         Number of tail lines. 10 by default, for follow mode\n")
	fmt.Printf("       --reset-grok-library-dir                               Save default GROK patterns to [%s]\n", defaultGrokLibraryDirs)
	fmt.Printf("  -t,  --template                                             Print a config YAML file template\n")
	fmt.Printf("  -V,  --version                                              Display app version information\n")
	fmt.Println()
}
