package grok

import (
	"path/filepath"

	"github.com/qiangyt/jog/pkg/res"
)

func ParseVjeantetPatternsStatikFile(name string) []Pattern {
	p := filepath.Join("/grok_vjeantet", name)
	res := res.New(p)
	patternsText := res.ReadString()

	return ParsePatterns(patternsText)
}

func ParseExtendedPatternsStatikFile(name string) []Pattern {
	p := filepath.Join("/grok_extended", name)
	res := res.New(p)
	patternsText := res.ReadString()

	return ParsePatterns(patternsText)
}

// CopyGrokVjeantestStatikFile ...
func CopyGrokVjeantestStatikFile(targetDir string, name string) {
	p := filepath.Join("/grok_vjeantet", name)
	res := res.New(p)

	res.CopyToFile(targetDir)
}

// CopyGrokExtendedStatikFile ...
func CopyGrokExtendedStatikFile(targetDir string, name string) {
	p := filepath.Join("/grok_extended", name)
	res := res.New(p)

	res.CopyToFile(targetDir)
}
