package main

import (
	"bufio"
	"io"
	"os"

	"github.com/pkg/errors"
)

// LogFormat ...
type LogFormat interface {
	Parse(event LogEvent) error
}

// ProcessLinesWithLocalFile ...
func ProcessLinesWithLocalFile(localFilePath string) {

	f, err := os.OpenFile(localFilePath, os.O_RDONLY, 0400)
	if err != nil {
		panic(errors.Wrapf(err, "failed to read file: %s", localFilePath))
	}
	defer f.Close()

	ProcessLinesWithReader(f)
}

// ProcessLinesWithReader ...
func ProcessLinesWithReader(reader io.Reader) {

	buf := bufio.NewReader(reader)

	for lineNo := 1; true; lineNo++ {
		raw, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				ProcessRawLine(lineNo, raw)
				return
			}
			panic(errors.Wrapf(err, "failed to read line %d", lineNo))
		}

		ProcessRawLine(lineNo, raw)
	}
}
