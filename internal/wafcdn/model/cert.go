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

type Cert struct {
	orm.ModelBase
	Remark     string   `json:"remark"`
	Email      string   `json:"email"`
	Wildcard   bool     `json:"wildcard"`
	Domain     string   `json:"domain"`
	PublicKey  string   `json:"publicKey"`
	PrivateKey string   `json:"privateKey"`
	CertInfo   CertInfo `json:"certInfo" gorm:"-"`
}

func (m *Cert) ParseCertInfo() error {
	if m.PublicKey == "" && m.PrivateKey == "" {
		return nil
	}
	block, _ := pem.Decode([]byte(m.PublicKey))
	if block == nil {
		return errors.New("证书内容无法解析")
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
