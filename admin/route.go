package admin

import (
	"github.com/penndev/galite/admin/middle"
	"github.com/penndev/galite/admin/system"
	"github.com/penndev/galite/config"

	"github.com/gin-gonic/gin"
)

func InitAdminRoute(r *gin.RouterGroup) {
	r.GET("/captcha", system.Captcha) //后台登录验证码
	r.POST("/login", system.Login)    //后台登录验证码

	r.Use(middle.JWTAuth(config.JWTSecret))
	r.PUT("/change-passwd", system.ChangePasswd) //后台登录验证码
	route := middle.NewRoleRouter(r, system.Role(false))
	// 后台脚手架鉴权控制功能
	route.GET("/system/admin", system.AdminList)
	route.POST("/system/admin", system.AdminAdd)
	route.PUT("/system/admin", system.AdminUpdate)
	route.DELETE("/system/admin", system.AdminDelete)
	route.GET("/system/admin/access-log", system.AccessLog)
	route.GET("/system/role", system.RoleList)
	route.POST("/system/role", system.RoleAdd)
	route.PUT("/system/role", system.RoleUpdate)
	route.DELETE("/system/role", system.RoleDelete)
	route.GET("/system/role/route", route.GETRoutes)
}
