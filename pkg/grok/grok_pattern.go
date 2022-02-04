package grok

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/qiangyt/jog/pkg/res"
	"github.com/qiangyt/jog/pkg/util"
)

type GrokPattern struct {
	Name string
	Expr string
}

func ParseGrokPatterns(patternsText string) []GrokPattern {
	r := make([]GrokPattern, 0)

	buf := bytes.NewBufferString(patternsText)

	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		l := scanner.Text()
		l = strings.TrimSpace(l)
		if len(l) > 0 && l[0] != '#' {
			nameAndExpr := strings.SplitN(l, " ", 2)

			p := GrokPattern{}
			p.Name = nameAndExpr[0]
			p.Expr = nameAndExpr[1]

			r = append(r, p)
		}
	}

	return r
}

func ParseVjeantetGrokPatternsStatikFile(name string) []GrokPattern {
	p := filepath.Join("/grok_vjeantet", name)
	res := res.New(p)
	patternsText := res.ReadString()

	return ParseGrokPatterns(patternsText)
}

func ParseExtendedGrokPatternsStatikFile(name string) []GrokPattern {
	p := filepath.Join("/grok_extended", name)
	res := res.New(p)
	patternsText := res.ReadString()

	return ParseGrokPatterns(patternsText)
}

func LoadAllGrokPatterns() []GrokPattern {
	r := []GrokPattern{}

	r = append(r, ParseVjeantetGrokPatternsStatikFile("aws")...)
	r = append(r, ParseVjeantetGrokPatternsStatikFile("bacula")...)
	r = append(r, ParseVjeantetGrokPatternsStatikFile("bro")...)
	r = append(r, ParseVjeantetGrokPatternsStatikFile("exim")...)
	r = append(r, ParseVjeantetGrokPatternsStatikFile("firewalls")...)
	r = append(r, ParseVjeantetGrokPatternsStatikFile("grok-patterns")...)
	r = append(r, ParseVjeantetGrokPatternsStatikFile("haproxy")...)
	r = append(r, ParseVjeantetGrokPatternsStatikFile("java")...)
	r = append(r, ParseVjeantetGrokPatternsStatikFile("junos")...)
	r = append(r, ParseVjeantetGrokPatternsStatikFile("linux-syslog")...)
	r = append(r, ParseVjeantetGrokPatternsStatikFile("mcollective-patterns")...)
	r = append(r, ParseVjeantetGrokPatternsStatikFile("mcollective")...)
	r = append(r, ParseVjeantetGrokPatternsStatikFile("mongodb")...)
	r = append(r, ParseVjeantetGrokPatternsStatikFile("nagios")...)
	r = append(r, ParseVjeantetGrokPatternsStatikFile("postgresql")...)
	r = append(r, ParseVjeantetGrokPatternsStatikFile("rails")...)
	r = append(r, ParseVjeantetGrokPatternsStatikFile("redis")...)
	r = append(r, ParseVjeantetGrokPatternsStatikFile("ruby")...)

	return r
}

func MergeGrokPatterns(allPatterns map[string]GrokPattern, patternsText string) {
	newPatterns := ParseGrokPatterns(patternsText)
	for _, pattern := range newPatterns {
		name := pattern.Name
		if existingOne, alreadyExists := allPatterns[name]; alreadyExists == true {
			panic(fmt.Errorf("duplicated grok pattern. name: %s. existing: %s. duplicated: %s", name, existingOne.Expr, pattern.Expr))
		}
		allPatterns[name] = pattern
	}
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

// DefaultGrokLibraryDirs ...
func DefaultGrokLibraryDirs(expand bool) []string {
	return []string{
		util.JogHomeDir(expand, "grok_vjeantet"),
		util.JogHomeDir(expand, "grok_extended"),
	}
}

// ResetDefaultGrokLibraryDir ...
func ResetDefaultGrokLibraryDir() {
	dirVjeantet := util.JogHomeDir(true, "grok_vjeantet")
	util.RemoveDir(dirVjeantet)

	dirExtended := util.JogHomeDir(true, "grok_extended")
	util.RemoveDir(dirExtended)

	InitDefaultGrokLibraryDir()
}

// InitDefaultGrokLibraryDir ...
func InitDefaultGrokLibraryDir() {
	jogHomeDir := util.JogHomeDir(true)

	if util.DirExists(filepath.Join(jogHomeDir, "grok_vjeantet")) == false {
		CopyGrokVjeantestStatikFile(jogHomeDir, "LICENSE")
		CopyGrokVjeantestStatikFile(jogHomeDir, "README.md")

		CopyGrokVjeantestStatikFile(jogHomeDir, "aws")
		CopyGrokVjeantestStatikFile(jogHomeDir, "bro")
		CopyGrokVjeantestStatikFile(jogHomeDir, "firewalls")
		CopyGrokVjeantestStatikFile(jogHomeDir, "haproxy")
		CopyGrokVjeantestStatikFile(jogHomeDir, "junos")
		CopyGrokVjeantestStatikFile(jogHomeDir, "linux-syslog")
		CopyGrokVjeantestStatikFile(jogHomeDir, "mcollective-patterns")
		CopyGrokVjeantestStatikFile(jogHomeDir, "nagios")
		CopyGrokVjeantestStatikFile(jogHomeDir, "rails")
		CopyGrokVjeantestStatikFile(jogHomeDir, "redis")
		CopyGrokVjeantestStatikFile(jogHomeDir, "bacula")
		CopyGrokVjeantestStatikFile(jogHomeDir, "exim")
		CopyGrokVjeantestStatikFile(jogHomeDir, "grok-patterns")
		CopyGrokVjeantestStatikFile(jogHomeDir, "java")
		CopyGrokVjeantestStatikFile(jogHomeDir, "mcollective")
		CopyGrokVjeantestStatikFile(jogHomeDir, "mongodb")
		CopyGrokVjeantestStatikFile(jogHomeDir, "postgresql")
		CopyGrokVjeantestStatikFile(jogHomeDir, "ruby")
	}

	dirExtended := util.JogHomeDir(true, "grok_extended")
	if util.DirExists(dirExtended) == false {
		CopyGrokExtendedStatikFile(jogHomeDir, "pm2")
	}

	util.MkdirAll(util.JogHomeDir(true, "grok_mine"))

}
