package ginhelper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouterGroup struct {
	*gin.RouterGroup
}

// 进行语法糖封装可以进行快捷下级路由组装
func (e *RouterGroup) GroupPush(relativePath string, r func(*RouterGroup)) {
	rg := &RouterGroup{RouterGroup: e.Group(relativePath)}
	r(rg)
}

type Engine struct {
	*gin.Engine
}

// 进行语法糖封装可以进行快捷下级路由组装
func (e *Engine) GroupPush(relativePath string, r func(*RouterGroup)) {
	rg := &RouterGroup{RouterGroup: e.Group(relativePath)}
	r(rg)
}

type RouteItem struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

// 进行路由拦截，对同一级路由下的部分路由做管理。
// 验证权限  - 后台验证是否有权限访问等。
type RoleRoute struct {
	*RouterGroup
	list   []RouteItem
	listRG []*RoleRoute // 用于存储路由列表
}

func (r *RoleRoute) GET(relativePath string, handlers ...gin.HandlerFunc) {
	r.list = append(r.list, RouteItem{Method: http.MethodGet, Path: r.BasePath() + relativePath})
	r.RouterGroup.GET(relativePath, handlers...)
}

func (r *RoleRoute) POST(relativePath string, handlers ...gin.HandlerFunc) {
	r.list = append(r.list, RouteItem{Method: http.MethodPost, Path: r.BasePath() + relativePath})
	r.RouterGroup.POST(relativePath, handlers...)
}

func (r *RoleRoute) DELETE(relativePath string, handlers ...gin.HandlerFunc) {
	r.list = append(r.list, RouteItem{Method: http.MethodDelete, Path: r.BasePath() + relativePath})
	r.RouterGroup.DELETE(relativePath, handlers...)
}

func (r *RoleRoute) PUT(relativePath string, handlers ...gin.HandlerFunc) {
	r.list = append(r.list, RouteItem{Method: http.MethodPut, Path: r.BasePath() + relativePath})
	r.RouterGroup.PUT(relativePath, handlers...)
}

func (r *RoleRoute) RouteList() []RouteItem {
	var routes []RouteItem
	routes = append(routes, r.list...)
	for _, sub := range r.listRG {
		routes = append(routes, sub.RouteList()...)
	}
	return routes
}

// 进行语法糖封装可以进行快捷下级路由组装
func (r *RoleRoute) GroupPush(relativePath string, rr func(*RoleRoute)) {
	rg := &RoleRoute{
		RouterGroup: &RouterGroup{RouterGroup: r.Group(relativePath)},
		list:        []RouteItem{},
		listRG:      []*RoleRoute{},
	}
	r.listRG = append(rg.listRG, rg)
	rr(rg)
}

// 通过对 gin router.Group 进行封装，来控制全部的路由信息
func (rg *RouterGroup) RoleRoute(middleware ...gin.HandlerFunc) *RoleRoute {
	rg.Use(middleware...) // 使用角色鉴权中间件
	return &RoleRoute{
		RouterGroup: rg,
		list:        []RouteItem{},
		listRG:      []*RoleRoute{},
	}
}
