package main

import "github.com/qiangyt/jog/convert"

//go:generate go run script/include_static.go

func main() {
	convert.Main()
}
