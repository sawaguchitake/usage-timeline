package reader

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
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
		begin := parseDateFlexible(row[2])
		if begin.IsZero() {
			return nil, fmt.Errorf("invalid begin date: %s", row[2])
		}
		end := parseDateFlexible(row[3]) // 終了日は空欄許容
		records = append(records, UsageRecord{
			No:          i,
			EquipmentID: row[0],
			User:        row[1],
			BeginDate:   begin,
			EndDate:     end,
		})
	}
	return records, nil
}

// 1桁月・日にも対応した日付パース関数
func parseDateFlexible(s string) time.Time {
	if s == "" {
		return time.Time{}
	}

	layouts := []string{"2006-1-2", "2006-01-02", "2006/1/2", "2006/01/02"}
	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, s)
		if err == nil {
			return t
		}
	}
	return time.Time{}
}
