package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
)

// ProcessRawLine ...
func ProcessRawLine(cfg Config, cmdLine CommandLine, lineNo int, rawLine string) {
	event := ParseAsRecord(cfg, lineNo, rawLine)
	var line = event.AsFlatLine(cfg)
	if len(line) > 0 {
		fmt.Println(line)
	}
}

// ProcessLocalFile ...
func ProcessLocalFile(cfg Config, cmdLine CommandLine, follow bool, localFilePath string) {
	var offset int64 = 0
	var lineNo int = 1

	if !follow {
		ReadLocalFile(cfg, cmdLine, localFilePath, offset, lineNo)
		return
	}

	ticker := time.NewTicker(time.Millisecond * 500)
	for range ticker.C {
		offset, lineNo = ReadLocalFile(cfg, cmdLine, localFilePath, offset, lineNo)
	}
}

// ReadLocalFile ...
func ReadLocalFile(cfg Config, cmdLine CommandLine, localFilePath string, offset int64, lineNo int) (int64, int) {
	f, err := os.Open(localFilePath)
	if err != nil {
		panic(errors.Wrapf(err, "failed to open: %s", localFilePath))
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		panic(errors.Wrapf(err, "failed to stat: %s", localFilePath))
	}
	fSize := fi.Size()

	if offset > 0 {
		if fSize == offset {
			return fSize, lineNo
		}

		_, err := f.Seek(offset, 0)
		if err != nil {
			panic(errors.Wrapf(err, "failed to seek: %s/%v", localFilePath, offset))
		}
	}

	lineNo = ProcessReader(cfg, cmdLine, f, lineNo)

	fi, err = f.Stat()
	if err != nil {
		panic(errors.Wrapf(err, "failed to stat: %s", localFilePath))
	}
	return fi.Size(), lineNo
}

// ProcessReader ...
func ProcessReader(cfg Config, cmdLine CommandLine, reader io.Reader, lineNo int) int {

	buf := bufio.NewReader(reader)

	for ; true; lineNo++ {
		rawLine, err := buf.ReadString('\n')
		len := len(rawLine)

		if len != 0 {
			// trim the tail \n
			if rawLine[len-1] == '\n' {
				rawLine = rawLine[:len-1]
			}
		}

		if err != nil {
			if err == io.EOF {
				log.Printf("got EOF, line %d\n", lineNo)
				ProcessRawLine(cfg, cmdLine, lineNo, rawLine)
				return lineNo + 1
			}
			panic(errors.Wrapf(err, "failed to read line %d", lineNo))
		}

		ProcessRawLine(cfg, cmdLine, lineNo, rawLine)
	}

	return lineNo
}
