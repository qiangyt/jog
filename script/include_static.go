package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const (
	// BaseDir ...
	BaseDir = "./static"
)

// Reads all static files in the current folder
// and encodes them as strings literals in static/static.go
func main() {
	fs, _ := ioutil.ReadDir(BaseDir)

	out, _ := os.Create("static/static_files.go")
	out.WriteString("package static \n\nconst (\n")

	for _, f := range fs {
		fName := f.Name()
		if fName == "static_files.go" {
			continue
		}

		fmt.Println("Including static file: " + fName)

		varName := strings.ReplaceAll(fName, ".", "_")
		out.WriteString("  " + varName + " = `")
		f, _ := os.Open(BaseDir + "/" + fName)
		io.Copy(out, f)
		out.WriteString("`\n")

	}
	out.WriteString(")\n")
}
