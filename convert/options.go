package convert

import (
	"fmt"
	"time"

	"github.com/gookit/goutil/strutil"
	"github.com/qiangyt/jog/util"
	"github.com/tj/go-naturaldate"
)

// OptionsT ...
type OptionsT struct {
	LogFilePath    string
	ConfigFilePath string

	ConfigItemPath  string
	ConfigItemValue string
	FollowMode      bool
	NumberOfLines   int

	levelFilterTexts []string
	levelFilters     []Enum

	beforeFilterText string
	BeforeFilter     *time.Time

	afterFilterText string
	AfterFilter     *time.Time

	OutputRawJSON bool

	GrokPatternsUsed []string
	GrokPatterns     []string

	OpenWebGUI bool
}

// Options ...
type Options = *OptionsT

// PrintConfigTemplate ...
func (i Options) PrintConfigTemplate() {
	fmt.Println(BuildDefaultConfigYAML())
}

// InitGroks ...
func (i Options) InitGroks(cfg Config) {
	if len(i.GrokPatternsUsed) == 0 {
		// try to uses default patterns
		i.GrokPatternsUsed = cfg.Grok.Uses
	}

	i.GrokPatterns = make([]string, len(i.GrokPatternsUsed))
	for index, patternName := range i.GrokPatternsUsed {
		i.GrokPatterns[index] = "%{" + patternName + "}"
	}
}

func (i Options) isGrokEnabled() bool {
	return len(i.GrokPatterns) > 0
}

// GetLevelFilters ...
func (i Options) GetLevelFilters() []Enum {
	return i.levelFilters
}

// InitLevelFilters ...
func (i Options) InitLevelFilters(levelFieldEnums EnumMap) {
	if len(i.levelFilterTexts) == 0 {
		i.levelFilters = make([]Enum, 0)
		return
	}

	for _, levelFilterText := range i.levelFilterTexts {
		levelFilter := levelFieldEnums.GetEnum(levelFilterText)
		i.levelFilters = append(i.levelFilters, levelFilter)
	}
}

// InitTimestampFilters ...
func (i Options) InitTimestampFilters(ctx util.JogContext, timestampField Field) {
	now := time.Now()

	if len(i.beforeFilterText) > 0 {
		f, err := naturaldate.Parse(i.beforeFilterText, now, naturaldate.WithDirection(naturaldate.Past))
		if err != nil {
			ctx.LogWarn("failed to parse before-time filter as natural timestamp, so try absolute parse", "beforeFilter", i.beforeFilterText)
			f = timestampField.ParseTimestamp(i.beforeFilterText)
		}
		ctx.LogInfo("before-time filter", "beforeFlter", f)
		i.BeforeFilter = &f
	}
	if len(i.afterFilterText) > 0 {
		f, err := naturaldate.Parse(i.afterFilterText, now, naturaldate.WithDirection(naturaldate.Past))
		if err != nil {
			ctx.LogWarn("failed to parse after-time filter as natural timestamp, so try absolute parse", "afterFilter", i.afterFilterText)
			f = timestampField.ParseTimestamp(i.afterFilterText)
		}
		ctx.LogInfo("after-time filter", "afterFilter", f)
		i.AfterFilter = &f

		if i.BeforeFilter != nil {
			if i.BeforeFilter.Before(*i.AfterFilter) {
				panic(fmt.Errorf("before-time filter (%s) shouldn't be before after-time filter (%s)", i.beforeFilterText, i.afterFilterText))
			}
		}
	}
}

// HasTimestampFilter ...
func (i Options) HasTimestampFilter() bool {
	return i.BeforeFilter != nil || i.AfterFilter != nil
}

// NewOptionsWithCommandLine ...
func NewOptionsWithCommandLine(args []string) (bool, Options) {

	r := &OptionsT{
		FollowMode:       false,
		NumberOfLines:    -1,
		levelFilterTexts: make([]string, 0),
		GrokPatternsUsed: make([]string, 0),
		OpenWebGUI:       false,
	}
	var err error
	var hasNumberOfLines = false

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg[0:1] == "-" {
			if arg == "-c" || arg == "--config" {
				if i+1 >= len(args) {
					util.PrintErrorHint("Missing config file path")
					return false, nil
				}

				r.ConfigFilePath = args[i+1]
				i++
			} else if arg == "-cs" || arg == "--config-set" {
				if i+1 >= len(args) {
					util.PrintErrorHint("Missing config item expression")
					return false, nil
				}

				r.ConfigItemPath, r.ConfigItemValue, err = util.ParseConfigExpression(args[i+1])
				if err != nil {
					util.PrintErrorHint("%v", err)
					return false, nil
				}
				i++
			} else if arg == "-cg" || arg == "--config-get" {
				if i+1 >= len(args) {
					util.PrintErrorHint("Missing config item path")
					return false, nil
				}

				r.ConfigItemPath = args[i+1]
				i++
			} else if arg == "-f" || arg == "--follow" {
				r.FollowMode = true
			} else if arg == "-n" || arg == "--lines" {
				if i+1 >= len(args) {
					util.PrintErrorHint("Missing lines argument")
					return false, nil
				}

				r.NumberOfLines = strutil.MustInt(args[i+1])
				hasNumberOfLines = true
				i++
			} else if arg == "-t" || arg == "--template" {
				r.PrintConfigTemplate()
				return false, nil
			} else if arg == "-j" || arg == "--json" {
				r.OutputRawJSON = true
			} else if arg == "-l" || arg == "--level" {
				if i+1 >= len(args) {
					util.PrintErrorHint("Missing level argument")
					return false, nil
				}

				r.levelFilterTexts = append(r.levelFilterTexts, args[i+1])
				i++
			} else if arg == "-g" || arg == "--grok" {
				if i+1 >= len(args) {
					util.PrintErrorHint("Missing grok argument")
					return false, nil
				}

				r.GrokPatternsUsed = append(r.GrokPatternsUsed, args[i+1])
				i++
			} else if arg == "--reset-grok-library-dir" {
				util.ResetDefaultGrokLibraryDir()
				return false, nil
			} else if arg == "-a" || arg == "--after" {
				if i+1 >= len(args) {
					util.PrintErrorHint("Missing after argument")
					return false, nil
				}

				r.afterFilterText = args[i+1]
				i++
			} else if arg == "-b" || arg == "--before" {
				if i+1 >= len(args) {
					util.PrintErrorHint("Missing before argument")
					return false, nil
				}

				r.beforeFilterText = args[i+1]
				i++
			} else if arg == "-w" || arg == "--web-gui" {
				r.OpenWebGUI = true
			} else {
				util.PrintErrorHint("Unknown option: '%s'", arg)
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
