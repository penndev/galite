package admin

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/penndev/galite/internal/admin/bind"
	"github.com/penndev/galite/internal/lib"
	"github.com/penndev/galite/internal/wafcdn/model"
	"github.com/penndev/gopkg/acme"
	"golang.org/x/net/publicsuffix"
)

type acmeDNSSession struct {
	Auth   *acme.Auth
	Tasks  []acme.ChallengeTask
	CertID uint
	Remark string
	Email  string
	At     time.Time
}

var acmeDNSSessions sync.Map

func CertList(c *gin.Context) {
	param := &bindCertParam{}
	if err := c.BindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}
	var total int64
	var list []model.Cert
	m := param.Param()
	m.List(&total, &list)
	for i := range list {
		list[i].ParseCertInfo()
	}
	c.JSON(http.StatusOK, bind.DataList{Total: total, Data: list})
}

func CertSave(c *gin.Context) {
	param := &model.Cert{}
	if err := c.BindJSON(param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}
	if err := param.ParseCertInfo(); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "证书错误" + err.Error()})
		return
	}
	if err := normalizeCertWildcard(param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if err := param.DB().Save(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "存储失败(" + err.Error() + ")"})
		return
	}
	c.JSON(http.StatusOK, bind.Message{Message: "完成"})
}

func normalizeCertWildcard(cert *model.Cert) error {
	for _, name := range cert.CertInfo.DNSNames {
		if strings.HasPrefix(name, "*.") && cert.Domain == "" {
			cert.Wildcard = true
			cert.Domain = strings.TrimPrefix(name, "*.")
			break
		}
	}
	if !cert.Wildcard {
		cert.Domain = ""
		return nil
	}
	apex, err := publicsuffix.EffectiveTLDPlusOne(strings.TrimSpace(cert.Domain))
	if err != nil || apex == "" {
		return errors.New("请填写泛解析主域名，例如 example.com")
	}
	cert.Domain = apex
	if cert.PublicKey == "" {
		return nil
	}
	want := "*." + apex
	for _, name := range cert.CertInfo.DNSNames {
		if strings.EqualFold(name, want) {
			return nil
		}
	}
	return errors.New("证书未包含 " + want)
}

func CertDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误"})
		return
	}
	var count int64
	(&model.Domain{}).DB().Model(&model.Domain{}).Where("cert_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "仍有域名绑定该证书"})
		return
	}
	param := &model.Cert{}
	param.ID = uint(id)
	if err := param.DB().Delete(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "删除失败(" + err.Error() + ")"})
		return
	}
	c.JSON(http.StatusOK, bind.Message{Message: "完成"})
}

func CertAcme(c *gin.Context) {
	var body struct {
		ID       uint   `json:"id"`
		Remark   string `json:"remark"`
		Email    string `json:"email"`
		Wildcard bool   `json:"wildcard"`
		Domain   string `json:"domain"`
		Session  string `json:"session"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}
	if body.Session != "" {
		finishCertAcmeDNS(c, body.Session)
		return
	}
	domainName := strings.TrimSpace(body.Domain)
	if strings.HasPrefix(domainName, "*.") {
		startCertAcmeDNS(c, body.ID, body.Remark, body.Email, domainName)
		return
	}
	if domainName == "" || body.Email == "" {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "必须设置域名与邮箱"})
		return
	}
	cert := &model.Cert{}
	if body.ID > 0 {
		if err := cert.DB().First(cert, body.ID).Error; err != nil {
			c.JSON(http.StatusBadRequest, bind.Message{Message: "证书不存在"})
			return
		}
	}
	cert.Email = body.Email
	if cert.Remark == "" {
		cert.Remark = body.Remark
	}
	if cert.Remark == "" {
		cert.Remark = domainName
	}

	preToken := "pre_" + base64.RawURLEncoding.EncodeToString([]byte(domainName))
	lib.Cache.SetAny("acme:"+preToken, domainName, 5*time.Minute)
	resp, err := http.Get("http://" + domainName + "/.well-known/acme-challenge/" + preToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "无法访问验证URL: " + err.Error()})
		return
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "读取内容失败: " + err.Error()})
		return
	}
	if resp.StatusCode != http.StatusOK || strings.TrimSpace(string(respBody)) != domainName {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "本地验证域名所有权失败"})
		return
	}

	auth := &acme.Auth{
		Domain: []string{domainName},
		Email:  cert.Email,
	}
	tasks, err := auth.AuthorizeOrder()
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	for i, task := range tasks {
		switch task.Type {
		case acme.ChallengeHTTP01:
			lib.Cache.SetAny("acme:"+task.Token, task.KeyAuth, 5*time.Minute)
			tasks[i].Status = true
		}
	}
	issued, err := auth.CreateOrderCert(tasks)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	cert.PublicKey = string(issued.Cert)
	cert.PrivateKey = string(issued.Key)
	if err := cert.ParseCertInfo(); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "证书错误" + err.Error()})
		return
	}
	if err := normalizeCertWildcard(cert); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if err := cert.DB().Save(cert).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "保存证书失败(" + err.Error() + ")"})
		return
	}
	c.JSON(http.StatusOK, bind.Message{Message: "证书已生成并保存"})
}

func startCertAcmeDNS(c *gin.Context, certID uint, remark, email, domainName string) {
	if email == "" {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "必须设置域名与邮箱"})
		return
	}
	auth := &acme.Auth{
		Domain: []string{domainName},
		Email:  email,
	}
	tasks, err := auth.AuthorizeOrder()
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	dnsTasks := make([]acme.ChallengeTask, 0)
	records := make([]gin.H, 0)
	for _, task := range tasks {
		if task.Type != acme.ChallengeDNS01 {
			continue
		}
		dnsTasks = append(dnsTasks, task)
		records = append(records, gin.H{
			"name":  "_acme-challenge." + task.Domain,
			"value": task.KeyAuth,
		})
	}
	if len(dnsTasks) == 0 {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "未返回 DNS 验证"})
		return
	}
	id := uuid.NewString()
	acmeDNSSessions.Store(id, &acmeDNSSession{
		Auth:   auth,
		Tasks:  dnsTasks,
		CertID: certID,
		Remark: remark,
		Email:  email,
		At:     time.Now(),
	})
	c.JSON(http.StatusOK, gin.H{
		"dns":     true,
		"session": id,
		"records": records,
	})
}

func finishCertAcmeDNS(c *gin.Context, sessionID string) {
	value, ok := acmeDNSSessions.Load(sessionID)
	if !ok {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "验证已过期，请重新申请"})
		return
	}
	sess := value.(*acmeDNSSession)
	if time.Since(sess.At) > 30*time.Minute {
		acmeDNSSessions.Delete(sessionID)
		c.JSON(http.StatusBadRequest, bind.Message{Message: "验证已过期，请重新申请"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	sess.Auth.Ctx = ctx
	for i := range sess.Tasks {
		sess.Tasks[i].Status = true
	}
	issued, err := sess.Auth.CreateOrderCert(sess.Tasks)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	cert := &model.Cert{}
	if sess.CertID > 0 {
		if err := cert.DB().First(cert, sess.CertID).Error; err != nil {
			c.JSON(http.StatusBadRequest, bind.Message{Message: "证书不存在"})
			return
		}
	}
	cert.Email = sess.Email
	if cert.Remark == "" {
		cert.Remark = sess.Remark
	}
	if cert.Remark == "" && len(issued.Domain) > 0 {
		cert.Remark = issued.Domain[0]
	}
	cert.PublicKey = string(issued.Cert)
	cert.PrivateKey = string(issued.Key)
	if err := cert.ParseCertInfo(); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "证书错误" + err.Error()})
		return
	}
	if err := normalizeCertWildcard(cert); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if err := cert.DB().Save(cert).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "保存证书失败(" + err.Error() + ")"})
		return
	}
	acmeDNSSessions.Delete(sessionID)
	c.JSON(http.StatusOK, bind.Message{Message: "证书已生成并保存"})
}
