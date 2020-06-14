package main

import (
	"log"

	"github.com/guromityan/go-esd/app"
	"github.com/guromityan/go-esd/lib"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// コマンドオプション
	mdFile  = kingpin.Flag("mdfile", "Markdown file to convert.").Short('f').Required().ExistingFile()
	destDir = kingpin.Flag("dest", "Existing directory to output the file.").Short('d').ExistingDir()
)

func main() {
	// 引数をパース
	kingpin.Parse()

	// Markdown をパース
	data, err := lib.MDParser(*mdFile)
	if err != nil {
		log.Fatalln(err)
	}

	// 出力先パスを設定
	err = data.SetPath(*destDir)
	if err != nil {
		log.Fatalln(err)
	}

	// Excel に変換する構造体を定義
	ts, err := lib.NewTestSpec(data)
	if err != nil {
		log.Fatalln(err)
	}

	// Excel にデータを書込
	app.SetData(ts)

	// 保存して終了
	err = ts.Save()
	if err != nil {
		log.Fatalln(err)
	}
}
