package admin

import (
	"github.com/penndev/galite/admin/bind"
	"github.com/penndev/galite/admin/middle"
	"github.com/penndev/galite/admin/system"
	"github.com/penndev/galite/config"

	"github.com/gin-gonic/gin"
)

func InitAdminRoute(r *gin.RouterGroup) {
	// 未登录开放接口
	r.GET("/captcha", system.Captcha) // 获取验证码
	r.POST("/login", system.Login)    // 用户登录验证
	/*
	 * 验证登录状态
	 */
	r.Use(middle.JWTAuth(config.JWTSecret))

	r.PUT("/change-passwd", system.ChangePasswd)

	/*
	 * 权限验证接口
	 */
	route := bind.NewRoleRouter(r, middle.Role(true))

	// 后台脚手架鉴权控制功能
	route.GET("/system/role/route", route.GETRoutes) // 通过对路由包装来动态返回全接口
	route.GET("/system/role", system.RoleList)
	route.POST("/system/role", system.RoleAdd)
	route.PUT("/system/role", system.RoleUpdate)
	route.DELETE("/system/role", system.RoleDelete)
	route.GET("/system/admin", system.AdminList)
	route.POST("/system/admin", system.AdminAdd)
	route.PUT("/system/admin", system.AdminUpdate)
	route.DELETE("/system/admin", system.AdminDelete)
	route.GET("/system/admin/access-log", system.AdminAccessLog)
}
