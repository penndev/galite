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

	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})

	// 对签名与加解密进行测试
	// 对签名进行验证
	// 对请求进行解密
	// 对返回的body进行加密处理
	var onError = func(c *gin.Context, status int, err error) {
		c.JSON(status, gin.H{
			"message": err.Error(),
		})
		c.Abort()
	}
	var signMiddle = middle.Signature(middle.SignatureConfig{
		Key:         []byte(os.Getenv("APP_SECRET")), // 签名密钥
		Hash:        sha256.New,                      // HMAC算法
		SignName:    "sign",                          // 签名参数名称
		ExpiredName: "expired",                       // 过期时间参数名称
		OnError:     onError,                         // 自定义异常处理结果
	})
	var encryptMiddle = middle.Encryption(middle.EncryptionConfig{
		Secret:          os.Getenv("APP_SECRET"),
		IvName:          "X-Iv",
		ContentTypeName: "application/x-buffer",
		OnError:         onError, // 自定义异常处理结果
	})
	engine.POST("/encrypt", signMiddle, encryptMiddle, func(ctx *gin.Context) {
		var req struct {
			Message string `json:"message"`
		}
		_ = ctx.ShouldBindJSON(&req)
		ctx.JSON(200, gin.H{
			"message": req.Message,
		})
	})
	//.

	return engine
}
