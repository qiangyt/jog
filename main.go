package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	logFile := InitLogger()
	defer logFile.Close()

	cfg := LoadConfig()

	var filePath string
	filePath = "./example_logs/logstash.log"

	if len(filePath) == 0 && len(os.Args) == 2 {
		filePath = os.Args[1]
	}

	if len(filePath) == 0 {
		log.Println("Read log lines from stdin")
		ProcessLinesWithReader(cfg, os.Stdin)
	} else {
		log.Printf("processing local file: %s\n", filePath)
		ProcessLinesWithLocalFile(cfg, filePath)
	}

	fmt.Println()
}
