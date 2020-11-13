package config

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/qiangyt/jog/static/grok_extended"
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

// Grok ...
type Grok = *GrokT

// SaveGrokPatternFile ...
func SaveGrokPatternFile(dir string, patternFileName string, patternFileContent string) {
	util.ReplaceFile(filepath.Join(dir, patternFileName), []byte(patternFileContent))
}

// DefaultGrokLibraryDirs ...
func DefaultGrokLibraryDirs(expand bool) []string {
	return []string{
		JogHomeDir(expand, "grok_vjeantet"),
		JogHomeDir(expand, "grok_extended"),
	}
}

// ResetDefaultGrokLibraryDir ...
func ResetDefaultGrokLibraryDir() {
	dirVjeantet := JogHomeDir(true, "grok_vjeantet")
	util.RemoveDir(dirVjeantet)

	dirExtended := JogHomeDir(true, "grok_extended")
	util.RemoveDir(dirExtended)

	InitDefaultGrokLibraryDir()
}

// InitDefaultGrokLibraryDir ...
func InitDefaultGrokLibraryDir() {
	jogHomeDir := JogHomeDir(true)

	licensePath := filepath.Join(jogHomeDir, "grok_vjeantet.LICENSE")
	util.WriteFileIfNotFound(licensePath, []byte(grok_vjeantet.LICENSE))

	readmePath := filepath.Join(jogHomeDir, "grok_vjeantet.README.md")
	util.WriteFileIfNotFound(readmePath, []byte(grok_vjeantet.README_md))

	dirVjeantet := JogHomeDir(true, "grok_vjeantet")
	if util.DirExists(dirVjeantet) == false {
		util.MkdirAll(dirVjeantet)

		SaveGrokPatternFile(dirVjeantet, "aws", grok_vjeantet.Aws)
		SaveGrokPatternFile(dirVjeantet, "bro", grok_vjeantet.Bro)
		SaveGrokPatternFile(dirVjeantet, "firewalls", grok_vjeantet.Firewalls)
		SaveGrokPatternFile(dirVjeantet, "haproxy", grok_vjeantet.Haproxy)
		SaveGrokPatternFile(dirVjeantet, "junos", grok_vjeantet.Junos)
		SaveGrokPatternFile(dirVjeantet, "linux-syslog", grok_vjeantet.Linux_syslog)
		SaveGrokPatternFile(dirVjeantet, "mcollective-patterns", grok_vjeantet.Mcollective_patterns)
		SaveGrokPatternFile(dirVjeantet, "nagios", grok_vjeantet.Nagios)
		SaveGrokPatternFile(dirVjeantet, "rails", grok_vjeantet.Rails)
		SaveGrokPatternFile(dirVjeantet, "redis", grok_vjeantet.Redis)
		SaveGrokPatternFile(dirVjeantet, "bacula", grok_vjeantet.Bacula)
		SaveGrokPatternFile(dirVjeantet, "exim", grok_vjeantet.Exim)
		SaveGrokPatternFile(dirVjeantet, "grok-patterns", grok_vjeantet.Grok_patterns)
		SaveGrokPatternFile(dirVjeantet, "java", grok_vjeantet.Java)
		SaveGrokPatternFile(dirVjeantet, "mcollective", grok_vjeantet.Mcollective)
		SaveGrokPatternFile(dirVjeantet, "mongodb", grok_vjeantet.Mongodb)
		SaveGrokPatternFile(dirVjeantet, "postgresql", grok_vjeantet.Postgresql)
		SaveGrokPatternFile(dirVjeantet, "ruby", grok_vjeantet.Ruby)
	}

	dirExtended := JogHomeDir(true, "grok_extended")
	if util.DirExists(dirExtended) == false {
		util.MkdirAll(dirExtended)

		SaveGrokPatternFile(dirExtended, "pm2", grok_extended.Pm2)
	}

	util.MkdirAll(JogHomeDir(true, "grok_mine"))

}

// Init ...
func (i Grok) Init(cfg Configuration) {

	InitDefaultGrokLibraryDir()

	i.grok, _ = grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})

	for _, patternsDir := range i.LibraryDirs {
		dir := util.ExpandPath(patternsDir)
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
