package middle

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/admin/bind"
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
	c.JSON(http.StatusOK, bind.DataList{Data: r.list})
}

func (r *RoleRoute) GET(relativePath string, handlers ...gin.HandlerFunc) {
	r.list = append(r.list, RouteItem{Method: "GET", Path: r.BasePath() + relativePath})
	r.RouterGroup.GET(relativePath, handlers...)
}

func (r *RoleRoute) POST(relativePath string, handlers ...gin.HandlerFunc) {
	r.list = append(r.list, RouteItem{Method: "POST", Path: r.BasePath() + relativePath})
	r.RouterGroup.POST(relativePath, handlers...)
}

func (r *RoleRoute) DELETE(relativePath string, handlers ...gin.HandlerFunc) {
	r.list = append(r.list, RouteItem{Method: "DELETE", Path: r.BasePath() + relativePath})
	r.RouterGroup.DELETE(relativePath, handlers...)
}

func (r *RoleRoute) PUT(relativePath string, handlers ...gin.HandlerFunc) {
	r.list = append(r.list, RouteItem{Method: "PUT", Path: r.BasePath() + relativePath})
	r.RouterGroup.PUT(relativePath, handlers...)
}
