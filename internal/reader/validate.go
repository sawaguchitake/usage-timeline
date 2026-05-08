package reader

import (
	"time"
)

// validateDatePeriod は、開始日と終了日の期間が有効かどうかを検証します。
// 開始日または終了日が空の場合は有効とみなします。
// 開始日が終了日より後の場合は無効とみなします。
// 開始日と終了日が同じ年・月でない場合は無効とみなします。
func validateDatePeriod(start, end time.Time) bool {
	if start.IsZero() || end.IsZero() {
		return true
	}
	if start.After(end) {
		return false
	}
	return start.Year() == end.Year() && start.Month() == end.Month()
}
