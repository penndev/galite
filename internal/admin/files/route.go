package files

import (
	"github.com/penndev/galite/pkg/ginhelper"
)

// InitRoute 注册系统文件管理接口。
func InitRoute(route *ginhelper.RoleRoute) {
	route.GET("/list", handleList)
	route.GET("/content", handleGetContent)
	route.GET("/download", handleDownload)

	route.POST("/dir", handleMkdir)
	route.POST("/upload", handleUpload)

	route.PUT("/content", handlePutContent)
	route.PUT("/rename", handleRename)
	route.PUT("/chmod", handleChmod)

	route.DELETE("/file", handleDelete)
}
