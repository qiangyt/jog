package convert

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gookit/goutil/strutil"
	"github.com/qiangyt/jog/util"
)

// ElementT ...
type ElementT struct {
	Color       util.Color
	Print       bool
	PrintFormat string `yaml:"print-format"`
}

// Element ...
type Element = *ElementT

// UnmarshalYAML ...
func (i Element) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return DynObject4YAML(i, unmarshal)
}

// MarshalYAML ...
func (i Element) MarshalYAML() (interface{}, error) {
	return DynObject2YAML(i)
}

// Init ...
func (i Element) Init(cfg Config) {

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

	printFormatV := util.ExtractFromMap(m, "print-format")
	if printFormatV != nil {
		printFormatT := strutil.MustString(printFormatV)
		if validPrintFormat(printFormatT) {
			i.PrintFormat = printFormatT
		} else {
			return fmt.Errorf("invalid print-format: %s", printFormatT)
		}
	}

	return nil
}

/* validPrintFormat check print-format if it's valid and meaningful
 * only verbs `s` and `v` are valid at the moment
 * `%5.s` is valid, but not meaningful, because the output will be empty, will not be accepted
 * `%.5s` is valid, but not very meaningful, but will be accepted
 */
func validPrintFormat(printFormat string) bool {
	var re = regexp.MustCompile(`%(-{0,1}\d{1,}){0,1}(\.\d{1,}){0,1}([sv])`)
	return re.MatchString(printFormat)
}

// ToMap ...
func (i Element) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	r["color"] = i.Color.String()
	r["print"] = i.Print
	r["print-format"] = i.PrintFormat
	return r
}

// Reset ...
func (i Element) Reset() {
	i.Color = &util.ColorT{}
	i.Color.Set("OpReset")

	i.Print = true

	i.PrintFormat = "%s"
}

// GetColor ...
func (i Element) GetColor(value string) util.Color {
	return i.Color
}

// IsEnabled ...
func (i Element) IsEnabled() bool {
	return i.Print
}

// PrintTo ...
func (i Element) PrintTo(color util.Color, builder *strings.Builder, a string) {
	a = shortenValue(a, i.PrintFormat)
	if color == nil {
		builder.WriteString(fmt.Sprintf(i.PrintFormat, a))
	} else {
		builder.WriteString(color.Sprintf(i.PrintFormat, a))
	}
}

// shortenValue shortens the value to maxWidth -3 chars if necessary, shortened values will be postfixed by three dots
func shortenValue(inValue string, printFormat string) string {
	idx := strings.Index(printFormat, ".")
	if idx >= 0 {
		width, err := strconv.Atoi(printFormat[1:idx])
		if err == nil && len([]rune(inValue)) > util.Abs(width) && util.Abs(width) > 3 {
			return fmt.Sprint(inValue[:util.Abs(width)-3], "...")
		}
	}
	return inValue
}
