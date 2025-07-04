package middle

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/admin/bind"
	"github.com/penndev/galite/internal/admin/model/system"
)

// 用户菜单鉴权
func Role() gin.HandlerFunc {
	return func(c *gin.Context) {
		admin, err := system.SysAdminGetByID(c.GetInt("adminID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, bind.Message{Message: "用户鉴权失败(1)"})
			c.Abort()
			return
		}
		if admin.Status != 1 {
			c.JSON(http.StatusUnauthorized, bind.Message{Message: "用户状态错误"})
			c.Abort()
			return
		}
		if admin.SysRoleID != nil && *admin.SysRoleID != 0 && admin.SysRole.Status != 1 {
			c.JSON(http.StatusUnauthorized, bind.Message{Message: "角色状态错误"})
			c.Abort()
			return
		}

		// 没设置权限则默认为超级管理员
		if admin.SysRoleID != nil && *admin.SysRoleID > 0 {
			routes := admin.SysRole.Route
			pass := false
			for _, route := range routes {
				if route.Method == c.Request.Method && route.Path == c.Request.URL.Path {
					pass = true
					break
				}
			}
			if !pass {
				c.JSON(http.StatusBadRequest, bind.Message{Message: "用户鉴权失败(2)"})
				c.Abort()
				return
			}
		}

		access := &system.SysAccessLog{
			SysAdminID: admin.ID,
			Method:     c.Request.Method,
			Path:       fmt.Sprint(c.Request.URL),
			IP:         c.ClientIP(),
		}
		c.Set("accessLog", true) // 设置访问日志标志
		c.Next()
		if c.GetBool("accessLog") {
			// 是否记录访问日志
			// httpRequest, err := httputil.DumpRequest(c.Request, false)
			// if err != nil {
			// 	access.Payload = "日志记录失败: " + err.Error()
			// } else {
			// 	access.Payload = string(httpRequest)
			// }
			access.Status = c.Writer.Status()
			if err := access.Bind(access).Create(access).Error; err != nil {
				c.JSON(http.StatusBadRequest, bind.Message{Message: "日志记录失败:" + err.Error()})
				c.Abort()
				return
			}
		}

	}
}
