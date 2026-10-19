package admin

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/admin/bind"
	"github.com/penndev/galite/internal/wafcdn/model"
)

// 添加新的站点
func DomainAdd(c *gin.Context) {
	param := &model.Domain{}
	if err := c.BindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}
	param.Domain = strings.TrimSpace(param.Domain)
	if param.Domain == "" {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "请填写域名"})
		return
	}
	if param.SiteID == nil || *param.SiteID == 0 {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "站点不存在"})
		return
	}
	if err := (&model.Site{}).DB().First(&model.Site{}, *param.SiteID).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "站点不存在"})
		return
	}
	if param.SSLForce {
		if _, ok := model.FindCertForHost(param.Domain); !ok {
			c.JSON(http.StatusBadRequest, bind.Message{Message: "没有覆盖该域名的证书，不能开启自动跳转 HTTPS"})
			return
		}
	}
	if err := param.DB().Omit("Site").Save(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "存储失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.Message{Message: "完成"})
	}
}

// 获取站点列表
func DomainList(c *gin.Context) {
	param := &bindDomainParam{}
	if err := c.BindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}
	var total int64
	var list []model.Domain

	m := param.Param() //处理筛选
	m.List(&total, &list)
	c.JSON(http.StatusOK, bind.DataList{Total: total, Data: list})
}

// 删除资料
func DomainDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}
	param := &model.Domain{}
	param.ID = uint(id)
	if err := param.DB().Delete(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "删除失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.Message{Message: "完成"})
	}
}
