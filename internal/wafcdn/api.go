package wafcdn

import (
	"github.com/penndev/galite/internal/wafcdn/api"
	"github.com/penndev/galite/pkg/ginhelper"
)

// 注册api接口
// wafcdn.InitApiRoute(engine.Group("/@wafcdn"))
func InitApiRoute(r *ginhelper.RouterGroup) {
	r.GET("/acme", api.HandleAcme)          // 申请证书
	r.GET("/purge", api.HandlePurgeCache)   // 用户主动刷新缓存
	r.GET("/ssl", api.HandleSSL)            // 获取证书配置信息
	r.GET("/domain", api.HandleDomain)      // 配置域名信息
	r.GET("/ip-verify", api.HandleIpVerify) // IP校验
	r.PUT("/cache", api.HandlePutCache)     // 缓存文件成功报告
	r.PUT("/log", api.HandlePutLog)         // 批量日志提交
}
