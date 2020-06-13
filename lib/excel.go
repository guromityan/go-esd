package lib

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/360EntSecGroup-Skylar/excelize"
)

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

type TestSpec struct {
	Data Tests
	File *excelize.File
}

func NewTestSpec(data *Tests) (*TestSpec, error) {
	filename := fmt.Sprintf("esd-%v.xlsx", getNow())
	if data.path == "" {
		data.path, _ = os.Getwd()
	}
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

func (ts *TestSpec) GetSetColWidthFunc(sheet string) func(col string, width int) error {
	setSheetName(sheet, ts.File)

	return func(col string, width int) error {
		w := float64(width) * 2.5
		return ts.File.SetColWidth(sheet, col, col, w)
	}
}

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

func getNow() string {
	// t := time.Now()
	// return t.Format("20060102150405")
	return "testfile"
}
