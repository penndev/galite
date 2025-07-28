package admin

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/admin/bind"
	"github.com/penndev/galite/internal/wafcdn/model"
	"github.com/penndev/gopkg/ip2region"
)

// 添加新的站点

func SiteAdd(c *gin.Context) {
	param := &model.Site{}
	if err := c.BindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}
	if err := param.Bind(param).Create(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "创建失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.Message{Message: "完成"})
	}
}

// 获取站点列表
func SiteList(c *gin.Context) {
	param := &bindSiteParam{}
	if err := c.BindQuery(&param); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误"})
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
		log.Println(err)
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误"})
		return
	}
	if err := param.Bind(param).Updates(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "更新失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.Message{Message: "完成"})
	}
}

// 删除资料
func SiteDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误"})
		return
	}
	param := &model.Site{}
	param.ID = uint(id)
	if err := param.Bind(param).Delete(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "删除失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.Message{Message: "完成"})
	}
}

// 获取 国家与省份
func IPRegion(c *gin.Context) {
	var region ip2region.IPRegion
	c.BindQuery(&region)

	// 筛选数据列表
	var current *ip2region.Region
	var data []string // 返回数据内容

	// 如果未选择国家，返回所有国家名称
	if region.Country == "" {
		for _, item := range ip2region.RegionList {
			data = append(data, item.Name)
		}
		c.JSON(http.StatusOK, bind.DataList{Data: data})
		return
	} else {
		for _, item := range ip2region.RegionList {
			if region.Country == item.Name {
				current = &item
				break
			}
		}
	}

	// 如果未选择省份，返回国家下的所有省份
	if region.Province == "" {
		for _, province := range current.Children {
			data = append(data, province.Name)
		}
		c.JSON(http.StatusOK, bind.DataList{Data: data})
		return
	} else {
		for _, item := range current.Children {
			if region.Province == item.Name {
				current = &item
				break
			}
		}
	}

	// 如果未选择城市，返回省份下的所有城市
	if region.City == "" {
		for _, child := range current.Children {
			data = append(data, child.Name)
		}
		c.JSON(http.StatusOK, bind.DataList{Data: data})
		return
	} else {
		for _, item := range current.Children {
			if region.City == item.Name {
				current = &item
				break
			}
		}
	}

	for _, child := range current.Children {
		data = append(data, child.Name)
	}
	c.JSON(http.StatusOK, bind.DataList{Data: data})
}
