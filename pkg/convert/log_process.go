package convert

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
	_util "github.com/qiangyt/jog/pkg/util"
)

var errReadTimeout = errors.New("read timeout")
var readTimeout time.Duration = time.Millisecond * 200
var followCheckInterval = time.Millisecond * 200

// ProcessRawLine ...
func ProcessRawLine(ctx ConvertContext, lineNo int, rawLine string) {
	record := ParseAsRecord(ctx, lineNo, rawLine)
	if !record.MatchesFilters(ctx) {
		return
	}

	var line string
	if ctx.Options().OutputRawJSON {
		line = record.Raw
	} else {
		line = record.AsFlatLine(ctx.Config())
	}

	if len(line) > 0 {
		fmt.Println(line)
	}
}

// ProcessLocalFile ...
func ProcessLocalFile(ctx ConvertContext) {
	var offset int64 = 0
	var lineNo int = 1

	if !ctx.Options().FollowMode {
		ReadLocalFile(ctx, offset, lineNo)
		return
	}

	ticker := time.NewTicker(followCheckInterval)
	for range ticker.C {
		offset, lineNo = ReadLocalFile(ctx, offset, lineNo)
	}
}

// ReadLocalFile ...
func ReadLocalFile(ctx ConvertContext, offset int64, lineNo int) (int64, int) {
	localFilePath := ctx.Options().LogFilePath

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

		if offset+1 < fSize {
			tmp := make([]byte, 1)
			if _, err := f.ReadAt(tmp, offset); err != nil {
				panic(errors.Wrapf(err, "failed to read at: %s/%v", localFilePath, offset+1))
			}
			if tmp[0] == '\n' {
				offset = offset + 1
			}
		}

		_, err := f.Seek(offset, 0)
		if err != nil {
			panic(errors.Wrapf(err, "failed to seek: %s/%v", localFilePath, offset))
		}
	}

	lineNo = ProcessReader(ctx, f, lineNo)

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
func ProcessReader(ctx ConvertContext, reader io.Reader, lineNo int) int {
	buf := bufio.NewReader(reader)
	isEOF := false

	if lineNo == 1 && ctx.Options().NumberOfLines > 0 {

		// skip 'options.NumberOfLines' of lines
		tailQueue := _util.NewTailQueue(ctx.Options().NumberOfLines)
		timer := time.NewTimer(readTimeout)

		for {
			rawLine, err := readRawLineWithTimeout(timer, buf)
			if err != nil {
				timer.Stop()

				if err == errReadTimeout {
					isEOF = false
				} else if err != io.EOF {
					panic(errors.Wrapf(err, "failed to read line %d", lineNo))
				} else {
					isEOF = true

					ctx.LogInfo("got EOF", "lineNo", lineNo)

					if len(rawLine) > 0 {
						if rawLine[0] != '\n' {
							tailQueue.Add(rawLine)
							lineNo++
						}
					}
				}

				break
			}

			tailQueue.Add(rawLine)
			lineNo++
		}

		lineNo = lineNo - tailQueue.Count()

		for ; !tailQueue.IsEmpty(); lineNo++ {
			rawLine := tailQueue.Kick().(string)
			ProcessRawLine(ctx, lineNo, rawLine)
		}
	}

	if isEOF {
		return lineNo
	}

	for {
		rawLine, err := readRawLine(buf)

		if err != nil {
			if err != io.EOF {
				panic(errors.Wrapf(err, "failed to read line %d", lineNo))
			}

			ctx.LogInfo("got EOF", "lineNo", lineNo)

			if len(rawLine) > 0 {
				if rawLine[0] != '\n' {
					ProcessRawLine(ctx, lineNo, rawLine)
					lineNo++
				}
			}
			return lineNo
		}

		ProcessRawLine(ctx, lineNo, rawLine)
		lineNo++
	}
}
