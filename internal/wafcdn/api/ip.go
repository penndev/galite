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
	if len(ipStruct.Region) > 0 {
		if lib.IPRegion == nil {
			c.JSON(http.StatusForbidden, bind.Message{Message: "IPREGION_DB 未配置或文件不存在"})
			return
		}
		info, err := lib.IPRegion.Find(c.Query("ip"))
		if err != nil {
			c.JSON(http.StatusForbidden, bind.Message{Message: err.Error()})
			return
		}
		var names []string
		for a := &info.Area; a != nil && a.Name != ""; a = a.Parent {
			names = append(names, a.Name)
		}
		for i, j := 0, len(names)-1; i < j; i, j = i+1, j-1 {
			names[i], names[j] = names[j], names[i]
		}
		for len(names) < 4 {
			names = append(names, "")
		}
		reqRegion := model.IPRegion{
			Country:  names[0],
			Province: names[1],
			City:     names[2],
			County:   names[3],
			ISP:      info.ISP.Name,
		}
		for _, region := range ipStruct.Region {
			if region.Country != "" && region.Country != reqRegion.Country { // 国家
				continue
			}
			if region.Province != "" && region.Province != reqRegion.Province { // 省份
				continue
			}
			if region.City != "" && region.City != reqRegion.City { // 城市
				continue
			}
			if region.County != "" && region.County != reqRegion.County { // 县级
				continue
			}
			// 如果匹配到允许则是允许拒绝则是拒绝。
			if ipStruct.Allowed { // 如果能匹配到任意一条规则，则根据 Allowed 返回
				c.JSON(http.StatusOK, bind.Message{Message: "ok"})
			} else {
				c.JSON(http.StatusForbidden, gin.H{"region": reqRegion})
			}
			return
		}
	}
	// 如果没有匹配到则允许是拒绝， 拒绝是允许
	if !ipStruct.Allowed { // 如果能匹配到任意一条规则，则根据 Allowed 返回
		c.JSON(http.StatusOK, bind.Message{Message: "ok"})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"region": "deny"})
	}
}
