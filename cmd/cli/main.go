package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"path/filepath"

	"github.com/sawaguchitake/usage-timeline/internal/utils"
	"github.com/sawaguchitake/usage-timeline/internal/reader"
)

var (
	// 色判定用: 土曜=青, 日曜=赤, 他=デフォルト
	colorReset = "\x1b[0m"
	colorBlue  = "\x1b[34m"
	colorRed   = "\x1b[31m"
)

func main() {
	file := "usage.csv"
	sheet := ""
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	if len(os.Args) > 2 {
		sheet = os.Args[2]
	}

	var records []reader.UsageRecord
	var err error
	if sheet != "" && strings.ToLower(filepath.Ext(file)) == ".xlsx" {
		records, err = utils.ReadRecord(file, reader.Options{SheetName: sheet})
	} else {
		records, err = utils.ReadRecord(file)
	}
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	utils.SortRecords(records)

	minDate, maxDate := utils.GetDatePeriod(records)
	dateLabels, weekLabels := utils.MakeLabels(minDate, maxDate)

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
