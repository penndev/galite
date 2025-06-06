package analyze

import (
	"github.com/gin-gonic/gin"
)

func InitApiRoute(r *gin.RouterGroup) {
	r.GET("/start", Start) // 应用启动 网页打开获取网站配置
	r.GET("/action")       // 操作动作 上报下载操作
}
