package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/sawaguchitake/usage-timeline/internal/reader"
)

var (
	// 色判定用: 土曜=青, 日曜=赤, 他=デフォルト
	colorReset = "\x1b[0m"
	colorBlue  = "\x1b[34m"
	colorRed   = "\x1b[31m"
)

func main() {
	csvFile := "usage.csv"
	if len(os.Args) > 1 {
		csvFile = os.Args[1]
	}
	records, err := reader.FromCSV(csvFile)
	if err != nil {
		log.Fatalf("process csv: %v", err)
	}

	if len(records) == 0 {
		log.Fatalf("no records found in CSV file: %s", csvFile)
	}

	// ガントチャートの期間（全機器の使用期間の最小・最大）を決定
	minDate, maxDate := records[0].BeginDate, records[0].EndDate
	for _, u := range records {
		if u.BeginDate.Before(minDate) {
			minDate = u.BeginDate
		}
		if u.EndDate.After(maxDate) {
			maxDate = u.EndDate
		}
	}

	// 機器ID昇順、同一IDの場合は使用開始日昇順、開始日が同じ場合は終了日が早い順でソート
	sort.Slice(records, func(i, j int) bool {
		if records[i].EquipmentID == records[j].EquipmentID {
			if records[i].BeginDate.Equal(records[j].BeginDate) {
				return records[i].EndDate.Before(records[j].EndDate)
			}
			return records[i].BeginDate.Before(records[j].BeginDate)
		}
		return records[i].EquipmentID < records[j].EquipmentID
	})

	// 日付ラベル（日のみ）を縦方向に表示
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

	nameCol := "User Name"
	nameWidth := 10

	// 日付ラベル
	fmt.Print(padName(nameCol, nameWidth) + "| ")
	for i, label := range dateLabels {
		w := weekLabels[i]
		color := getWeekColor(w)
		fmt.Print(color, label, colorReset, " ")
	}
	fmt.Println()

	// 曜日ラベル行
	fmt.Print(padName("", nameWidth) + "| ")
	for _, w := range weekLabels {
		color := getWeekColor(w)
		fmt.Print(color, w, colorReset, " ")
	}
	fmt.Println()

	// 各機器の使用記録ガントチャート表示（使用者名単位）
	prevID := ""
	for _, u := range records {
		if u.EquipmentID != prevID {
			// 機器IDが変わったらセパレータ行を挿入
			fmt.Println(strings.Repeat("-", nameWidth) + "+-" + strings.Repeat("---", len(weekLabels)))
		}
		fmt.Print(padName(u.User, nameWidth) + "| ")
		idx := 0
		isEndless := u.EndDate.IsZero()
		for d := minDate; !d.After(maxDate); d = d.AddDate(0, 0, 1) {
			w := weekLabels[idx]
			color := getWeekColor(w)
			if !d.Before(u.BeginDate) {
				if isEndless {
					fmt.Print(color, "??", colorReset, " ")
				} else if !d.After(u.EndDate) {
					fmt.Print(color, "**", colorReset, " ")
				} else {
					fmt.Print(color, "  ", colorReset, " ")
				}
			} else {
				fmt.Print(color, "  ", colorReset, " ")
			}
			idx++
		}
		fmt.Println()
		prevID = u.EquipmentID
	}
}

// 全角文字を考慮して指定幅でパディングする関数
func padName(name string, width int) string {
	w := 0
	for _, r := range name {
		if r >= 0x2E80 && r <= 0x9FFF { // CJK統合漢字・ひらがな・カタカナ等
			w += 2
		} else {
			w += 1
		}
	}
	pad := width - w
	if pad > 0 {
		return name + strings.Repeat(" ", pad)
	}
	return name
}

// 曜日ラベルに対応する色を返す関数
func getWeekColor(w string) string {
	switch w {
	case "Sa":
		return colorBlue
	case "Su":
		return colorRed
	default:
		return colorReset
	}
}
