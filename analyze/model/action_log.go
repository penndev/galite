package model

import "github.com/penndev/galite/model/orm"

// StartLog 表示启动记录表
type AnalyzeActionLog struct {
	orm.Model
	Type              string
	Description       string
	AnalyzeStartLogID uint
}
