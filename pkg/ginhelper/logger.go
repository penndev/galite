package ginhelper

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger(logger *zap.Logger) gin.HandlerFunc {
	loggerLevel := logger.Level()
	if loggerLevel > zap.WarnLevel {
		return func(c *gin.Context) {}
	}
	return func(c *gin.Context) {
		start := time.Now()
		query := c.Request.URL.Path
		c.Next()
		cost := time.Since(start)
		if loggerLevel == zap.WarnLevel && cost > 2*time.Second {
			logger.Warn("gin/warn", zap.Duration("cost", cost),
				zap.String("method", c.Request.Method), zap.String("query", query), zap.Int("status", c.Writer.Status()),
				zap.String("ip", c.ClientIP()), zap.String("user-agent", c.Request.UserAgent()))

		} else {
			logger.Info("gin/info", zap.Duration("cost", cost),
				zap.String("method", c.Request.Method), zap.String("query", query), zap.Int("status", c.Writer.Status()),
				zap.String("ip", c.ClientIP()), zap.String("user-agent", c.Request.UserAgent()))
		}
	}
}
