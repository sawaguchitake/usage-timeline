package utils

import (
	"sort"

	"github.com/sawaguchitake/usage-timeline/internal/reader"
)

// SortRecords はUsageRecordのスライスを機器ID昇順、同一IDの場合は使用開始日昇順、開始日が同じ場合は終了日が早い順でソートします。
func SortRecords(records []reader.UsageRecord) {
	sort.Slice(records, func(i, j int) bool {
		if records[i].EquipmentID == records[j].EquipmentID {
			if records[i].BeginDate.Equal(records[j].BeginDate) {
				return records[i].EndDate.Before(records[j].EndDate)
			}
			return records[i].BeginDate.Before(records[j].BeginDate)
		}
		return records[i].EquipmentID < records[j].EquipmentID
	})
}
