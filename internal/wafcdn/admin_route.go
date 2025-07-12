package wafcdn

import (
	"github.com/penndev/galite/internal/wafcdn/admin"
	"github.com/penndev/galite/pkg/ginhelper"
)

// 注册admin接口
// wafcdn.InitAdminRoute(engine.Group("/wafcdn"))
func InitAdminRoute(route *ginhelper.RouterGroup) {
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
	route.GET("/cache/delete", admin.CacheList)
	route.GET("/cache/delete/list", admin.CacheDeleteList)
	route.PUT("/cache/delete/list", admin.CacheDeleteAdd)

	// 查看访问日志列表
	route.GET("/log", admin.LogList)

	route.GET("/openrsty/status", admin.OpenrestyStatus)
	route.PUT("/openrsty/start", admin.OpenrestyStart)
	route.PUT("/openrsty/stop", admin.OpenrestyStop)
	route.PUT("/openrsty/reload", admin.OpenrestyReload)
}
