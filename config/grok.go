package config

import (
	"github.com/qiangyt/jog/util"
	"github.com/vjeantet/grok"
)

// GrokT ...
type GrokT struct {
	grok     *grok.Grok
	Uses     []string          `yaml:"uses"`
	Patterns map[string]string `yaml:"patterns"`
}

// Grok ...
type Grok = *GrokT

// Init ...
func (i Grok) Init(cfg Configuration) {
	i.grok, _ = grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})
	//i.grokPatternsUsed []string
	//i.grokConfig.AddPattern("", "")
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
	i.Patterns = make(map[string]string)
}

// FromMap ...
func (i Grok) FromMap(m map[string]interface{}) error {
	i.Uses = util.ExtractFromMap(m, "uses").([]string)
	i.Patterns = util.ExtractFromMap(m, "patterns").(map[string]string)
	return nil
}

// ToMap ...
func (i Grok) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	r["uses"] = i.Uses

	patterns := make(map[string]string)
	for patternName, patternExpr := range i.Patterns {
		patterns[patternName] = patternExpr
	}
	r["patterns"] = patterns

	return r
}
