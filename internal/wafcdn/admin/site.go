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

func ensureSiteDir(site *model.Site) error {
	root := strings.Trim(strings.ReplaceAll(strings.TrimSpace(site.Root), "\\", "/"), "/")
	if root == "" {
		return errors.New("请填写站点目录")
	}
	site.Root = root
	if err := normalizeLocations(site); err != nil {
		return err
	}
	abs, err := FilesResolve(root)
	if err != nil {
		return err
	}
	return os.MkdirAll(abs, 0755)
}

func normalizeLocations(site *model.Site) error {
	seen := map[string]struct{}{}
	for i := range site.Locations {
		loc := &site.Locations[i]
		path := strings.TrimSpace(loc.Path)
		if path != "/" {
			path = strings.TrimRight(path, "/")
		}
		if path == "" || !strings.HasPrefix(path, "/") || strings.Contains(path, "..") || strings.Contains(path, "//") || strings.Contains(path, "?") || strings.Contains(path, "#") || strings.Contains(path, "\\") {
			return errors.New("路径不合法")
		}
		if _, ok := seen[path]; ok {
			return errors.New("路径重复: " + path)
		}
		seen[path] = struct{}{}
		loc.Path = path

		switch loc.Action {
		case "proxy":
			if loc.Proxy == nil {
				return errors.New("路径 " + path + " 的反代地址必须是 http(s)://主机:端口，路径可选")
			}
			server, err := normalizeProxyServer(loc.Proxy.Server)
			if err != nil {
				return errors.New("路径 " + path + " 的" + err.Error())
			}
			loc.Proxy.Server = server
			loc.Redirect = ""
			loc.Code = 0
		case "redirect":
			target, err := normalizeRedirect(loc.Redirect)
			if err != nil {
				return errors.New("路径 " + path + " 的" + err.Error())
			}
			loc.Redirect = target
			if loc.Code != 301 {
				loc.Code = 302
			}
			loc.Proxy = nil
		default:
			return errors.New("路径 " + path + " 的处理方式不正确")
		}
	}
	return nil
}

func normalizeRedirect(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "/") && !strings.HasPrefix(raw, "//") {
		if strings.Contains(raw, "..") || strings.Contains(raw, "\\") {
			return "", errors.New("重定向路径不合法")
		}
		return raw, nil
	}
	parsed, err := url.Parse(raw)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return "", errors.New("重定向地址必须是 http(s)://主机 或 /路径")
	}
	return raw, nil
}

func normalizeProxyServer(server string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(server))
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Hostname() == "" || parsed.Port() == "" {
		return "", errors.New("反代地址必须是 http(s)://主机:端口，路径可选")
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", errors.New("反代地址不能带参数")
	}
	path := parsed.EscapedPath()
	if path == "/" {
		// 单独的 / 表示换掉站点路径前缀
	} else if path != "" {
		path = strings.TrimRight(path, "/")
		if path == "" || !strings.HasPrefix(path, "/") || strings.Contains(path, "..") {
			return "", errors.New("反代路径不合法")
		}
	}
	return parsed.Scheme + "://" + parsed.Host + path, nil
}

// 添加新的站点

func SiteAdd(c *gin.Context) {
	param := &model.Site{}
	if err := c.BindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}
	if err := ensureSiteDir(param); err != nil {
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
	if err := ensureSiteDir(param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if err := param.DB().Model(param).Select("remark", "root", "locations", "security", "header").Updates(param).Error; err != nil {
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
