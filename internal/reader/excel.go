package reader

import (
	"fmt"
	"os"
	"sort"
	"time"

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
		fmt.Println("ファイルを開けません:", err)
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
	fmt.Println(sheetName)

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("no rows: %w", err)
	}

	parseDate := func(dateStr string) (time.Time, error) {
		if dateStr == "" {
			return time.Time{}, nil
		}
		return time.Parse("01-02-06", dateStr)
	}

	// 8行以降でB列に値がある全行を出力
	for i, row := range rows[7:] {
		if len(row) > 1 && row[1] != "" {
			startDate, err := parseDate(row[3])
			if err != nil {
				return nil, fmt.Errorf("parse start date on row %d: %w", i+8, err)
			}
			endDate, err := parseDate(row[4])
			if err != nil {
				return nil, fmt.Errorf("parse end date on row %d: %w", i+8, err)
			}
			notes := ""
			if len(row) > 7 {
				notes = row[7]
			}
			record := UsageRecord{
				No:          i + 8,
				EquipmentID: row[1],
				User:        row[2],
				BeginDate:   startDate,
				EndDate:     endDate,
				TargetUser:  row[5],
				Purpose:     row[6],
				Notes:       notes,
			}
			records = append(records, record)
		}
	}

	return records, nil
}
