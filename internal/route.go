package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/admin"
	"github.com/penndev/galite/internal/config"
	"github.com/penndev/galite/internal/logger"
	"github.com/penndev/galite/internal/wafcdn"
	"github.com/penndev/galite/pkg/middle"
)

// 开发模式与正常模式
func InitRoute() *gin.Engine {
	var engine *gin.Engine
	if config.Mode() == config.ModeDev {
		engine = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
		engine.Use(middle.Logger(logger.GinZapLogger))
		engine.Use(middle.Recovery(logger.GinZapLogger))
	}
	// 处理通用的中间件
	engine.Use(middle.CORS())

	// 后台请求路由
	admin.InitRoute(engine.Group("/admin"))

	wafcdn.InitApiRoute(engine.Group("/@wafcdn"))

	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})

	return engine
}
