package util

import (
	"path/filepath"
)

const (
	jogHomeDir = "~/.jog"
)

// JogHomeDir ...
func JogHomeDir(expand bool, children ...string) string {
	var r string

	if !expand {
		r = jogHomeDir
	} else {
		r = ExpandHomePath(jogHomeDir)
		MkdirAll(r)
	}

	return filepath.Join(r, filepath.Join(children...))
}
