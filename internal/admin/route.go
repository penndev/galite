package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/admin/bind"
	"github.com/penndev/galite/internal/admin/middle"
	"github.com/penndev/galite/internal/admin/system"
	"github.com/penndev/galite/internal/config"
	"github.com/penndev/galite/internal/admin/files"
	"github.com/penndev/galite/internal/wafcdn"
	"github.com/penndev/galite/pkg/ginhelper"
)

func InitRoute(r *ginhelper.RouterGroup) {
	// 未登录开放接口
	r.GET("/captcha", system.Captcha)     // 获取验证码
	r.POST("/login", system.Login)        // 用户登录验证
	r.POST("/login-2fa", system.LoginOTP) // 用户登录两步验证
	/*
	 * 验证登录状态
	 * r.Context.GetInt("adminID") 获取登录用户的ID
	 */
	r.Use(middle.JWTAuth(config.Secret()))

	r.PUT("/change-passwd", system.ChangePasswd)
	r.PUT("/change-otp", system.ChangeOTP)
	r.GET("/otp/secret", system.GetOTPSecret)    // 创建OTP验证器 （谷歌验证器）
	r.PUT("/otp/secret", system.VerifyOTPSecret) //	校验OTP验证器 （谷歌验证器）
	/*
	 * 权限验证接口
	 * route.Context.Set("accessLog", false) 被middle.Role控制
	 */
	route := r.RoleRoute(middle.Role())

	// 后台脚手架鉴权控制功能
	route.GET("/system/role/route", func(c *gin.Context) {
		c.JSON(http.StatusOK, bind.DataList{Data: route.RouteList()})
	}) // 通过对路由包装来动态返回全接口
	route.GET("/system/role", system.RoleList)
	route.POST("/system/role", system.RoleAdd)
	route.PUT("/system/role", system.RoleUpdate)
	route.DELETE("/system/role", system.RoleDelete)
	route.GET("/system/admin", system.AdminList)
	route.POST("/system/admin", system.AdminAdd)
	route.PUT("/system/admin", system.AdminUpdate)
	route.DELETE("/system/admin", system.AdminDelete)
	route.GET("/system/admin/access-log", system.AdminAccessLog)

	// 挂载wafcdn后台管理
	route.GroupPush("/wafcdn", wafcdn.InitAdminRoute)

	// 文件管理
	route.GroupPush("/files", files.InitRoute)
}
