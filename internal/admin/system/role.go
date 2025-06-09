package system

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/admin/bind"
	"github.com/penndev/galite/internal/admin/model/system"
)

type RoleRoute struct {
	*gin.RouterGroup
	list []system.RouteItem
}

func (r *RoleRoute) GET(relativePath string, handlers ...gin.HandlerFunc) {
	r.list = append(r.list, system.RouteItem{Method: http.MethodGet, Path: r.BasePath() + relativePath})
	r.RouterGroup.GET(relativePath, handlers...)
}

func (r *RoleRoute) POST(relativePath string, handlers ...gin.HandlerFunc) {
	r.list = append(r.list, system.RouteItem{Method: http.MethodPost, Path: r.BasePath() + relativePath})
	r.RouterGroup.POST(relativePath, handlers...)
}

func (r *RoleRoute) DELETE(relativePath string, handlers ...gin.HandlerFunc) {
	r.list = append(r.list, system.RouteItem{Method: http.MethodDelete, Path: r.BasePath() + relativePath})
	r.RouterGroup.DELETE(relativePath, handlers...)
}

func (r *RoleRoute) PUT(relativePath string, handlers ...gin.HandlerFunc) {
	r.list = append(r.list, system.RouteItem{Method: http.MethodPut, Path: r.BasePath() + relativePath})
	r.RouterGroup.PUT(relativePath, handlers...)
}

func (r *RoleRoute) GETRoutes(c *gin.Context) {
	c.JSON(http.StatusOK, bind.DataList{Data: r.list})
}

// 通过对 gin router.Group 进行封装，来控制全部的路由信息
func NewRoleRouter(r *gin.RouterGroup, middleware ...gin.HandlerFunc) *RoleRoute {
	r.Use(middleware...) // 使用角色鉴权中间件
	return &RoleRoute{
		RouterGroup: r,
		list:        []system.RouteItem{},
	}
}

func RoleList(c *gin.Context) {
	param := &bindSystemRoleParam{}
	if err := c.BindQuery(&param); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	var total int64
	var list []system.SysRole

	m := param.Param() //处理筛选
	m.List(&total, &list)
	c.JSON(http.StatusOK, bind.DataList{Total: total, Data: list})
}

func RoleAdd(c *gin.Context) {
	param := &system.SysRole{}
	if err := c.BindJSON(&param); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	if err := param.Bind(param).Create(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "创建失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.ErrorMessage{Message: "完成"})
	}
}

func RoleUpdate(c *gin.Context) {
	param := &system.SysRole{}
	if err := c.BindJSON(&param); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	if err := param.Bind(param).Updates(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "更新失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.ErrorMessage{Message: "完成"})
	}
}

func RoleDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	param := &system.SysRole{}
	param.ID = uint(id)
	if err := param.Bind(param).Delete(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "删除失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.ErrorMessage{Message: "完成"})
	}
}
