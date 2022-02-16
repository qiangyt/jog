package grok

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	_io "github.com/qiangyt/jog/pkg/io"
)

type Pattern struct {
	Name string
	Expr string
}

func ParsePatterns(patternsText string) []Pattern {
	r := make([]Pattern, 0)

	buf := bytes.NewBufferString(patternsText)

	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		l := scanner.Text()
		l = strings.TrimSpace(l)
		if len(l) > 0 && l[0] != '#' {
			nameAndExpr := strings.SplitN(l, " ", 2)

			p := Pattern{}
			p.Name = nameAndExpr[0]
			p.Expr = nameAndExpr[1]

			r = append(r, p)
		}
	}

	return r
}

func LoadAllPatterns() []Pattern {
	r := []Pattern{}

	r = append(r, ParseVjeantetPatternsStatikFile("aws")...)
	r = append(r, ParseVjeantetPatternsStatikFile("bacula")...)
	r = append(r, ParseVjeantetPatternsStatikFile("bro")...)
	r = append(r, ParseVjeantetPatternsStatikFile("exim")...)
	r = append(r, ParseVjeantetPatternsStatikFile("firewalls")...)
	r = append(r, ParseVjeantetPatternsStatikFile("grok-patterns")...)
	r = append(r, ParseVjeantetPatternsStatikFile("haproxy")...)
	r = append(r, ParseVjeantetPatternsStatikFile("java")...)
	r = append(r, ParseVjeantetPatternsStatikFile("junos")...)
	r = append(r, ParseVjeantetPatternsStatikFile("linux-syslog")...)
	r = append(r, ParseVjeantetPatternsStatikFile("mcollective-patterns")...)
	r = append(r, ParseVjeantetPatternsStatikFile("mcollective")...)
	r = append(r, ParseVjeantetPatternsStatikFile("mongodb")...)
	r = append(r, ParseVjeantetPatternsStatikFile("nagios")...)
	r = append(r, ParseVjeantetPatternsStatikFile("postgresql")...)
	r = append(r, ParseVjeantetPatternsStatikFile("rails")...)
	r = append(r, ParseVjeantetPatternsStatikFile("redis")...)
	r = append(r, ParseVjeantetPatternsStatikFile("ruby")...)

	return r
}

func MergePatterns(allPatterns map[string]Pattern, patternsText string) {
	newPatterns := ParsePatterns(patternsText)
	for _, pattern := range newPatterns {
		name := pattern.Name
		if existingOne, alreadyExists := allPatterns[name]; alreadyExists == true {
			panic(fmt.Errorf("duplicated grok pattern. name: %s. existing: %s. duplicated: %s", name, existingOne.Expr, pattern.Expr))
		}
		allPatterns[name] = pattern
	}
}

// DefaultGrokLibraryDirs ...
func DefaultGrokLibraryDirs(expand bool) []string {
	return []string{
		_io.JogHomeDir(expand, "grok_vjeantet"),
		_io.JogHomeDir(expand, "grok_extended"),
	}
}

// ResetDefaultGrokLibraryDir ...
func ResetDefaultGrokLibraryDir() {
	dirVjeantet := _io.JogHomeDir(true, "grok_vjeantet")
	_io.RemoveDir(dirVjeantet)

	dirExtended := _io.JogHomeDir(true, "grok_extended")
	_io.RemoveDir(dirExtended)

	InitDefaultGrokLibraryDir()
}

// InitDefaultGrokLibraryDir ...
func InitDefaultGrokLibraryDir() {
	jogHomeDir := _io.JogHomeDir(true)

	if _io.DirExists(filepath.Join(jogHomeDir, "grok_vjeantet")) == false {
		//TODO: read directory in a loop
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

	dirExtended := _io.JogHomeDir(true, "grok_extended")
	if _io.DirExists(dirExtended) == false {
		CopyGrokExtendedStatikFile(jogHomeDir, "pm2")
	}

	_io.MkdirAll(_io.JogHomeDir(true, "grok_mine"))

}
