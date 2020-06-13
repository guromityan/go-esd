package main

import (
	"log"

	"github.com/guromityan/go-esd/app"
	"github.com/guromityan/go-esd/lib"
)

const (
	aHead = "連番"
	bHead = "項番"
	cHead = "カテゴリ"
	dHead = "ケース"
	eHead = "確認手順"
	fHead = "期待値"
	gHead = "結果"
	hHead = "確認日"
	iHead = "確認者"
	jHead = "備考"
)

var header []string = []string{
	aHead, bHead, cHead, dHead, eHead, fHead, gHead, hHead, iHead, jHead,
}

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
	f := ts.GetSetCellValFunc("Sheet1")
	app.SetHeader(header, f)

	err = ts.Save()
	if err != nil {
		log.Fatalln(err)
	}
}
