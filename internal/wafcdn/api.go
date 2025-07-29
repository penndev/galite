package wafcdn

import (
	"github.com/penndev/galite/internal/wafcdn/api"
	"github.com/penndev/galite/pkg/ginhelper"
)

// 注册api接口
// wafcdn.InitApiRoute(engine.Group("/@wafcdn"))
func InitApiRoute(r *ginhelper.RouterGroup) {
	r.GET("/ssl", api.HandleSSL)       // 获取证书配置信息
	r.GET("/domain", api.HandleDomain) // 配置域名信息
	r.GET("/ipcheck", api.HandleIpCheck)
	r.GET("/cache", api.HandleGetCache) // 查询缓存信息
	r.PUT("/cache", api.HandlePutCache) // 缓存文件成功报告
	r.PUT("/log", api.HandlePutLog)     // 日志提交
}
