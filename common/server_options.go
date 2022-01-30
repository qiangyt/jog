package common

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/pkg/errors"
	"github.com/qiangyt/jog/static"
	"github.com/qiangyt/jog/util"
)

// ServerOptionsT ...
type ServerOptionsT struct {
	LogFilePath    string
	ConfigFilePath string

	ConfigItemPath  string
	ConfigItemValue string
}

// ServerOptions ...
type ServerOptions = *ServerOptionsT

// BuildDefaultServerConfigurationYAML ...
func BuildDefaultServerConfigurationYAML() string {
	tmpl, err := template.New("default server configuration YAML").Parse(static.DefaultServer_yml)
	if err != nil {
		panic(errors.Wrap(err, "failed to parse default server configuration YAML as template"))
	}

	vars := map[string]string{}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, vars)
	if err != nil {
		panic(errors.Wrap(err, "failed to execute default configuration YAML as template"))
	}
	return buf.String()
}

// PrintConfigTemplate ...
func (i ServerOptions) PrintConfigTemplate() {
	fmt.Println(BuildDefaultServerConfigurationYAML())
}

// ServerOptionsWithCommandLine ...
func ServerOptionsWithCommandLine(globalOptions GlobalOptions) (bool, ServerOptions) {

	r := &ServerOptionsT{}
	var err error

	args := globalOptions.SubArgs

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg[0:1] == "-" {
			if arg == "-c" || arg == "--config" {
				if i+1 >= len(args) {
					globalOptions.PrintErrorHint("Missing config file path")
					return false, nil
				}

				r.ConfigFilePath = args[i+1]
				i++
			} else if arg == "-cs" || arg == "--config-set" {
				if i+1 >= len(args) {
					globalOptions.PrintErrorHint("Missing config item expression")
					return false, nil
				}

				r.ConfigItemPath, r.ConfigItemValue, err = util.ParseConfigExpression(args[i+1])
				if err != nil {
					globalOptions.PrintErrorHint("%v", err)
					return false, nil
				}
				i++
			} else if arg == "-cg" || arg == "--config-get" {
				if i+1 >= len(args) {
					globalOptions.PrintErrorHint("Missing config item path")
					return false, nil
				}

				r.ConfigItemPath = args[i+1]
				i++
			} else if arg == "-t" || arg == "--template" {
				r.PrintConfigTemplate()
				return false, nil
			} else {
				globalOptions.PrintErrorHint("Unknown option: '%s'", arg)
				return false, nil
			}
		} else {
			r.LogFilePath = arg
		}
	}

	return true, r
}
