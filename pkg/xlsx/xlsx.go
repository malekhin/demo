package xlsx

import (
	"demo/internal/common"
	"fmt"

	"github.com/xuri/excelize/v2"
)

const DefaultSheet = "Sheet1"

type Xlsx struct {
	file   *excelize.File
	sheets []*Sheet
}

func New() *Xlsx {
	return &Xlsx{file: excelize.NewFile()}
}

func (x *Xlsx) AddSheet(sheetName string, isActive bool) (*Sheet, error) {
	sheetName = cut(sheetName)

	sheetIndex, err := x.file.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("%w: NewSheet", err)
	}
	if isActive {
		x.file.SetActiveSheet(sheetIndex)
	}

	sheet := &Sheet{
		name: sheetName,
	}
	x.sheets = append(x.sheets, sheet)

	return sheet, nil
}

func (x *Xlsx) Build() (*Buffer, error) {
	for _, sheet := range x.sheets {
		colIndex := 1
		for _, cellValues := range sheet.data {
			for cellIndex, cellValue := range cellValues {
				colName, err := excelize.ColumnNumberToName(cellIndex + 1)
				if err != nil {
					return nil, common.Wrap(err, "excelize.ColumnNumberToName")
				}
				cell := fmt.Sprintf("%s%d", colName, colIndex)
				err = x.file.SetCellValue(sheet.name, cell, cellValue)
				if err != nil {
					return nil, common.Wrap(err, "SetCellValue")
				}
			}
			colIndex++
		}
	}

	// This function will be invalid when only one worksheet is left.
	err := x.file.DeleteSheet(DefaultSheet)
	if err != nil {
		return nil, fmt.Errorf("%w: DeleteSheet", err)
	}

	var b Buffer
	if err := x.file.Write(&b); err != nil {
		return nil, fmt.Errorf("%w: Write", err)
	}

	return &b, nil
}

func cut(s string) string {
	r := []rune(s)
	if len(r) > 31 {
		s = string(r[:31])
		return s
	}
	return s
}
