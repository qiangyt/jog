package log

import (
	"os"

	"github.com/pkg/errors"
	jogio "github.com/qiangyt/jog/pkg/io"
)

// FileT implements io.Writer
type FileT struct {
	path string
	file *os.File
}

// File ...
type File = *FileT

// Write ...
func (i File) Write(p []byte) (int, error) {
	return i.file.Write(p)
}

func (i File) Path() string {
	return i.path
}

func (i File) File() *os.File {
	return i.file
}

// Open ...
func NewFile(path string) File {
	r := &FileT{path: path}

	create := true
	if fi := jogio.FileStat(path, false); fi != nil {
		if fi.Size() >= 100*1024*1024 {
			jogio.RemoveFile(path)
		} else {
			create = false
		}
	}

	flg := os.O_RDWR | os.O_EXCL // | os.O_SYNC
	if create {
		flg = flg | os.O_CREATE
	} else {
		flg = flg | os.O_APPEND
	}

	f, err := os.OpenFile(path, flg, 0666)
	if err != nil {
		panic(errors.Wrapf(err, "failed to create/open log file: %s", path))
	}

	r.file = f

	return r
}

// Close ...
func (i File) Close() {
	if i.file != nil {
		i.file.Close()
	}
}
