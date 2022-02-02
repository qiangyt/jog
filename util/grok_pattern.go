package util

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/qiangyt/jog/static/grok_extended"
	"github.com/qiangyt/jog/static/grok_vjeantet"
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

func LoadAllGrokPatterns() []GrokPattern {
	r := []GrokPattern{}

	r = append(r, ParseGrokPatterns(grok_vjeantet.Aws)...)
	r = append(r, ParseGrokPatterns(grok_vjeantet.Bacula)...)
	r = append(r, ParseGrokPatterns(grok_vjeantet.Bro)...)
	r = append(r, ParseGrokPatterns(grok_vjeantet.Exim)...)
	r = append(r, ParseGrokPatterns(grok_vjeantet.Firewalls)...)
	r = append(r, ParseGrokPatterns(grok_vjeantet.Grok_patterns)...)
	r = append(r, ParseGrokPatterns(grok_vjeantet.Haproxy)...)
	r = append(r, ParseGrokPatterns(grok_vjeantet.Java)...)
	r = append(r, ParseGrokPatterns(grok_vjeantet.Junos)...)
	r = append(r, ParseGrokPatterns(grok_vjeantet.Linux_syslog)...)
	r = append(r, ParseGrokPatterns(grok_vjeantet.Mcollective_patterns)...)
	r = append(r, ParseGrokPatterns(grok_vjeantet.Mcollective)...)
	r = append(r, ParseGrokPatterns(grok_vjeantet.Mongodb)...)
	r = append(r, ParseGrokPatterns(grok_vjeantet.Nagios)...)
	r = append(r, ParseGrokPatterns(grok_vjeantet.Postgresql)...)
	r = append(r, ParseGrokPatterns(grok_vjeantet.Rails)...)
	r = append(r, ParseGrokPatterns(grok_vjeantet.Redis)...)
	r = append(r, ParseGrokPatterns(grok_vjeantet.Ruby)...)

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

// SaveGrokPatternFile ...
func SaveGrokPatternFile(dir string, patternFileName string, patternFileContent string) {
	ReplaceFile(filepath.Join(dir, patternFileName), []byte(patternFileContent))
}

// DefaultGrokLibraryDirs ...
func DefaultGrokLibraryDirs(expand bool) []string {
	return []string{
		JogHomeDir(expand, "grok_vjeantet"),
		JogHomeDir(expand, "grok_extended"),
	}
}

// ResetDefaultGrokLibraryDir ...
func ResetDefaultGrokLibraryDir() {
	dirVjeantet := JogHomeDir(true, "grok_vjeantet")
	RemoveDir(dirVjeantet)

	dirExtended := JogHomeDir(true, "grok_extended")
	RemoveDir(dirExtended)

	InitDefaultGrokLibraryDir()
}

// InitDefaultGrokLibraryDir ...
func InitDefaultGrokLibraryDir() {
	jogHomeDir := JogHomeDir(true)

	licensePath := filepath.Join(jogHomeDir, "grok_vjeantet.LICENSE")
	WriteFileIfNotFound(licensePath, []byte(grok_vjeantet.LICENSE))

	readmePath := filepath.Join(jogHomeDir, "grok_vjeantet.README.md")
	WriteFileIfNotFound(readmePath, []byte(grok_vjeantet.README_md))

	dirVjeantet := JogHomeDir(true, "grok_vjeantet")
	if DirExists(dirVjeantet) == false {
		MkdirAll(dirVjeantet)

		SaveGrokPatternFile(dirVjeantet, "aws", grok_vjeantet.Aws)
		SaveGrokPatternFile(dirVjeantet, "bro", grok_vjeantet.Bro)
		SaveGrokPatternFile(dirVjeantet, "firewalls", grok_vjeantet.Firewalls)
		SaveGrokPatternFile(dirVjeantet, "haproxy", grok_vjeantet.Haproxy)
		SaveGrokPatternFile(dirVjeantet, "junos", grok_vjeantet.Junos)
		SaveGrokPatternFile(dirVjeantet, "linux-syslog", grok_vjeantet.Linux_syslog)
		SaveGrokPatternFile(dirVjeantet, "mcollective-patterns", grok_vjeantet.Mcollective_patterns)
		SaveGrokPatternFile(dirVjeantet, "nagios", grok_vjeantet.Nagios)
		SaveGrokPatternFile(dirVjeantet, "rails", grok_vjeantet.Rails)
		SaveGrokPatternFile(dirVjeantet, "redis", grok_vjeantet.Redis)
		SaveGrokPatternFile(dirVjeantet, "bacula", grok_vjeantet.Bacula)
		SaveGrokPatternFile(dirVjeantet, "exim", grok_vjeantet.Exim)
		SaveGrokPatternFile(dirVjeantet, "grok-patterns", grok_vjeantet.Grok_patterns)
		SaveGrokPatternFile(dirVjeantet, "java", grok_vjeantet.Java)
		SaveGrokPatternFile(dirVjeantet, "mcollective", grok_vjeantet.Mcollective)
		SaveGrokPatternFile(dirVjeantet, "mongodb", grok_vjeantet.Mongodb)
		SaveGrokPatternFile(dirVjeantet, "postgresql", grok_vjeantet.Postgresql)
		SaveGrokPatternFile(dirVjeantet, "ruby", grok_vjeantet.Ruby)
	}

	dirExtended := JogHomeDir(true, "grok_extended")
	if DirExists(dirExtended) == false {
		MkdirAll(dirExtended)

		SaveGrokPatternFile(dirExtended, "pm2", grok_extended.Pm2)
	}

	MkdirAll(JogHomeDir(true, "grok_mine"))

}
