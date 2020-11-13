package config

import (
	"fmt"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/qiangyt/jog/static/grok_vjeantet"
	"github.com/qiangyt/jog/util"
	"github.com/vjeantet/grok"
)

// GrokT ...
type GrokT struct {
	grok          *grok.Grok
	Uses          []string `yaml:"uses"`
	MatchesFields []string `yaml:"matches-fields"`
	LibraryDirs   []string `yaml:"library-dirs"`
}

// DefaultGrokLibraryDir ...
func DefaultGrokLibraryDir() string {
	return JogHomeDir("grok_vjeantet")
}

// Grok ...
type Grok = *GrokT

// SaveDefaultGrokPatternFile ...
func SaveDefaultGrokPatternFile(patternFileName string, patternFileContent string) {
	dir := DefaultGrokLibraryDir()
	util.ReplaceFile(filepath.Join(dir, patternFileName), []byte(patternFileContent))
}

// ResetDefaultGrokLibraryDir ...
func ResetDefaultGrokLibraryDir() {
	dir := DefaultGrokLibraryDir()
	util.RemoveDir(dir)

	InitDefaultGrokLibraryDir()
}

// InitDefaultGrokLibraryDir ...
func InitDefaultGrokLibraryDir() {
	jogHomeDir := JogHomeDir()

	licensePath := filepath.Join(jogHomeDir, "grok_vjeantet.LICENSE")
	util.WriteFileIfNotFound(licensePath, []byte(grok_vjeantet.LICENSE))

	readmePath := filepath.Join(jogHomeDir, "grok_vjeantet.README.md")
	util.WriteFileIfNotFound(readmePath, []byte(grok_vjeantet.README_md))

	dir := DefaultGrokLibraryDir()
	if util.DirExists(dir) == false {
		util.MkdirAll(dir)

		SaveDefaultGrokPatternFile("aws", grok_vjeantet.Aws)
		SaveDefaultGrokPatternFile("bro", grok_vjeantet.Bro)
		SaveDefaultGrokPatternFile("firewalls", grok_vjeantet.Firewalls)
		SaveDefaultGrokPatternFile("haproxy", grok_vjeantet.Haproxy)
		SaveDefaultGrokPatternFile("junos", grok_vjeantet.Junos)
		SaveDefaultGrokPatternFile("linux-syslog", grok_vjeantet.Linux_syslog)
		SaveDefaultGrokPatternFile("mcollective-patterns", grok_vjeantet.Mcollective_patterns)
		SaveDefaultGrokPatternFile("nagios", grok_vjeantet.Nagios)
		SaveDefaultGrokPatternFile("rails", grok_vjeantet.Rails)
		SaveDefaultGrokPatternFile("redis", grok_vjeantet.Redis)
		SaveDefaultGrokPatternFile("bacula", grok_vjeantet.Bacula)
		SaveDefaultGrokPatternFile("exim", grok_vjeantet.Exim)
		SaveDefaultGrokPatternFile("grok-patterns", grok_vjeantet.Grok_patterns)
		SaveDefaultGrokPatternFile("java", grok_vjeantet.Java)
		SaveDefaultGrokPatternFile("mcollective", grok_vjeantet.Mcollective)
		SaveDefaultGrokPatternFile("mongodb", grok_vjeantet.Mongodb)
		SaveDefaultGrokPatternFile("postgresql", grok_vjeantet.Postgresql)
		SaveDefaultGrokPatternFile("ruby", grok_vjeantet.Ruby)
	}

	util.MkdirAll(JogHomeDir("grok_extended"))
	util.MkdirAll(JogHomeDir("grok_mine"))

}

// Init ...
func (i Grok) Init(cfg Configuration) {

	InitDefaultGrokLibraryDir()

	i.grok, _ = grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})

	for _, patternsDir := range i.LibraryDirs {
		dir, err := homedir.Expand(patternsDir)
		if err != nil {
			panic(errors.Wrapf(err, "failed to get home dir: %s", patternsDir))
		}

		if util.DirExists(dir) == false {
			continue
		}
		i.grok.AddPatternsFromPath(dir)
	}
}

// UnmarshalYAML ...
func (i Grok) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return UnmarshalYAML(i, unmarshal)
}

// MarshalYAML ...
func (i Grok) MarshalYAML() (interface{}, error) {
	return MarshalYAML(i)
}

// Reset ...
func (i Grok) Reset() {
	i.Uses = make([]string, 0)
	i.LibraryDirs = []string{}
}

// FromMap ...
func (i Grok) FromMap(m map[string]interface{}) error {
	var err error

	i.LibraryDirs, err = util.ExtractStringSliceFromMap(m, "library-dirs")
	if err != nil {
		return errors.Wrap(err, "failed to parse grok.library-dirs")
	}

	i.Uses, err = util.ExtractStringSliceFromMap(m, "uses")
	if err != nil {
		return errors.Wrap(err, "failed to parse grok.uses")
	}

	i.MatchesFields, err = util.ExtractStringSliceFromMap(m, "matches-fields")
	if err != nil {
		return errors.Wrap(err, "failed to parse grok.matches-fields")
	}
	if len(i.MatchesFields) < 1 {
		return fmt.Errorf("grok.matches-fields must contains at least 1 standard fields")
	}

	// TODO: how to ensure i.Uses doesn't refer to a pattern that not exists ?
	// for _, usedPatternName := range i.Uses {
	// pattern := fmt.Sprintf("%%{%s}", usedPatternName)
	// if _, err := i.grok.Parse(pattern, ""); err != nil {
	// 	return fmt.Errorf("using pattern '%s' but not defined in available patterns", usedPatternName)
	//}
	// }

	return nil
}

// ToMap ...
func (i Grok) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	r["uses"] = i.Uses
	r["library-dirs"] = i.LibraryDirs
	r["matches-fields"] = i.MatchesFields

	return r
}

// Parse ...
func (i Grok) Parse(pattern string, line string) (map[string]string, error) {
	return i.grok.Parse(pattern, line)
}
