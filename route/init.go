package route

import (
	"github.com/penndev/wga/admin"
	"github.com/penndev/wga/config"
	"github.com/penndev/wga/route/middle"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	engine := gin.New()

	// 全局组件....
	engine.Use(middle.Logger(config.Logger))
	engine.Use(middle.Recovery(config.Logger))

	//
	engine.Use(middle.CORS())

	admin.InitAdminRoute(engine.Group("/admin"))

	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})

	return engine
}
