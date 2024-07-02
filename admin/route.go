package admin

import (
	"github.com/penndev/wga/admin/middle"
	"github.com/penndev/wga/admin/system"
	"github.com/penndev/wga/config"

	"github.com/gin-gonic/gin"
)

func InitAdminRoute(r *gin.RouterGroup) {
	r.GET("/captcha", system.Captcha) //后台登录验证码
	r.POST("/login", system.Login)    //后台登录验证码

	r.Use(middle.JWTAuth(config.JWTSecret))
	r.GET("/system/admin", system.SystemUserList)
}
