package conf

import (
	"fmt"

	"github.com/pkg/errors"
	_io "github.com/qiangyt/jog/pkg/io"
	_util "github.com/qiangyt/jog/pkg/util"
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

// Init ...
func (i Grok) Init() {
	i.grok, _ = grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})

	for _, patternsDir := range i.LibraryDirs {
		dir := _io.ExpandHomePath(patternsDir)
		if _io.DirExists(dir) == false {
			continue
		}
		i.grok.AddPatternsFromPath(dir)
	}
}

// UnmarshalYAML ...
func (i Grok) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return _util.DynObject4YAML(i, unmarshal)
}

// MarshalYAML ...
func (i Grok) MarshalYAML() (interface{}, error) {
	return _util.DynObject2YAML(i)
}

// Reset ...
func (i Grok) Reset() {
	i.Uses = make([]string, 0)
	i.LibraryDirs = []string{}
}

// FromMap ...
func (i Grok) FromMap(m map[string]interface{}) error {
	var err error

	i.LibraryDirs, err = _util.ExtractStringSliceFromMap(m, "library-dirs")
	if err != nil {
		return errors.Wrap(err, "failed to parse grok.library-dirs")
	}

	i.Uses, err = _util.ExtractStringSliceFromMap(m, "uses")
	if err != nil {
		return errors.Wrap(err, "failed to parse grok.uses")
	}

	i.MatchesFields, err = _util.ExtractStringSliceFromMap(m, "matches-fields")
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
