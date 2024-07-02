package middle

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")

		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Max-Age", "259200")
			c.Header("Access-Control-Allow-Methods", "*")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "*")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
