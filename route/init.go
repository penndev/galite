package route

import (
	"github.com/penndev/galite/admin"
	"github.com/penndev/galite/config"
	"github.com/penndev/galite/route/middle"

	"github.com/gin-gonic/gin"
)

// 开发模式与正常模式
func Init() *gin.Engine {
	var engine *gin.Engine
	if config.Mode == config.ModeDEV {
		engine = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
		engine.Use(middle.Logger(config.Logger))
		engine.Use(middle.Recovery(config.Logger))
	}
	// 处理通用的中间件
	engine.Use(middle.CORS())

	// 后台请求路由
	admin.InitRoute(engine.Group("/admin"))

	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})

	return engine
}
