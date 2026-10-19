package api

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/penndev/galite/internal/wafcdn/admin"
	"github.com/penndev/galite/internal/wafcdn/model"
	"github.com/penndev/gopkg/captcha/drag"
	"golang.org/x/net/publicsuffix"
)

// 对nginx提供接口 获取证书配置
// @url=/@wafcdn/domain?host=@host
// @return 配置信息
func findDomain(host string) (model.Domain, error) {
	domain := model.Domain{}
	err := domain.DB().Where("domain = ?", host).Preload("Site").Preload("Cert").First(&domain).Error
	if err == nil || !strings.Contains(host, ".") || strings.HasPrefix(host, "*.") {
		return domain, err
	}
	rest := host[strings.Index(host, ".")+1:]
	if !strings.Contains(rest, ".") {
		return domain, err
	}
	wild := "*." + rest
	domain = model.Domain{}
	err = domain.DB().Where("domain = ?", wild).Preload("Site").Preload("Cert").First(&domain).Error
	return domain, err
}

func HandleSSL(c *gin.Context) {
	host := c.Query("host")
	if host == "" {
		c.JSON(400, gin.H{
			"error": "host not found",
		})
		return
	}
	cert := model.Cert{}
	domain, err := findDomain(host)
	if err == nil && !domain.SSL {
		c.JSON(400, gin.H{"error": "ssl not found"})
		return
	}
	if err == nil && domain.Cert.PublicKey != "" && domain.Cert.PrivateKey != "" {
		cert = domain.Cert
	} else if wild, ok := findWildcardCert(host); ok {
		cert = wild
	} else {
		c.JSON(400, gin.H{"error": "ssl not found"})
		return
	}
	c.JSON(200, gin.H{
		"publicKey":  cert.PublicKey,
		"privateKey": cert.PrivateKey,
	})
}

func findWildcardCert(host string) (model.Cert, bool) {
	apex, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil || apex == "" || host == apex {
		return model.Cert{}, false
	}
	prefix := strings.TrimSuffix(host, "."+apex)
	if prefix == "" || strings.Contains(prefix, ".") {
		return model.Cert{}, false
	}
	cert := model.Cert{}
	err = cert.DB().Model(&cert).Where("wildcard = ? AND domain = ?", true, apex).First(&cert).Error
	if err != nil || cert.PublicKey == "" || cert.PrivateKey == "" {
		return model.Cert{}, false
	}
	return cert, true
}

// 对nginx提供接口 获取域名配置信息
// @url=/@wafcdn/domain?host=@host
// @return 配置信息
func HandleDomain(c *gin.Context) {
	host := c.Query("host")
	if host == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "host not found"})
		return
	}
	domain, err := findDomain(host)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 格式返回头
	respHeader := make(map[string]string)
	for _, item := range domain.Site.Header {
		respHeader[item.Name] = item.Value
	}
	proxyBody, err := proxyPayload(domain.Site.Proxy, domain.Site.Type == "proxy")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "替换域名为IP失败: " + err.Error()})
		return
	}

	staticRoot := ""
	locations := make([]gin.H, 0)
	if domain.Site.Type == "static" {
		if domain.Site.Static.Root == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "static root not set"})
			return
		}
		abs, err := admin.FilesResolve(domain.Site.Static.Root)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		staticRoot = filepath.ToSlash(abs)
		for _, loc := range domain.Site.Static.Locations {
			locProxy, err := proxyPayload(loc.Proxy, true)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "替换域名为IP失败: " + err.Error()})
				return
			}
			locations = append(locations, gin.H{
				"path":  loc.Path,
				"proxy": locProxy,
			})
		}
	}

	c.JSON(200, gin.H{
		"site":     domain.SiteID,
		"sslForce": domain.SSLForce,
		"type":     domain.Site.Type,
		"security": domain.Site.Security,
		"header":   respHeader,
		"static": gin.H{
			"root":      staticRoot,
			"locations": locations,
		},
		"proxy": proxyBody,
	})
}

func proxyPayload(proxy model.SiteProxy, resolveIP bool) (gin.H, error) {
	header := make(map[string]string)
	for _, item := range proxy.Header {
		header[item.Name] = item.Value
	}
	server := proxy.Server
	if resolveIP && server != "" {
		replaced, err := replaceDomainWithIP(server)
		if err != nil {
			return nil, err
		}
		server = replaced
	}
	return gin.H{
		"server":            server,
		"host":              proxy.Host,
		"keepaliveTimeout":  proxy.KeepaliveTimeout,
		"keepaliveRequests": proxy.KeepaliveRequests,
		"header":            header,
		"cache":             proxy.Cache,
		"cache_purge":       proxy.CachePurge,
	}, nil
}

// 对nginx提供接口 获取验证码
// @url=/@wafcdn/captcha?host=@host
// @return 验证码配置
func HandleCaptcha(c *gin.Context) {
	data, point, err := drag.New()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := uuid.New().String()
	c.JSON(http.StatusOK, gin.H{
		"id":          id,
		"imageBase64": data.ImageBase64,
		"imageWidth":  data.ImageWidth,
		"imageHeight": data.ImageHeight,
		"pieceWidth":  data.PieceWidth,
		"pieceHeight": data.PieceHeight,
		"pieceBase64": data.PieceBase64,
		"verifyX":     point.X,
		"verifyY":     point.Y,
	})
}
