package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"time"
)

// 機器使用記録（使用者単位）を表す構造体
type UsageRecord struct {
	EquipmentID string
	UserName    string
	BeginDate   time.Time
	EndDate     time.Time
}

func main() {
	csvFile := "usage.csv"
	if len(os.Args) > 1 {
		csvFile = os.Args[1]
	}
	file, err := os.Open(csvFile)
	if err != nil {
		fmt.Println("ファイルを開けません:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("CSV読み込みエラー:", err)
		return
	}

	var usages []UsageRecord
	for i, record := range records {
		if i == 0 {
			continue // ヘッダー行（id,user_name,begin_date,end_date）をスキップ
		}
		begin, err := parseDateFlexible(record[2])
		var end time.Time
		var err2 error
		if record[3] == "" {
			end = time.Time{} // 空欄の場合はゼロ値（未返却など）
		} else {
			end, err2 = parseDateFlexible(record[3])
			if err != nil || err2 != nil {
				fmt.Printf("日付パースエラー: %v, %v\n", err, err2)
				continue
			}
		}
		usages = append(usages, UsageRecord{
			EquipmentID: record[0],
			UserName:    record[1],
			BeginDate:   begin,
			EndDate:     end,
		})
	}

	// ガントチャートの期間（全機器の使用期間の最小・最大）を決定
	minDate, maxDate := usages[0].BeginDate, usages[0].EndDate
	for _, u := range usages {
		if u.BeginDate.Before(minDate) {
			minDate = u.BeginDate
		}
		if u.EndDate.After(maxDate) {
			maxDate = u.EndDate
		}
	}

	// 機器ID昇順、同一IDの場合は使用開始日昇順でソート
	sort.Slice(usages, func(i, j int) bool {
		if usages[i].EquipmentID == usages[j].EquipmentID {
			return usages[i].BeginDate.Before(usages[j].BeginDate)
		}
		return usages[i].EquipmentID < usages[j].EquipmentID
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
	// 色判定用: 土曜=青, 日曜=赤, 他=デフォルト
	colorReset := "\x1b[0m"
	colorBlue := "\x1b[34m"
	colorRed := "\x1b[31m"

	// 日付ラベル
	fmt.Print(padName(nameCol, nameWidth) + "| ")
	for i, label := range dateLabels {
		w := weekLabels[i]
		var color string
		switch w {
		case "Sa":
			color = colorBlue
		case "Su":
			color = colorRed
		default:
			color = colorReset
		}
		fmt.Print(color, label, colorReset, " ")
	}
	fmt.Println()

	// 曜日ラベル行
	fmt.Print(padName("", nameWidth) + "| ")
	for _, w := range weekLabels {
		var color string
		switch w {
		case "Sa":
			color = colorBlue
		case "Su":
			color = colorRed
		default:
			color = colorReset
		}
		fmt.Print(color, w, colorReset, " ")
	}
	fmt.Println()

	// 各機器の使用記録ガントチャート表示（使用者名単位）
	for _, u := range usages {
		fmt.Print(padName(u.UserName, nameWidth) + "| ")
		idx := 0
		isEndless := u.EndDate.IsZero()
		for d := minDate; !d.After(maxDate); d = d.AddDate(0, 0, 1) {
			w := weekLabels[idx]
			var color string
			if w == "Sa" {
				color = colorBlue
			} else if w == "Su" {
				color = colorRed
			} else {
				color = colorReset
			}
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
		return name + fmt.Sprintf(fmt.Sprintf("%%%ds", pad), " ")
	}
	return name
}

// 1桁月・日にも対応した日付パース関数
func parseDateFlexible(s string) (time.Time, error) {
	layouts := []string{
		"2006-1-2", "2006-01-02", "2006/1/2", "2006/01/02",
	}
	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}
	return t, err
}
