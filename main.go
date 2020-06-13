package main

import (
	"log"

	"github.com/guromityan/go-esd/app"
	"github.com/guromityan/go-esd/lib"
)

func main() {
	filename := "sampledata/test.md"
	data, err := lib.MDParser(filename)
	if err != nil {
		log.Fatalln(err)
	}
	ts, err := lib.NewTestSpec(data)
	if err != nil {
		log.Fatalln(err)
	}

	app.SetData(ts)

	err = ts.Save()
	if err != nil {
		log.Fatalln(err)
	}
}
