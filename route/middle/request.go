package middle

import (
	"github.com/gin-gonic/gin"
)

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

		uri := scheme + "://" + host // 固定Data目录。
		c.Set("requestHost", uri)
	}
}
