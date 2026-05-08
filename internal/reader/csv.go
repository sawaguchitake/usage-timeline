package reader

import (
	"encoding/csv"
	"fmt"
	"os"
)

// FromCSV は指定されたCSVファイルから機器使用記録を読み込み、UsageRecordのスライスを返します。
func FromCSV(filename string, _ Options) (records []UsageRecord, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for i, row := range rows {
		if i == 0 {
			continue // ヘッダー行（id,user_name,begin_date,end_date）をスキップ
		}
		begin, err := fromMultiFormatDate(row[2])
		if err != nil {
			return nil, fmt.Errorf("invalid begin date on row %d: %w", i+1, err)
		}
		end, err := fromMultiFormatDate(row[3]) // 終了日は空欄許容
		if err != nil {
			return nil, fmt.Errorf("invalid end date on row %d: %w", i+1, err)
		}
		if !validateDatePeriod(begin, end) {
			return nil, fmt.Errorf("invalid date period on row %d: begin=%v, end=%v", i+1, toDateString(begin), toDateString(end))
		}
		records = append(records, UsageRecord{
			No:          i + 1,
			EquipmentID: row[0],
			User:        row[1],
			BeginDate:   begin,
			EndDate:     end,
		})
	}
	return records, nil
}
