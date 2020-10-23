package main

import (
	"os"

	"github.com/gookit/color"
	"github.com/gookit/goutil/strutil"
	"github.com/qiangyt/jog/config"
)

// OptionsT ...
type OptionsT struct {
	LogFilePath      string
	ConfigFilePath   string
	Debug            bool
	ConfigItemPath   string
	ConfigItemValue  string
	FollowMode       bool
	NumberOfLines    int
	levelFilterTexts []string
	levelFilters     []config.Enum
}

// Options ...
type Options = *OptionsT

// GetLevelFilters ...
func (i Options) GetLevelFilters() []config.Enum {
	return i.levelFilters
}

// InitLevelFilters ...
func (i Options) InitLevelFilters(levelFieldEnums config.EnumMap) {
	if len(i.levelFilterTexts) == 0 {
		i.levelFilters = make([]config.Enum, 0)
		return
	}

	for _, levelFilterText := range i.levelFilterTexts {
		levelFilter := levelFieldEnums.GetEnum(levelFilterText)
		i.levelFilters = append(i.levelFilters, levelFilter)
	}
}

// OptionsWithCommandLine ...
func OptionsWithCommandLine() (bool, Options) {

	r := &OptionsT{
		Debug:            false,
		FollowMode:       false,
		NumberOfLines:    -1,
		levelFilterTexts: make([]string, 0),
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
			} else if arg == "-l" || arg == "--level" {
				if i+1 >= len(os.Args) {
					color.Red.Println("Missing level argument\n")
					PrintHelp()
					return false, nil
				}

				r.levelFilterTexts = append(r.levelFilterTexts, os.Args[i+1])
				i++
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
