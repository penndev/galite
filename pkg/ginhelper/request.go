package ginhelper

import (
	"github.com/gin-gonic/gin"
)

// @return Context.[Keys].RequestHost 请求Host
// example:
//
//	https://example.com:443 输出的示例
//
// 非信任获取
func RequestHost() gin.HandlerFunc {
	return func(c *gin.Context) {
		scheme := c.Request.Header.Get("X-Forwarded-Proto")
		if scheme == "" {
			scheme = "http"
		}

		host := c.Request.Header.Get("X-Forwarded-Host")
		if host == "" {
			host = c.Request.Host
		}

		port := c.Request.Header.Get("X-Forwarded-Port")
		if port != "" {
			host = host + ":" + port
		}

		uri := scheme + "://" + host // 固定Data目录
		c.Set("requestHost", uri)
	}
}
