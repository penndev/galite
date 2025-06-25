package wafcdn

import (
	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/wafcdn/admin"
	"github.com/penndev/galite/internal/wafcdn/api"
)

// 注册api接口
// wafcdn.InitApiRoute(engine.Group("/@wafcdn"))
func InitApiRoute(r *gin.RouterGroup) {
	r.GET("/ssl", handleSSL)            // 获取证书配置信息
	r.GET("/domain", handleDomain)      // 配置域名信息
	r.GET("/cache", api.HandleGetCache) // 查询缓存信息
	r.PUT("/cache", api.HandlePutCache) // 缓存文件成功报告
}

// 注册admin接口
// wafcdn.InitAdminRoute(engine.Group("/wafcdn"))
func InitAdminRoute(route *gin.RouterGroup) {
	// wafcdn
	route.GET("/site", admin.SiteList)
	route.POST("/site", admin.SiteAdd)
	route.PUT("/site", admin.SiteUpdate)
	route.DELETE("/site", admin.SiteDelete)
	// 管理域名
	route.GET("/domain", admin.DomainList)
	route.POST("/domain", admin.DomainAdd)
	route.PUT("/domain", admin.DomainUpdate)
	route.DELETE("/domain", admin.DomainDelete)

	// 查看缓存列表
	route.GET("/cache", admin.CacheList)
	route.DELETE("/cache", admin.CacheDelete)
}
