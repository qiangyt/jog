package main

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
func (me LogFile) Write(p []byte) (int, error) {
	return me.file.Write(p)
}

// Open ...
func (me LogFile) Open() {
	p := me.path

	create := true
	if fi := FileStat(p, false); fi != nil {
		if fi.Size() >= 10*1024*1024 {
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

	me.file = f

	log.SetOutput(me)
	log.SetPrefix(fmt.Sprintf("[%5d] ", os.Getpid()))

	if !create {
		me.file.Write([]byte("-------------------------------------------------------------------------------\n"))
	}
	log.Printf("started at: %v\n", time.Now())
}

// Close ...
func (me LogFile) Close() {
	if me.file != nil {
		me.file.Close()
	}
}

// InitLogger ...
func InitLogger() LogFile {
	r := &LogFileT{
		path: filepath.Join(ExeDirectory(), "jog.log"),
	}
	r.Open()
	return r
}
