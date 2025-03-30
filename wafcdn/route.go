package wafcdn

import (
	"github.com/gin-gonic/gin"
)

func InitRoute(r *gin.RouterGroup) {
	r.GET("/ssl", handleSSL)
	r.GET("/domain", handleDomain)
	r.GET("/cache", handleGetCache) // 查询是否存在缓存
	r.PUT("/cache", handlePutCache) // 创建缓存
}
