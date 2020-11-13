package util

import (
	"strings"
)

// Printable ...
type Printable interface {
	IsEnabled() bool
	GetColor(value string) Color
	PrintBody(color Color, builder *strings.Builder, body string)
}
