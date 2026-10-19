package model

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"

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

func dnsNameCovers(pattern, host string) bool {
	pattern = strings.ToLower(strings.TrimSpace(pattern))
	host = strings.ToLower(strings.TrimSpace(host))
	if pattern == "" || host == "" {
		return false
	}
	if pattern == host {
		return true
	}
	if !strings.HasPrefix(pattern, "*.") || strings.HasPrefix(host, "*.") {
		return false
	}
	suffix := pattern[1:]
	if !strings.HasSuffix(host, suffix) {
		return false
	}
	label := strings.TrimSuffix(host, suffix)
	return label != "" && !strings.Contains(label, ".")
}

func (m *Cert) covers(host string) bool {
	if err := m.ParseCertInfo(); err != nil {
		return false
	}
	for _, name := range m.CertInfo.DNSNames {
		if dnsNameCovers(name, host) {
			return true
		}
	}
	return false
}

func FindCertForHost(host string) (Cert, bool) {
	var list []Cert
	err := (&Cert{}).DB().Where("public_key <> '' AND private_key <> ''").Find(&list).Error
	if err != nil {
		return Cert{}, false
	}
	for i := range list {
		if list[i].covers(host) {
			return list[i], true
		}
	}
	return Cert{}, false
}
