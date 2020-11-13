package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

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

// Open ...
func (i LogFile) Open() {
	p := i.path

	create := true
	if fi := FileStat(p, false); fi != nil {
		if fi.Size() >= 100*1024*1024 {
			RemoveFile(p)
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

	f, err := os.OpenFile(p, flg, 0666)
	if err != nil {
		panic(errors.Wrapf(err, "failed to create/open log file: %s", p))
	}

	i.file = f

	log.SetOutput(i)
	log.SetPrefix(fmt.Sprintf("[%5d] ", os.Getpid()))

	if !create {
		i.file.Write([]byte("-------------------------------------------------------------------------------\n"))
	}
	log.Printf("started at: %v\n", time.Now())
}

// Close ...
func (i LogFile) Close() {
	if i.file != nil {
		i.file.Close()
	}
}

// InitLogger ...
func InitLogger(jogHomeDir string) LogFile {
	r := &LogFileT{
		path: filepath.Join(jogHomeDir, "jog.log"),
	}
	r.Open()
	return r
}
