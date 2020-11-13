package config

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/qiangyt/jog/static/grok_patterns"
	"github.com/qiangyt/jog/util"
	"github.com/vjeantet/grok"
)

// GrokT ...
type GrokT struct {
	grok        *grok.Grok
	Uses        []string `yaml:"uses"`
	LibraryDirs []string `yaml:"library-dirs"`
}

// DefaultGrokLibraryDir ...
func DefaultGrokLibraryDir() string {
	return JogHomeDir("grok-library")
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

	licensePath := filepath.Join(jogHomeDir, "vjeantet-grok.LICENSE")
	util.WriteFileIfNotFound(licensePath, []byte(grok_patterns.LICENSE))

	readmePath := filepath.Join(jogHomeDir, "vjeantet-grok.README.md")
	util.WriteFileIfNotFound(readmePath, []byte(grok_patterns.README_md))

	dir := DefaultGrokLibraryDir()
	if util.DirExists(dir) {
		return
	}
	util.MkdirAll(dir)

	SaveDefaultGrokPatternFile("aws", grok_patterns.Aws)
	SaveDefaultGrokPatternFile("bro", grok_patterns.Bro)
	SaveDefaultGrokPatternFile("firewalls", grok_patterns.Firewalls)
	SaveDefaultGrokPatternFile("haproxy", grok_patterns.Haproxy)
	SaveDefaultGrokPatternFile("junos", grok_patterns.Junos)
	SaveDefaultGrokPatternFile("linux-syslog", grok_patterns.Linux_syslog)
	SaveDefaultGrokPatternFile("mcollective-patterns", grok_patterns.Mcollective_patterns)
	SaveDefaultGrokPatternFile("nagios", grok_patterns.Nagios)
	SaveDefaultGrokPatternFile("rails", grok_patterns.Rails)
	SaveDefaultGrokPatternFile("redis", grok_patterns.Redis)
	SaveDefaultGrokPatternFile("bacula", grok_patterns.Bacula)
	SaveDefaultGrokPatternFile("exim", grok_patterns.Exim)
	SaveDefaultGrokPatternFile("grok-patterns", grok_patterns.Grok_patterns)
	SaveDefaultGrokPatternFile("java", grok_patterns.Java)
	SaveDefaultGrokPatternFile("mcollective", grok_patterns.Mcollective)
	SaveDefaultGrokPatternFile("mongodb", grok_patterns.Mongodb)
	SaveDefaultGrokPatternFile("postgresql", grok_patterns.Postgresql)
	SaveDefaultGrokPatternFile("ruby", grok_patterns.Ruby)
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

	return r
}
