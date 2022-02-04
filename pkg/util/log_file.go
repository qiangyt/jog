package util

import (
	"os"

	"github.com/pkg/errors"
)

// LogFileT implements io.Writer
type LogFileT struct {
	path string
	file *os.File
}

// LogFile ...
type LogFile = *LogFileT

// Write ...
func (i LogFile) Write(p []byte) (int, error) {
	return i.file.Write(p)
}

func (i LogFile) Path() string {
	return i.path
}

func (i LogFile) File() *os.File {
	return i.file
}

// Open ...
func NewLogFile(path string) LogFile {
	r := &LogFileT{path: path}

	create := true
	if fi := FileStat(path, false); fi != nil {
		if fi.Size() >= 100*1024*1024 {
			RemoveFile(path)
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
func (i LogFile) Close() {
	if i.file != nil {
		i.file.Close()
	}
}
