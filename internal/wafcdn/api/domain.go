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
)

// 对nginx提供接口 获取证书配置
// @url=/@wafcdn/domain?host=@host
// @return 配置信息
func findDomain(host string) (model.Domain, error) {
	domain := model.Domain{}
	err := domain.DB().Where("domain = ?", host).Preload("Site").First(&domain).Error
	if err == nil || !strings.Contains(host, ".") || strings.HasPrefix(host, "*.") {
		return domain, err
	}
	rest := host[strings.Index(host, ".")+1:]
	if !strings.Contains(rest, ".") {
		return domain, err
	}
	wild := "*." + rest
	domain = model.Domain{}
	err = domain.DB().Where("domain = ?", wild).Preload("Site").First(&domain).Error
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
	cert, ok := model.FindCertForHost(host)
	if !ok {
		c.JSON(400, gin.H{"error": "ssl not found"})
		return
	}
	c.JSON(200, gin.H{
		"publicKey":  cert.PublicKey,
		"privateKey": cert.PrivateKey,
	})
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
	if domain.Site.Root == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "site root not set"})
		return
	}
	abs, err := admin.FilesResolve(domain.Site.Root)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	locations := make([]gin.H, 0)
	for _, loc := range domain.Site.Locations {
		item := gin.H{
			"path":   loc.Path,
			"action": loc.Action,
		}
		if loc.Action == "proxy" {
			if loc.Proxy == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "路径 " + loc.Path + " 缺少回源配置"})
				return
			}
			locProxy, err := proxyPayload(*loc.Proxy, true)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "替换域名为IP失败: " + err.Error()})
				return
			}
			item["proxy"] = locProxy
		} else if loc.Action == "redirect" {
			item["redirect"] = loc.Redirect
			item["code"] = loc.Code
		}
		locations = append(locations, item)
	}

	c.JSON(200, gin.H{
		"site":      domain.SiteID,
		"sslForce":  domain.SSLForce,
		"security":  domain.Site.Security,
		"header":    respHeader,
		"root":      filepath.ToSlash(abs),
		"locations": locations,
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
