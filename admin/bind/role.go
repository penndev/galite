package bind

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouteItem struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

type RoleRoute struct {
	*gin.RouterGroup
	list []RouteItem
}

func NewRoleRouter(r *gin.RouterGroup, role gin.HandlerFunc) *RoleRoute {
	r.Use(role)
	return &RoleRoute{
		RouterGroup: r,
		list:        []RouteItem{},
	}
}

func (r *RoleRoute) GETRoutes(c *gin.Context) {
	c.JSON(http.StatusOK, DataList{Data: r.list})
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
