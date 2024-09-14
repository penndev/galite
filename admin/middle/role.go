package middle

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/admin/bind"
	"github.com/penndev/galite/model/system"
)

// 用户菜单鉴权
func Role(isLog bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		admin, err := system.SysAdminGetByID(c.GetString("jwtAuth"))
		if err != nil {
			c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "用户鉴权失败(1)"})
			c.Abort()
			return
		}
		if admin.Status != 1 {
			c.JSON(http.StatusUnauthorized, bind.ErrorMessage{Message: "用户状态错误"})
			c.Abort()
			return
		}
		if *admin.SysRoleID != 0 && admin.SysRole.Status != 1 {
			c.JSON(http.StatusUnauthorized, bind.ErrorMessage{Message: "角色状态错误"})
			c.Abort()
			return
		}

		// 没设置权限则默认为超级管理员
		if *admin.SysRoleID > 0 {
			routes := admin.SysRole.Route
			pass := false
			for _, route := range routes {
				if route.Method == c.Request.Method && route.Path == c.Request.URL.Path {
					pass = true
					break
				}
			}
			if !pass {
				c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "用户鉴权失败(2)"})
				c.Abort()
				return
			}
		}

		if isLog { // 记录日志
			access := &system.SysAccessLog{
				SysAdminID: admin.ID,
				Method:     c.Request.Method,
				Path:       fmt.Sprint(c.Request.URL),
				IP:         c.ClientIP(),
			}
			if err := access.Bind(access).Create(access).Error; err != nil {
				c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "日志记录失败:" + err.Error()})
				c.Abort()
				return
			}
			c.Next()
			httpRequest, _ := httputil.DumpRequest(c.Request, false)
			access.Payload = string(httpRequest)
			access.Status = c.Writer.Status()
			access.Bind(access).Updates(access)
		}
	}
}
