package model

import "github.com/penndev/galite/pkg/orm"

// StartLog 表示启动记录表
type AnalyzeStartLog struct {
	orm.ModelBase
	App                 string            `json:"app"`                    // 应用名
	Channel             string            `json:"channel"`                // 渠道码
	DeviceID            string            `json:"device_id"`              // 设备ID
	AnalyzeDeviceInfoID uint              `json:"analyze_device_info_id"` // 设备信息ID外键
	AnalyzeDeviceInfo   AnalyzeDeviceInfo // 设备信息关联
	Referer             string            `json:"referer"`                                  // 来源网址
	IP                  uint32            `json:"ip" gorm:"type:int unsigned;comment:ipv4"` // IP信息，存储为IPv4数字
}
