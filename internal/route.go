package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/admin"
	"github.com/penndev/galite/internal/config"
	"github.com/penndev/galite/internal/logger"
	"github.com/penndev/galite/pkg/ginhelper"
)

// 开发模式与正常模式
func InitRoute() *ginhelper.Engine {
	engine := &ginhelper.Engine{Engine: nil}
	if config.Mode() == config.ModeDev {
		engine.Engine = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		engine.Engine = gin.New()
		engine.Use(ginhelper.Logger(logger.GinZapLogger))
		engine.Use(ginhelper.Recovery(logger.GinZapLogger))
	}
	// 处理通用的中间件
	engine.Use(ginhelper.CORS())
	// 后台请求路由
	engine.GroupPush("/admin", admin.InitRoute)

	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})

	return engine
}
