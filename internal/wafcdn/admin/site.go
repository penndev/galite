package admin

import (
	"errors"
	"net/http"
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
	abs, err := FilesResolve(root)
	if err != nil {
		return err
	}
	return os.MkdirAll(abs, 0755)
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
	if err := param.DB().Create(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "创建失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.Message{Message: "完成"})
	}
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
	if err := param.DB().Updates(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "更新失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.Message{Message: "完成"})
	}
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
