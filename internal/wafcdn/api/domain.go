package api

import (
	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/wafcdn/model"
	"golang.org/x/net/publicsuffix"
)

// 对nginx提供接口 获取证书配置
// @url=/@wafcdn/domain?host=@host
// @return 配置信息
func HandleSSL(c *gin.Context) {
	host := c.Query("host")
	if host == "" {
		c.JSON(400, gin.H{
			"error": "host not found",
		})
		return
	}
	domain := model.Domain{}
	wildcardHost, _ := publicsuffix.EffectiveTLDPlusOne(host)
	result := domain.DB().Where("domain = ? or (domain = ? and wildcard = true)", host, wildcardHost).Last(&domain)
	if !domain.SSL || result.Error != nil {
		c.JSON(400, gin.H{
			"error": "ssl not found",
		})
		return
	}
	c.JSON(200, gin.H{
		"publicKey":  domain.PublicKey,
		"privateKey": domain.PrivateKey,
	})
}

// 对nginx提供接口 获取域名配置信息
// @url=/@wafcdn/domain?host=@host
// @return 配置信息
func HandleDomain(c *gin.Context) {
	if c.Query("host") == "" {
		c.JSON(400, gin.H{
			"error": "host not found",
		})
		return
	}

	host := c.Query("host")
	wildcardHost, _ := publicsuffix.EffectiveTLDPlusOne(host)
	domain := model.Domain{}
	result := domain.DB().Where("domain = ? or (domain = ? and wildcard = true)", host, wildcardHost).Preload("Site").Last(&domain)
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

	sslForce := false
	if domain.SSL {
		sslForce = domain.SSLForce
	}

	c.JSON(200, gin.H{
		"site":     domain.SiteID,
		"sslForce": sslForce,
		"type":     domain.Site.Type,
		"security": domain.Site.Security,
		"header":   respHeader,
		"proxy": gin.H{
			"server":            domain.Site.Proxy.Server,
			"host":              domain.Site.Proxy.Host,
			"keepaliveTimeout":  domain.Site.Proxy.KeepaliveTimeout,
			"keepaliveRequests": domain.Site.Proxy.KeepaliveRequests,
			"header":            proxyHeader,
			"cache":             domain.Site.Proxy.Cache,
		},
	})
}
