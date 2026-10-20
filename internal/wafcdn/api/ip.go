package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/admin/bind"
	"github.com/penndev/galite/internal/lib"
	"github.com/penndev/galite/internal/wafcdn/model"
	"github.com/penndev/galite/pkg/util"
)

func HandleIpVerify(c *gin.Context) {
	site := c.Query("site")
	siteID, err := util.StrConv[uint](site)
	if err != nil {
		c.JSON(http.StatusForbidden, bind.Message{Message: err.Error()})
		return
	}
	siteModel, err := model.GetSiteByID(siteID)
	if err != nil {
		c.JSON(http.StatusForbidden, bind.Message{Message: err.Error()})
		return
	}
	ipStruct := siteModel.Security.Ip
	if !ipStruct.Status {
		c.JSON(http.StatusForbidden, bind.Message{Message: "status not enable"})
		return
	}
	ip := c.Query("ip")
	areaName := ""
	ispName := ""
	if len(ipStruct.Region) > 0 || ipStruct.Allowed {
		if lib.IPRegion == nil {
			if len(ipStruct.Region) > 0 {
				c.JSON(http.StatusForbidden, bind.Message{Message: "IPREGION_DB 未配置或文件不存在"})
				return
			}
		} else if info, err := lib.IPRegion.Find(ip); err != nil {
			if len(ipStruct.Region) > 0 {
				c.JSON(http.StatusForbidden, bind.Message{Message: err.Error()})
				return
			}
		} else {
			areaName = info.Area.Name
			ispName = info.ISP.Name
		}
	}
	if len(ipStruct.Region) > 0 {
		for _, region := range ipStruct.Region {
			leaf := region.Country
			if region.Province != "" {
				leaf = region.Province
			}
			if region.City != "" {
				leaf = region.City
			}
			if region.County != "" {
				leaf = region.County
			}
			if leaf != "" && leaf != areaName {
				continue
			}
			// 如果匹配到允许则是允许拒绝则是拒绝。
			if ipStruct.Allowed { // 如果能匹配到任意一条规则，则根据 Allowed 返回
				c.JSON(http.StatusOK, bind.Message{Message: "ok"})
			} else {
				c.JSON(http.StatusForbidden, gin.H{"region": areaName, "ip": ip, "isp": ispName})
			}
			return
		}
	}
	// 如果没有匹配到则允许是拒绝， 拒绝是允许
	if !ipStruct.Allowed { // 如果能匹配到任意一条规则，则根据 Allowed 返回
		c.JSON(http.StatusOK, bind.Message{Message: "ok"})
	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"region": "deny",
			"ip":     ip,
			"area":   areaName,
			"isp":    ispName,
		})
	}
}
