package util

import (
	"strings"
)

// Printable ...
type Printable interface {
	IsEnabled() bool
	GetColor(value string) Color
	PrintTo(color Color, builder *strings.Builder, a string)
}
