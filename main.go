package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	Colors["Blue"].Println("Simple to use color")

	logFile := InitLogger()
	defer logFile.Close()

	var filePath string
	filePath = "./example_logs/logstash.log"

	if len(filePath) == 0 && len(os.Args) == 2 {
		filePath = os.Args[1]
	}

	if len(filePath) == 0 {
		log.Println("Read log lines from stdin")
		//ProcessLinesWithReader(os.Stdin)
	} else {
		log.Printf("processing local file: %s\n", filePath)
		//ProcessLinesWithLocalFile(filePath)
	}

	fmt.Println()
}
