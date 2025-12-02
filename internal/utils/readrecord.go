package utils

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sawaguchitake/usage-timeline/internal/reader"
)

// ReadRecord は指定されたファイルから使用記録を読み込み、UsageRecordのスライスを返します。
// opts はオプションを可変長で受け取ります。渡されなければデフォルトの Options{} を使用します。
func ReadRecord(file string, opts ...reader.Options) ([]reader.UsageRecord, error) {
	var options reader.Options
	if len(opts) > 0 {
		options = opts[0]
	}

	ext := strings.ToLower(filepath.Ext(file))
	var records []reader.UsageRecord
	var err error

	switch ext {
	case ".csv":
		records, err = reader.FromCSV(file, options)
	case ".xlsx":
		records, err = reader.FromExcel(file, options)
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
