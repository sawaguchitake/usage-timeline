package utils

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sawaguchitake/usage-timeline/internal/reader"
)

// ReadRecord は指定されたファイルから使用記録を読み込み、UsageRecordのスライスを返します。
func ReadRecord(file string) ([]reader.UsageRecord, error) {
	ext := strings.ToLower(filepath.Ext(file))
	var records []reader.UsageRecord
	var err error

	switch ext {
	case ".csv":
		records, err = reader.FromCSV(file)
	case ".xlsx":
		records, err = reader.FromExcel(file)
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("process file: %v", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no records found in file: %s", file)
	}
	return records, nil
}
