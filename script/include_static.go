package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func includeDir(staticGoParentDir string, staticFileParentDir string, dirName string) {
	staticGoDir := filepath.Join(staticGoParentDir, dirName)
	os.MkdirAll(staticGoDir, os.ModePerm)

	staticFilesDir := filepath.Join(staticFileParentDir, dirName)
	fs, _ := ioutil.ReadDir(staticFilesDir)

	for _, f := range fs {
		fName := f.Name()

		if f.IsDir() {
			includeDir(staticGoDir, staticFilesDir, fName)
		} else {
			includeFile(staticGoDir, staticFilesDir, fName)
		}
	}
}

func includeFile(staticGoParentDir string, staticFileParentDir string, fName string) {
	packageName := staticGoParentDir[strings.LastIndex(staticGoParentDir, "/")+1:]

	fPath := filepath.Join(staticFileParentDir, fName)
	fmt.Println("Including static file: " + fPath)

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

	out, _ := os.Create(filepath.Join(staticGoParentDir, fName+".go"))
	defer out.Close()
	out.WriteString("package " + packageName + " \n\nconst (\n")

	varName := fTitle
	if len(fExt) > 0 {
		varName = varName + "_" + fExt
	}
	varName = strings.ToUpper(varName[:1]) + varName[1:]
	varName = strings.ReplaceAll(varName, "-", "_")
	out.WriteString("  // " + varName + " ...\n")
	out.WriteString("  " + varName + " string = `\n")

	contentBytes, _ := ioutil.ReadFile(fPath)
	content := string(contentBytes)
	content = strings.ReplaceAll(content, "`", "` + \"`\" + `")
	out.WriteString(content)

	out.WriteString("`\n")
	out.WriteString(")\n")
}

// Reads all static files in the current folder
// and encodes them as strings literals in static/static.go
func main() {
	fs, _ := ioutil.ReadDir("./static_files")

	for _, f := range fs {
		fName := f.Name()
		if f.IsDir() {
			includeDir("./static", "./static_files", fName)
		} else {
			includeFile("./static", "./static_files", fName)
		}
	}
}
