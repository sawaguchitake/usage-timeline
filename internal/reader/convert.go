package reader

import (
	"fmt"
	"time"
)

// fromMultiFormatDate は複数のフォーマットに対応して日付文字列を解析します。
func fromMultiFormatDate(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}

	layouts := []string{"2006-1-2", "2006-01-02", "2006/1/2", "2006/01/02"}
	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("failed to parse date with any supported format: %w", err)
}

// fromExcelDate はExcelの日付形式を解析してtime.Timeを返します。
func fromExcelDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse("01-02-06", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// toDateString はtime.Timeを"YYYY-MM-DD"形式の文字列に変換します。ゼロ値の場合は空文字を返します。
func toDateString(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format("2006-01-02")
}
