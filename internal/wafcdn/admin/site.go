package admin

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/admin/bind"
	"github.com/penndev/galite/internal/lib"
	"github.com/penndev/galite/internal/wafcdn/model"
)

func ensureSiteStaticDir(site *model.Site) error {
	if site.Type != "static" {
		return nil
	}
	root := strings.Trim(strings.ReplaceAll(strings.TrimSpace(site.Static.Root), "\\", "/"), "/")
	if root == "" {
		return errors.New("请填写静态站点目录名")
	}
	site.Static.Root = root
	if err := normalizeStaticLocations(site); err != nil {
		return err
	}
	abs, err := FilesResolve(root)
	if err != nil {
		return err
	}
	return os.MkdirAll(abs, 0755)
}

func normalizeStaticLocations(site *model.Site) error {
	seen := map[string]struct{}{}
	for i := range site.Static.Locations {
		loc := &site.Static.Locations[i]
		path := strings.TrimSpace(loc.Path)
		if path == "" {
			return errors.New("路径不能为空")
		}
		if !strings.HasPrefix(path, "/") {
			return errors.New("路径必须以 / 开头")
		}
		path = strings.TrimRight(path, "/")
		if path == "" || strings.Contains(path, "..") || strings.Contains(path, "//") || strings.Contains(path, "?") || strings.Contains(path, "#") || strings.Contains(path, "\\") {
			return errors.New("路径不合法")
		}
		if _, ok := seen[path]; ok {
			return errors.New("路径重复: " + path)
		}
		seen[path] = struct{}{}
		loc.Path = path

		server, err := normalizeProxyServer(loc.Proxy.Server)
		if err != nil {
			return errors.New("路径 " + path + " 的" + err.Error())
		}
		loc.Proxy.Server = server
	}
	return nil
}

func normalizeProxyServer(server string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(server))
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Hostname() == "" || parsed.Port() == "" {
		return "", errors.New("反代地址必须是 http(s)://主机:端口，路径可选")
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", errors.New("反代地址不能带参数")
	}
	path := strings.TrimRight(parsed.EscapedPath(), "/")
	if path == "/" || path == "" {
		path = ""
	} else if !strings.HasPrefix(path, "/") || strings.Contains(path, "..") {
		return "", errors.New("反代路径不合法")
	}
	return parsed.Scheme + "://" + parsed.Host + path, nil
}

func normalizeSiteProxy(site *model.Site) error {
	if site.Type != "proxy" {
		return nil
	}
	server, err := normalizeProxyServer(site.Proxy.Server)
	if err != nil {
		return err
	}
	site.Proxy.Server = server
	return nil
}

// 添加新的站点

func SiteAdd(c *gin.Context) {
	param := &model.Site{}
	if err := c.BindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}
	if err := ensureSiteStaticDir(param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if err := normalizeSiteProxy(param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if err := param.DB().Create(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "创建失败(" + err.Error() + ")"})
		return
	}
	c.JSON(http.StatusOK, bind.Message{Message: "完成"})
}

// 获取站点列表
func SiteList(c *gin.Context) {
	param := &bindSiteParam{}
	if err := c.BindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}
	var total int64
	var list []model.Site

	m := param.Param() //处理筛选
	m.List(&total, &list)
	c.JSON(http.StatusOK, bind.DataList{Total: total, Data: list})
}

// 更新资料
func SiteUpdate(c *gin.Context) {
	param := &model.Site{}
	if err := c.BindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}
	if err := ensureSiteStaticDir(param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if err := normalizeSiteProxy(param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if err := param.DB().Updates(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "更新失败(" + err.Error() + ")"})
		return
	}
	c.JSON(http.StatusOK, bind.Message{Message: "完成"})
}

// 删除资料
func SiteDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}
	param := &model.Site{}
	param.ID = uint(id)
	if err := param.DB().Delete(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "删除失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.Message{Message: "完成"})
	}
}

// 获取 国家与省份
func IPRegion(c *gin.Context) {
	if lib.IPRegion == nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "IPREGION_DB 未配置或文件不存在"})
		return
	}
	var region model.IPRegion
	c.BindQuery(&region)

	parentID := uint32(0)
	names := make([]string, 0, 3)
	if region.Country != "" {
		names = append(names, region.Country)
	}
	if region.Province != "" {
		names = append(names, region.Province)
	}
	if region.City != "" {
		names = append(names, region.City)
	}
	for _, name := range names {
		found := false
		for _, item := range lib.IPRegion.Areas(parentID) {
			if item.Name == name {
				parentID = item.ID
				found = true
				break
			}
		}
		if !found {
			c.JSON(http.StatusOK, bind.DataList{Data: []string{}})
			return
		}
	}

	children := lib.IPRegion.Areas(parentID)
	data := make([]string, 0, len(children))
	for _, item := range children {
		data = append(data, item.Name)
	}
	c.JSON(http.StatusOK, bind.DataList{Data: data})
}
