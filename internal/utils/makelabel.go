package utils

import (
	"time"
)

// MakeLabels は指定された期間の各日の「日」ラベルと「曜日」ラベルのスライスを生成します。
func MakeLabels(minDate time.Time, maxDate time.Time) ([]string, []string) {
	dateLabels := []string{}
	weekLabels := []string{}
	for d := minDate; !d.After(maxDate); d = d.AddDate(0, 0, 1) {
		dateLabels = append(dateLabels, d.Format("02"))
		// 曜日を2文字（例: Mo, Tu, We, Th, Fr, Sa, Su）で取得
		w := d.Weekday()
		var w2 string
		switch w {
		case time.Monday:
			w2 = "Mo"
		case time.Tuesday:
			w2 = "Tu"
		case time.Wednesday:
			w2 = "We"
		case time.Thursday:
			w2 = "Th"
		case time.Friday:
			w2 = "Fr"
		case time.Saturday:
			w2 = "Sa"
		case time.Sunday:
			w2 = "Su"
		}
		weekLabels = append(weekLabels, w2)
	}
	return dateLabels, weekLabels
}
