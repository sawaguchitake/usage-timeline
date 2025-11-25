package utils

import (
	"time"

	"github.com/sawaguchitake/usage-timeline/internal/reader"
)

// GetDatePeriod は使用記録のスライスから、全期間の最小開始日と最大終了日を取得します。
func GetDatePeriod(records []reader.UsageRecord) (time.Time, time.Time) {
	minDate, maxDate := records[0].BeginDate, records[0].EndDate
	for _, u := range records {
		if u.BeginDate.Before(minDate) {
			minDate = u.BeginDate
		}
		if u.EndDate.After(maxDate) {
			maxDate = u.EndDate
		}
	}
	return minDate, maxDate
}
