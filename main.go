package main

import (
	"fmt"

	"github.com/gookit/color"
)

func main() {
	color.Blue.Println("Simple to use color")

	log := InitLogger()
	defer log.Close()

	/*if len(os.Args) == 2 {
		filePath := os.Args[1]
		ProcessLinesWithLocalFile(filePath)
	} else {
		log.Println("Read log lines from stdin")
		ProcessLinesWithReader(os.Stdin)
	}*/

	ProcessLinesWithLocalFile("/Users/i508673/github/qiangyt/json2log/example_logs/bts.log")
	fmt.Println()
}
