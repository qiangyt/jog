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
	BaseDir = "./static_files"
)

// Reads all static files in the current folder
// and encodes them as strings literals in static/static.go
func main() {
	os.MkdirAll("static", os.ModePerm)

	fs, _ := ioutil.ReadDir(BaseDir)

	for _, f := range fs {
		fName := f.Name()

		fmt.Println("Including static file: " + fName)

		var fTitle, fExt string

		indexOfLastDot := strings.LastIndex(fName, ".")
		if indexOfLastDot >= 0 {
			fTitle = fName[:indexOfLastDot]
			fExt = fName[indexOfLastDot+1:]
		} else {
			fTitle = fName
			fExt = ""
		}

		indexOfLastSlash := strings.LastIndex(fTitle, "/")
		if indexOfLastSlash >= 0 {
			fTitle = fTitle[indexOfLastSlash+1:]
		}

		// TODO: this code does not work for text files that contain: %  `

		out, _ := os.Create("static/" + fName + ".go")
		out.WriteString("package static \n\nconst (\n")

		varName := fTitle
		if len(fExt) > 0 {
			varName = varName + "_" + fExt
		}
		out.WriteString("  " + varName + " = `")

		f, _ := os.Open(BaseDir + "/" + fName)
		io.Copy(out, f)
		out.WriteString("`\n")
		out.WriteString(")\n")
	}
}
