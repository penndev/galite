package admin

import (
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/admin/bind"
	"github.com/penndev/galite/internal/lib"
	"github.com/penndev/galite/internal/wafcdn/model"
	"github.com/penndev/gopkg/acme"
	"golang.org/x/net/publicsuffix"
)

// 添加新的站点
func DomainAdd(c *gin.Context) {
	param := &model.Domain{}
	if err := c.BindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}
	if param.Wildcard {
		var err error
		if param.Domain, err = publicsuffix.EffectiveTLDPlusOne(param.Domain); err != nil {
			c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
			return
		}
	}
	if err := param.ParseCertInfo(); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}
	if err := param.DB().Save(param).Error; err != nil {
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
	// 解析证书
	for i := range list {
		if list[i].PublicKey != "" && list[i].PrivateKey != "" {
			list[i].ParseCertInfo() // 解析证书信息
		}
	}
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

func DomainAcme(c *gin.Context) {
	param := &model.Domain{}
	if err := c.BindJSON(param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}

	if param.Domain == "" || param.SSLEmail == "" {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "必须设置域名与邮箱"})
		return
	}

	// 本地验证域名所有权
	preToken := "pre_" + base64.RawURLEncoding.EncodeToString([]byte(param.Domain))
	lib.Cache.SetAny("acme:"+preToken, param.Domain, 5*time.Minute) // 缓存5分钟
	resp, err := http.Get("http://" + param.Domain + "/.well-known/acme-challenge/" + preToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "无法访问验证URL: " + err.Error()})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "读取内容失败: " + err.Error()})
		return
	}
	if resp.StatusCode != http.StatusOK || strings.TrimSpace(string(body)) != param.Domain {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "本地验证域名所有权失败"})
		return
	}

	auth := &acme.Auth{
		Domain: []string{param.Domain},
		Email:  "your@email.com",
	}
	tasks, err := auth.AuthorizeOrder()
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	for i, task := range tasks {
		switch task.Type {
		case acme.ChallengeHTTP01:
			lib.Cache.SetAny("acme:"+task.Token, task.KeyAuth, 5*time.Minute) // 缓存5分钟
			tasks[i].Status = true
		}
	}

	cert, err := auth.CreateOrderCert(tasks)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	param.PublicKey = string(cert.Cert)
	param.PrivateKey = string(cert.Key)
	if err := param.DB().Save(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "保存证书失败(" + err.Error() + ")"})
		return
	}
	c.JSON(http.StatusOK, bind.Message{Message: "证书已生成并保存"})
}
