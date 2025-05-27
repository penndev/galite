package route

import (
	"crypto/sha256"
	"os"

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

	// 处理签名中间件
	signMiddle := middle.Signature(middle.SignatureConfig{
		Key:         []byte(os.Getenv("APP_SECRET")), // 签名密钥
		Hash:        sha256.New,                      // HMAC算法
		SignName:    "sign",                          // 签名参数名称
		ExpiredName: "expired",                       // 过期时间参数名称
		// 自定义异常处理结果
		OnError: func(c *gin.Context, status int, err error) {
			c.JSON(status, gin.H{
				"message": err.Error(),
			})
		},
	})

	engine.GET("/ping", signMiddle, func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})

	return engine
}
