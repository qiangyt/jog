package config

import (
	"strings"

	"github.com/gookit/goutil/strutil"
	"github.com/qiangyt/jog/util"
)

// ElementT ...
type ElementT struct {
	Color  util.Color
	Print  bool
	Before string
	After  string
}

// Element ...
type Element = *ElementT

// UnmarshalYAML ...
func (i Element) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return util.UnmarshalYAML(i, unmarshal)
}

// MarshalYAML ...
func (i Element) MarshalYAML() (interface{}, error) {
	return util.MarshalYAML(i)
}

// FromMap ...
func (i Element) FromMap(m map[string]interface{}) error {
	colorV := util.ExtractFromMap(m, "color")
	if colorV != nil {
		if err := util.UnmashalYAMLAgain(colorV, &i.Color); err != nil {
			return err
		}
	}

	printV := util.ExtractFromMap(m, "print")
	if printV != nil {
		i.Print = util.ToBool(printV)
	}

	beforeV := util.ExtractFromMap(m, "before")
	if beforeV != nil {
		i.Before = strutil.MustString(beforeV)
	}

	afterV := util.ExtractFromMap(m, "after")
	if afterV != nil {
		i.After = strutil.MustString(afterV)
	}

	return nil
}

// ToMap ...
func (i Element) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	r["color"] = i.Color.String()
	r["print"] = i.Print
	r["before"] = i.Before
	r["after"] = i.After
	return r
}

// Reset ...
func (i Element) Reset() {
	i.Color = &util.ColorT{}
	i.Color.Set("OpReset")

	i.Print = true
	i.Before = ""
	i.After = ""
}

// GetColor ...
func (i Element) GetColor(value string) util.Color {
	return i.Color
}

// IsEnabled ...
func (i Element) IsEnabled() bool {
	return i.Print
}

// PrintBody ...
func (i Element) PrintBody(color util.Color, builder *strings.Builder, body string) {
	if color == nil {
		builder.WriteString(body)
	} else {
		builder.WriteString(color.Sprint(body))
	}
}

// PrintBefore ...
func (i Element) PrintBefore(color util.Color, builder *strings.Builder) {
	if len(i.Before) == 0 {
		return
	}
	if color == nil {
		builder.WriteString(i.Before)
	} else {
		builder.WriteString(color.Sprint(i.Before))
	}
}

// PrintAfter ...
func (i Element) PrintAfter(color util.Color, builder *strings.Builder) {
	if len(i.After) == 0 {
		return
	}
	if color == nil {
		builder.WriteString(i.After)
	} else {
		builder.WriteString(color.Sprint(i.After))
	}
}
