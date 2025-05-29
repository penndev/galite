package wafcdn

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/wafcdn/model"
)

// 对nginx提供接口 获取证书配置
// @url=/@wafcdn/domain?host=@host
// @return 配置信息
func handleSSL(c *gin.Context) {
	host := c.Query("host")
	if host == "" {
		c.JSON(400, gin.H{
			"error": "host not found",
		})
		return
	}
	domain := model.Domain{}
	result := domain.Bind(&domain).Where("name = ?", host).Last(&domain)
	if !domain.SSL || result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"publickey":  domain.PublicKey,
		"privatekey": domain.PrivateKey,
	})
}

// 对nginx提供接口 获取域名配置信息
// @url=/@wafcdn/domain?host=@host
// @return 配置信息
func handleDomain(c *gin.Context) {
	if c.Query("host") == "" {
		c.JSON(400, gin.H{
			"error": "host not found",
		})
		return
	}
	domain := model.Domain{}
	result := domain.Bind(&domain).Where("name = ?", c.Query("host")).Preload("Site").Last(&domain)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	// 格式返回头
	respHeader := make(map[string]string)
	for _, item := range domain.Site.Header {
		respHeader[item.Name] = item.Value
	}

	proxyHeader := make(map[string]string)
	for _, item := range domain.Site.Proxy.Header {
		proxyHeader[item.Name] = item.Value
	}

	sslforce := false
	if domain.SSL {
		sslforce = domain.SSLForce
	}
	c.JSON(200, gin.H{
		"site":     domain.SiteID,
		"sslforce": sslforce,
		"type":     domain.Site.Type,
		"security": domain.Site.Security,
		"header":   respHeader,
		"proxy": gin.H{
			"server":             domain.Site.Proxy.Server,
			"host":               domain.Site.Proxy.Host,
			"keepalive_timeout":  domain.Site.Proxy.KeepaliveTimeout,
			"keepalive_requests": domain.Site.Proxy.KeepaliveRequests,
			"header":             proxyHeader,
			"cache":              domain.Site.Proxy.Cache,
		},
	})
}

func handleGetCache(c *gin.Context) {
	param := &model.Cache{}
	if err := c.BindQuery(param); err != nil {
		log.Println("参数错误", err.Error())
		c.JSON(400, gin.H{"message": "参数错误" + err.Error()})
		return
	}
	if err := param.Bind(param).Where("site_id = ? and method = ? and uri = ?", param.SiteID, param.Method, param.Uri).First(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "查询失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, param)
	}
}

func handlePutCache(c *gin.Context) {
	param := &model.Cache{}
	if err := c.BindJSON(param); err != nil {
		log.Println("参数错误", err.Error())
		c.JSON(400, gin.H{"message": "参数错误" + err.Error()})
		return
	}
	// 不存在则创建
	if err := param.Bind(param).Where("site_id = ? and method = ? and uri = ?", param.SiteID, param.Method, param.Uri).Assign(*param).FirstOrCreate(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "创建失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "完成"})
	}
}
