package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/admin/bind"
	"github.com/penndev/galite/internal/wafcdn/model"
	"github.com/penndev/galite/pkg/util"
	"github.com/penndev/gopkg/ip2region"
)

func HandlePutLog(c *gin.Context) {
	param := &model.Log{}
	if err := c.BindJSON(param); err != nil {
		log.Println("参数错误", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误" + err.Error()})
		return
	}
	param.Gorm().Create(param)
	c.JSON(http.StatusOK, gin.H{"Message": "完成"})
}

func HandleIpCheck(c *gin.Context) {
	site := c.Query("site")
	siteID, err := util.StrConv[uint](site)
	if err != nil {
		c.JSON(http.StatusNotFound, bind.Message{})
		return
	}
	siteModel, err := model.GetSiteByID(siteID)
	if err != nil {
		c.JSON(http.StatusNotFound, bind.Message{})
		return
	}
	ipStruct := siteModel.Security.Ip
	if !ipStruct.Status {
		c.JSON(http.StatusNotFound, bind.Message{})
		return
	}
	if len(ipStruct.Region) > 0 {
		reqRegion := ip2region.Find(c.Query("ip"))
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
			if ipStruct.Allowed { // 如果能匹配到任意一条规则，则根据 Allowed 返回
				c.JSON(http.StatusOK, bind.Message{Message: "ok"})
			} else {
				c.JSON(http.StatusNotFound, bind.Message{Message: "fail"})
			}
			return
		}
	}

	if !ipStruct.Allowed { // 最终效果需要取反
		c.JSON(http.StatusOK, bind.Message{Message: "ok"})
	} else {
		c.JSON(http.StatusNotFound, bind.Message{Message: "fail"})
	}
}
