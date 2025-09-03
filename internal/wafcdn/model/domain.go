package model

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/penndev/galite/pkg/orm"
)

type CertInfo struct {
	Issuer    string   `json:"issuer"`
	Subject   string   `json:"subject"`
	NotBefore string   `json:"not_before"`
	NotAfter  string   `json:"not_after"`
	DNSNames  []string `json:"dns_names"`
}

type Domain struct {
	orm.ModelBase
	Remark string `json:"remark"` // 备注

	Domain   string `json:"domain"`   // 站点域名
	Wildcard bool   `json:"wildcard"` // 通配符域名

	SiteID *uint `json:"SiteId"` // 必须用指针因为外键关联问题 foreign key constraint
	Site   Site  `json:"Site"`

	SSL        bool     `json:"ssl"`      // 是否启用ssl
	SSLEmail   string   `json:"sslEmail"` // 申请证书的邮箱
	SSLForce   bool     `json:"sslForce"` // 强制https
	PublicKey  string   `json:"publicKey"`
	PrivateKey string   `json:"privateKey"`
	CertInfo   CertInfo `json:"certInfo" gorm:"-"` // 证书信息，不存数据库
}

func (m *Domain) ParseCertInfo() error {
	// 不开启证书则不验证
	if !m.SSL {
		return nil
	}
	block, _ := pem.Decode([]byte(m.PublicKey))
	if block == nil {
		return errors.New("error pem byte")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}
	m.CertInfo = CertInfo{
		Issuer:    cert.Issuer.String(),
		Subject:   cert.Subject.String(),
		NotBefore: cert.NotBefore.Format("2006-01-02 15:04:05"),
		NotAfter:  cert.NotAfter.Format("2006-01-02 15:04:05"),
		DNSNames:  cert.DNSNames,
	}

	return nil
}
