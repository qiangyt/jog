package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/qiangyt/jog/config"
	"github.com/qiangyt/jog/util"
)

var errReadTimeout = errors.New("read timeout")
var readTimeout time.Duration = time.Millisecond * 200
var followCheckInterval = time.Millisecond * 200

// ProcessRawLine ...
func ProcessRawLine(cfg config.Configuration, options Options, lineNo int, rawLine string) {
	record := ParseAsRecord(cfg, options, lineNo, rawLine)
	if !record.MatchesFilters(cfg, options) {
		return
	}

	var line string
	if options.OutputRawJSON {
		line = record.Raw
	} else {
		line = record.AsFlatLine(cfg)
	}

	if len(line) > 0 {
		fmt.Println(line)
	}
}

// ProcessLocalFile ...
func ProcessLocalFile(cfg config.Configuration, options Options, follow bool, localFilePath string) {
	var offset int64 = 0
	var lineNo int = 1

	if !follow {
		ReadLocalFile(cfg, options, localFilePath, offset, lineNo)
		return
	}

	ticker := time.NewTicker(followCheckInterval)
	for range ticker.C {
		offset, lineNo = ReadLocalFile(cfg, options, localFilePath, offset, lineNo)
	}
}

// ReadLocalFile ...
func ReadLocalFile(cfg config.Configuration, options Options, localFilePath string, offset int64, lineNo int) (int64, int) {
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

	lineNo = ProcessReader(cfg, options, f, lineNo)

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
func ProcessReader(cfg config.Configuration, options Options, reader io.Reader, lineNo int) int {
	buf := bufio.NewReader(reader)
	isEOF := false

	if lineNo == 1 && options.NumberOfLines > 0 {

		// skip 'options.NumberOfLines' of lines
		tailQueue := util.NewTailQueue(options.NumberOfLines)
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
			ProcessRawLine(cfg, options, lineNo, rawLine)
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
			ProcessRawLine(cfg, options, lineNo, rawLine)
			return lineNo + 1
		}

		ProcessRawLine(cfg, options, lineNo, rawLine)
	}

	return lineNo
}
