package wafcdn

import (
	"github.com/penndev/galite/internal/wafcdn/admin"
	"github.com/penndev/galite/pkg/ginhelper"
)

// 注册admin接口
// wafcdn.InitAdminRoute(engine.Group("/wafcdn"))
func InitAdminRoute(route *ginhelper.RoleRoute) {
	route.GET("/site", admin.SiteList)
	route.POST("/site", admin.SiteAdd)
	route.PUT("/site", admin.SiteUpdate)
	route.DELETE("/site", admin.SiteDelete)
	route.GET("/ip-region", admin.IPRegion)

	// 静态文件管理
	route.GET("/files", admin.FilesList)
	route.POST("/files", admin.FilesMkdir)
	route.PUT("/files", admin.FilesRename)
	route.DELETE("/files", admin.FilesDelete)
	route.GET("/files/file", admin.FilesDownload)
	route.POST("/files/file", admin.FilesUpload)

	// 证书
	route.GET("/cert", admin.CertList)
	route.POST("/cert", admin.CertSave)
	route.PUT("/cert", admin.CertSave)
	route.DELETE("/cert", admin.CertDelete)
	route.POST("/cert/acme", admin.CertAcme)

	// 管理域名
	route.GET("/domain", admin.DomainList)
	route.POST("/domain", admin.DomainAdd)
	route.PUT("/domain", admin.DomainAdd)
	route.DELETE("/domain", admin.DomainDelete)
	// 查看缓存列表
	route.GET("/cache", admin.CacheList)
	route.DELETE("/cache", admin.CacheDelete)

	route.GET("/cache/delete", admin.CacheList)
	route.GET("/cache/delete/list", admin.CacheDeleteList)
	route.PUT("/cache/delete/list", admin.CacheDeleteAdd)
	route.DELETE("/cache-clear", admin.CacheDeleteClear)

	// 查看访问日志列表
	route.GET("/log", admin.LogList)
	route.DELETE("/log-clear", admin.LogClear)
	route.GET("/openresty/status", admin.OpenrestyStatus)
	route.PUT("/openresty/start", admin.OpenrestyStart)
	route.PUT("/openresty/stop", admin.OpenrestyStop)
	route.PUT("/openresty/reload", admin.OpenrestyReload)
}
