package main

import (
	"fmt"
)

func main() {
	/*if len(os.Args) == 2 {
		filePath := os.Args[1]
		ProcessLinesWithLocalFile(filePath)
	} else {
		ProcessLinesWithReader(os.Stdin)
	}*/

	ProcessLinesWithLocalFile("/Users/i508673/github/qiangyt/json2log/example_logs/bts.log")
	fmt.Println()
}
