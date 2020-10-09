package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/qiangyt/jog/util"
)

var errReadTimeout = errors.New("read timeout")
var readTimeout time.Duration = time.Millisecond * 200
var followCheckInterval = time.Millisecond * 200

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

	ticker := time.NewTicker(followCheckInterval)
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
		if fSize <= offset {
			return fSize, lineNo
		}

		_, err := f.Seek(offset+1, 0)
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

func readRawLineWithTimeout(timer *time.Timer, buf *bufio.Reader) (string, error) {
	type ReadResult struct {
		line string
		err  error
	}
	ch := make(chan ReadResult)

	go func() {
		line, err := readRawLine(buf)
		ch <- ReadResult{line, err}
	}()

	timer.Reset(readTimeout)

	select {
	case result := <-ch:
		return result.line, result.err
	case <-timer.C:
		return "", errReadTimeout
	}
}

func readRawLine(buf *bufio.Reader) (string, error) {
	rawLine, err := buf.ReadString('\n')
	len := len(rawLine)

	if len != 0 {
		// trim the tail \n
		if rawLine[len-1] == '\n' {
			rawLine = rawLine[:len-1]
		}
	}

	return rawLine, err
}

// ProcessReader ...
func ProcessReader(cfg Config, cmdLine CommandLine, reader io.Reader, lineNo int) int {
	buf := bufio.NewReader(reader)
	isEOF := false

	if lineNo == 1 && cmdLine.NumberOfLines > 0 {

		// skip 'cmdLine.NumberOfLines' of lines
		tailQueue := util.NewTailQueue(cmdLine.NumberOfLines)
		timer := time.NewTimer(readTimeout)

		for ; true; lineNo++ {
			rawLine, err := readRawLineWithTimeout(timer, buf)
			if err != nil {
				timer.Stop()

				if err == errReadTimeout {
					isEOF = false
				} else if err != io.EOF {
					panic(errors.Wrapf(err, "failed to read line %d", lineNo))
				} else {
					isEOF = true

					tailQueue.Add(rawLine)
					log.Printf("got EOF, line %d\n", lineNo)
					lineNo++
				}

				break
			}

			tailQueue.Add(rawLine)
		}

		lineNo = lineNo - tailQueue.Count()

		for ; !tailQueue.IsEmpty(); lineNo++ {
			rawLine := tailQueue.Kick().(string)
			ProcessRawLine(cfg, cmdLine, lineNo, rawLine)
		}
	}

	if isEOF {
		return lineNo
	}

	for ; true; lineNo++ {
		rawLine, err := readRawLine(buf)

		if err != nil {
			if err != io.EOF {
				panic(errors.Wrapf(err, "failed to read line %d", lineNo))
			}

			log.Printf("got EOF, line %d\n", lineNo)
			ProcessRawLine(cfg, cmdLine, lineNo, rawLine)
			return lineNo + 1
		}

		ProcessRawLine(cfg, cmdLine, lineNo, rawLine)
	}

	return lineNo
}