package app

import (
	"fmt"

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

const startRow = 3

var header []string = []string{
	aHead, bHead, cHead, dHead, eHead, fHead, gHead, hHead, iHead, jHead,
}

func SetData(ts *lib.TestSpec) {
	genrs := ts.Data.Genres

	for _, g := range genrs {
		rowNum := startRow
		// シート名、ヘッダの設定
		setCellVal := ts.GetSetCellValFunc(g.Genre)
		mergeCell := ts.GetMergeCellFunc(g.Genre)

		setHeaders(header, setCellVal)
		for _, c := range g.Categories {
			// カテゴリの設定
			setCellVal(3, rowNum+1, c.Category)
			categoryRowNum := rowNum + 1

			for _, cs := range c.Cases {
				// ケースの設定
				setCellVal(4, rowNum+1, cs.Case)
				// ステップの設定
				steps := ""
				for si, s := range cs.Steps {
					if si == 0 {
						steps = fmt.Sprintf("%v. %v", si+1, s)
					} else {
						steps = fmt.Sprintf("%v\n%v. %v", steps, si+1, s)
					}
				}
				setCellVal(5, rowNum+1, steps)
				mergeCell(4, rowNum+1, 4, rowNum+len(cs.Checks))
				mergeCell(5, rowNum+1, 5, rowNum+len(cs.Checks))

				for chi, ch := range cs.Checks {
					rowNum++
					// 連番の設定
					setCellVal(1, rowNum, rowNum-startRow)
					// 項番の設定
					setCellVal(2, rowNum, fmt.Sprintf("%v-%v-%v", c.Num+1, cs.Num+1, chi+1))
					// チェックの設定
					setCellVal(6, rowNum, ch)
				}
			}
			mergeCell(3, categoryRowNum, 3, rowNum)
		}
	}
}

func setHeaders(header []string, f func(x, y int, val interface{}) error) {
	for i, h := range header {
		f(i+1, 3, h)
	}
}
