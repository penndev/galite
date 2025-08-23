package model

import (
	"crypto/x509"
	"encoding/pem"

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
	Domain string `json:"domain" binding:"hostname_rfc1123"` // 站点域名
	Remark string `json:"remark"`                            // 备注

	SiteID *uint `json:"SiteId"` // 必须用指针因为外键关联问题 foreign key constraint
	Site   Site  `json:"Site"`

	SSL        bool     `json:"ssl"`                      // 是否启用ssl
	SSLEmail   string   `json:"sslEmail" binding:"email"` // 申请证书的邮箱，必须邮箱格式
	SSLForce   bool     `json:"sslForce"`                 // 强制https
	PublicKey  string   `json:"publicKey"`
	PrivateKey string   `json:"privateKey"`
	CertInfo   CertInfo `json:"certInfo" gorm:"-"` // 证书信息，不存数据库
}

func (m *Domain) ParseCertInfo() error {

	var block *pem.Block
	block, _ = pem.Decode([]byte(m.PublicKey))
	// if block == nil {
	// 	break
	// }
	// if block.Type != "CERTIFICATE" {
	// 	continue
	// }
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
