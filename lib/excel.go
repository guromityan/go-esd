package lib

import (
	"fmt"
	"path/filepath"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const rowHeightBase = 25

const cellLineStyle = `{
	"border":[
		{"type":"bottom","color":"000000","style":1},
		{"type":"top",   "color":"000000","style":1},
		{"type":"left",  "color":"000000","style":1},
		{"type":"right", "color":"000000","style":1}
	],
	"alignment":{
		"horizontal": "left",
		"vertical": "center",
		"wrap_text": true
		}
	}`

// TestSpec テストデータ
type TestSpec struct {
	Data Tests
	File *excelize.File
}

// NewTestSpec テストデータ生成
func NewTestSpec(data *Tests) (*TestSpec, error) {
	filename := fmt.Sprintf("esd-%v.xlsx", data.Name)
	filepath := filepath.Clean(fmt.Sprintf("%v/%v", data.path, filename))
	fmt.Println(filepath)
	file, err := newExcelFile(filepath)
	if err != nil {
		return nil, err
	}
	return &TestSpec{
		Data: *data,
		File: file,
	}, nil
}

// Save 全てのシートのアクティブセルを A1 に設定して保存
func (ts *TestSpec) Save() error {
	ss := ts.File.GetSheetList()
	for _, s := range ss {
		// セルの値を再計算
		rows, err := ts.File.GetRows(s)
		if err != nil {
			return err
		}
		for r, row := range rows {
			for c := range row {
				axis, _ := excelize.CoordinatesToCellName(r+1, c+1)
				ts.File.CalcCellValue(s, axis)
			}
		}
		ts.File.GetCellValue(s, "A1")
	}
	ts.File.UpdateLinkedValue()
	ts.File.SetActiveSheet(0)
	err := ts.File.Save()
	if err != nil {
		return fmt.Errorf("Cloud not update file: %v)", err)
	}
	return nil
}

// GetSetCellValFunc セルに値を設定する関数を提供
func (ts *TestSpec) GetSetCellValFunc(sheet string) func(x, y int, val interface{}) error {
	setSheetName(sheet, ts.File)

	return func(x, y int, val interface{}) error {
		axis, err := excelize.CoordinatesToCellName(x, y)
		if err != nil {
			return err
		}
		err = ts.File.SetCellValue(sheet, axis, val)
		if err != nil {
			return err
		}
		return nil
	}
}

// GetMergeCellFunc セルを結合する関数を提供
func (ts *TestSpec) GetMergeCellFunc(sheet string) func(x1, y1, x2, y2 int) error {
	setSheetName(sheet, ts.File)

	return func(x1, y1, x2, y2 int) error {
		axis1, err := excelize.CoordinatesToCellName(x1, y1)
		if err != nil {
			return err
		}

		axis2, err := excelize.CoordinatesToCellName(x2, y2)
		if err != nil {
			return err
		}

		err = ts.File.MergeCell(sheet, axis1, axis2)
		if err != nil {
			return err
		}

		return nil
	}
}

// GetSetColWidthFunc カラムの幅を設定する関数を提供
func (ts *TestSpec) GetSetColWidthFunc(sheet string) func(col string, width int) error {
	setSheetName(sheet, ts.File)

	return func(col string, width int) error {
		w := float64(width) * 2.5
		return ts.File.SetColWidth(sheet, col, col, w)
	}
}

// GetSetRowHeightFunc 行の高さを設定する関数を提供
func (ts *TestSpec) GetSetRowHeightFunc(sheet string) func(srow, erow, steps int) error {
	setSheetName(sheet, ts.File)

	return func(srow, erow, steps int) error {
		diff := erow - srow + 1
		h := float64((rowHeightBase * steps / diff) - steps)
		// 最低限の高さを確保
		if h < rowHeightBase {
			h = rowHeightBase
		}
		for row := srow; row <= erow; row++ {
			ts.File.SetRowHeight(sheet, row, h)
		}
		return nil
	}
}

// GetSetStyleFunc セルのスタイルを設定する関数を提供
func (ts *TestSpec) GetSetStyleFunc(sheet string) func(axis1, axis2 string) error {
	setSheetName(sheet, ts.File)

	return func(axis1, axis2 string) error {
		style, err := ts.File.NewStyle(cellLineStyle)
		if err != nil {
			return err
		}
		err = ts.File.SetCellStyle(sheet, axis1, axis2, style)
		if err != nil {
			return err
		}
		return nil
	}
}

func setSheetName(name string, f *excelize.File) {
	ss := f.GetSheetList()
	isExist := false
	n := 0
	for i, s := range ss {
		if s == name {
			isExist = true
			n = i
			break
		}
	}
	if !isExist {
		n = f.NewSheet(name)
	}
	f.SetActiveSheet(n)
}

func newExcelFile(filename string) (*excelize.File, error) {
	f := excelize.NewFile()
	f.Path = filename
	err := f.Save()
	if err != nil {
		return nil, fmt.Errorf("Cloud not Save file %v: %v", filename, err)
	}
	return f, nil
}
