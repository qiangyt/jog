package convert

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/qiangyt/jog/convert/config"
	"github.com/qiangyt/jog/jsonpath"
	"github.com/qiangyt/jog/static"
	"github.com/qiangyt/jog/util"
	"gopkg.in/yaml.v2"
)

// PrintVersion ...
func PrintVersion() {
	fmt.Println(static.AppVersion)
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

// ParseConfigExpression ...
func ParseConfigExpression(expr string) (string, string, error) {
	arr := strings.Split(expr, "=")
	if len(arr) != 2 {
		return "", "", fmt.Errorf("invalid config item expression: <%s>", expr)
	}
	return arr[0], arr[1], nil
}

// ReadConfig ...
func ReadConfig(configFilePath string) config.Configuration {
	if len(configFilePath) == 0 {
		return config.WithDefaultYamlFile()
	}
	return config.WithYamlFile(configFilePath)
}

// PrintConfigItem ...
func PrintConfigItem(m map[string]interface{}, configItemPath string) {
	item, err := jsonpath.Get(m, configItemPath)
	if err != nil {
		panic(errors.Wrap(err, ""))
	}
	out, err := yaml.Marshal(item)
	if err != nil {
		panic(errors.Wrap(err, ""))
	}
	fmt.Print(string(out))
}

// SetConfigItem ...
func SetConfigItem(cfg config.Configuration, m map[string]interface{}, configItemPath string, configItemValue string) {
	if err := jsonpath.Set(m, configItemPath, configItemValue); err != nil {
		panic(errors.Wrap(err, ""))
	}
	if err := cfg.FromMap(m); err != nil {
		panic(errors.Wrap(err, ""))
	}
}

func Main() {
	config.InitDefaultGrokLibraryDir()

	ok, options := OptionsWithCommandLine()
	if !ok {
		return
	}

	if !options.Debug {
		defer func() {
			if p := recover(); p != nil {
				color.Red.Printf("%v\n\n", p)
				os.Exit(1)
				return
			}
		}()
	}

	logFile := util.InitLogger(config.JogHomeDir(true))
	defer logFile.Close()

	cfg := ReadConfig(options.ConfigFilePath)

	if len(options.ConfigItemPath) > 0 {
		m := cfg.ToMap()
		if len(options.ConfigItemValue) > 0 {
			SetConfigItem(cfg, m, options.ConfigItemPath, options.ConfigItemValue)
		} else {
			PrintConfigItem(m, options.ConfigItemPath)
			return
		}
	}

	if cfg.LevelField != nil {
		options.InitLevelFilters(cfg.LevelField.Enums)
	}
	if cfg.TimestampField != nil {
		options.InitTimestampFilters(cfg.TimestampField)
	}

	options.InitGroks(cfg)

	if len(options.LogFilePath) == 0 {
		log.Println("read JSON log lines from stdin")
		ProcessReader(cfg, options, os.Stdin, 1)
	} else {
		log.Printf("processing local JSON log file: %s\n", options.LogFilePath)
		ProcessLocalFile(cfg, options, options.FollowMode, options.LogFilePath)
	}

	fmt.Println()
}
