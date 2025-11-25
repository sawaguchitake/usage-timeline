package reader

import "time"

// 機器使用記録（使用者単位）を表す構造体
type UsageRecord struct {
	No          int       // 連番
	EquipmentID string    // 機器ID
	User        string    // 使用者
	BeginDate   time.Time // 開始日
	EndDate     time.Time // 終了日
	TargetUser  string    // 対象ユーザ
	Purpose     string    // 使用目的／個人情報の有無
	Notes       string    // 備考
}
