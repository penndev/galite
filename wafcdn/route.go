package wafcdn

import (
	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/wafcdn/admin"
)

// 注册api接口
// wafcdn.InitApiRoute(engine.Group("/@wafcdn"))
func InitApiRoute(r *gin.RouterGroup) {
	r.GET("/ssl", handleSSL)
	r.GET("/domain", handleDomain)
	r.GET("/cache", handleGetCache) // 查询是否存在缓存
	r.PUT("/cache", handlePutCache) // 创建缓存
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
}
