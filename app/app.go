package app

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/guromityan/go-esd/lib"
)

var startRow int = 3

type rowConf struct {
	Header string
	Width  int
	Column string
}

var rowConfs []rowConf = []rowConf{
	{Header: "連番", Width: 5, Column: "A"},
	{Header: "項版", Width: 5, Column: "B"},
	{Header: "カテゴリ", Width: 5, Column: "C"},
	{Header: "ケース", Width: 5, Column: "D"},
	{Header: "確認手順", Width: 5, Column: "E"},
	{Header: "期待値", Width: 5, Column: "F"},
	{Header: "結果", Width: 5, Column: "G"},
	{Header: "確認日", Width: 5, Column: "H"},
	{Header: "確認者", Width: 5, Column: "I"},
	{Header: "備考", Width: 5, Column: "J"},
}

func SetData(ts *lib.TestSpec) error {
	genrs := ts.Data.Genres

	for _, g := range genrs {
		rowNum := startRow
		// シート名、ヘッダの設定
		setCellVal := ts.GetSetCellValFunc(g.Genre)
		mergeCell := ts.GetMergeCellFunc(g.Genre)
		setStyle := ts.GetSetStyleFunc(g.Genre)

		setHeaders(rowConfs, setCellVal)
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
					termNum := fmt.Sprintf("%v-%v-%v", c.Num+1, cs.Num+1, chi+1)
					setCellVal(2, rowNum, termNum)
					// チェックの設定
					setCellVal(6, rowNum, ch)
				}
			}
			mergeCell(3, categoryRowNum, 3, rowNum)
		}

		startAxis, _ := excelize.CoordinatesToCellName(1, startRow)
		endAxis, _ := excelize.CoordinatesToCellName(10, rowNum)
		setStyle(startAxis, endAxis)
	}
	return nil
}

func setHeaders(rowConfs []rowConf, f func(x, y int, val interface{}) error) {
	for i, rc := range rowConfs {
		f(i+1, 3, rc.Header)
	}
}
