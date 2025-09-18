package api

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/penndev/galite/internal/wafcdn/model"
	"github.com/penndev/gopkg/captcha2"
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
	host := c.Query("host")
	if host == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "host not found"})
		return
	}
	wildcardHost, _ := publicsuffix.EffectiveTLDPlusOne(host)
	domain := model.Domain{}
	result := domain.DB().Where("domain = ? or (domain = ? and wildcard = true)", host, wildcardHost).Preload("Site").Last(&domain)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
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

	c.JSON(200, gin.H{
		"site":     domain.SiteID,
		"sslForce": domain.SSLForce,
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
			"cache_purge":       domain.Site.Proxy.CachePurge,
		},
	})
}

// 对nginx提供接口 获取验证码
// @url=/@wafcdn/captcha?host=@host
// @return 验证码配置
func HandleCaptcha(c *gin.Context) {

	img := &captcha2.NewDragImg{
		ImageWidth:  300,
		ImageHeight: 150,
	}
	img.DragDraw()

	bufImage := new(bytes.Buffer)
	if err := png.Encode(bufImage, img.Image); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bufPiece := new(bytes.Buffer)
	if err := png.Encode(bufPiece, img.Piece); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New().String()
	c.JSON(http.StatusOK, gin.H{
		"id":          id,
		"imageBase64": "data:image/png;base64," + base64.StdEncoding.EncodeToString(bufImage.Bytes()),
		"imageWidth":  img.ImageWidth,
		"imageHeight": img.ImageHeight,
		"pieceWidth":  img.PieceWidth,
		"pieceHeight": img.PieceHeight,
		"pieceBase64": "data:image/png;base64," + base64.StdEncoding.EncodeToString(bufPiece.Bytes()),
		"verifyX":     img.PieceX,
		"verifyY":     img.PieceY,
	})

}
