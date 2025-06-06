package model

import "github.com/penndev/galite/model/orm"

// 访问的设备信息存在大量的冗余拆分独立建表
type AnalyzeDeviceInfo struct {
	orm.Model
	Md5          []byte `json:"md5" gorm:"type:binary(16);unique;not null"`
	UserAgent    string `json:"userAgent"`
	GPUVendor    string `json:"gpuVendor"`
	GPUModel     string `json:"gpuModel"`
	ScreenWidth  int    `json:"screenWidth"`
	ScreenHeight int    `json:"screenHeight"`
	CPUCores     int    `json:"cpuCores"`
	MemoryGB     int    `json:"memoryGb"`
	OSPlatform   string `json:"osPlatform"`
	OSLanguage   string `json:"osLanguage"`
	OSTimezone   string `json:"osTimezone"`
}
