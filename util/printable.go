package util

import (
	"strings"
)

// Printable ...
type Printable interface {
	IsEnabled() bool
	GetColor(value string) Color
	PrintBefore(color Color, builder *strings.Builder)
	PrintBody(color Color, builder *strings.Builder, body string)
	PrintAfter(color Color, builder *strings.Builder)
}
