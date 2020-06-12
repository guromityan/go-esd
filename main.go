package main

import (
	"github.com/guromityan/go-esd/lib"
)

func main() {
	filename := "sampledata/test.md"

	lib.MDParser(filename)
}
