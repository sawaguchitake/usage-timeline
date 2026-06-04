package reader

import (
	"fmt"
	"os"
	"sort"

	"github.com/xuri/excelize/v2"
)

// FromExcel は指定されたExcelファイルから機器使用記録を読み込み、UsageRecordのスライスを返します。
// 1. 全シートのうち、シート名を降順ソートした際の最初のシートを選択。
// 2. そのシートの8行以降でB列に値がある全行を出力。
// 3. Record構造体のスライスに格納。
// エラーが発生した場合はエラーを返します。
func FromExcel(filename string, options Options) ([]UsageRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("open excel: %w", err)
	}
	defer f.Close()

	var sheetName = options.SheetName
	if sheetName == "" {
		sheetNames := f.GetSheetList()
		if len(sheetNames) == 0 {
			return nil, fmt.Errorf("no sheets in the workbook")
		}

		// シート名を降順ソートして最初のシートを選択
		sort.Sort(sort.Reverse(sort.StringSlice(sheetNames)))
		sheetName = sheetNames[0]
	}

	return getRecords(sheetName, f)
}

// GetSheetList は指定されたExcelファイルからシート名のリストを取得し、
// シート名を降順ソートしたスライスを返します。
// エラーが発生した場合はエラーを返します。
func GetSheetList(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("open excel: %w", err)
	}
	defer f.Close()

	sheetNames := f.GetSheetList()
	if len(sheetNames) == 0 {
		return nil, fmt.Errorf("no sheets in the workbook")
	}

	// シート名を降順ソート
	sort.Sort(sort.Reverse(sort.StringSlice(sheetNames)))

	return sheetNames, nil
}

// getRecords は指定されたシート名とExcelファイルからUsageRecordのスライスを取得します。
// シートの8行以降でB列に値がある全行を処理し、UsageRecordに格納します。
// エラーが発生した場合はエラーを返します。
func getRecords(sheetName string, f *excelize.File) (records []UsageRecord, err error) {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("no rows: %w", err)
	}

	// 明細は8行目以降。それより行数が少ないシートは明細なしとして扱う。
	if len(rows) <= 7 {
		return nil, nil
	}

	// 8行以降でB列に値がある全行を出力
	for i, row := range rows[7:] {
		if cell(row, 1) != "" {
			beginDate, err := fromExcelDate(cell(row, 3))
			if err != nil {
				return nil, fmt.Errorf("invalid begin date on row %d: %w", i+8, err)
			}
			endDate, err := fromExcelDate(cell(row, 4))
			if err != nil {
				return nil, fmt.Errorf("invalid end date on row %d: %w", i+8, err)
			}
			if !validateDatePeriod(beginDate, endDate) {
				return nil, fmt.Errorf("invalid date period on row %d: begin=%v, end=%v", i+8, toDateString(beginDate), toDateString(endDate))
			}
			record := UsageRecord{
				No:          i + 8,
				EquipmentID: cell(row, 1),
				User:        cell(row, 2),
				BeginDate:   beginDate,
				EndDate:     endDate,
				TargetUser:  cell(row, 5),
				Purpose:     cell(row, 6),
				Notes:       cell(row, 7),
			}
			records = append(records, record)
		}
	}

	return records, nil
}

// cell は row の index 位置のセル値を返す。範囲外は空文字を返す。
// excelize の GetRows は行末の空セルを切り詰めるため、列アクセスは必ずこれを経由する。
func cell(row []string, i int) string {
	if i < len(row) {
		return row[i]
	}
	return ""
}
