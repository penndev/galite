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

func (m *Cert) Normalize() error {
	m.Domain = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(m.Domain)), "*.")
	if m.Domain == "" || strings.Contains(m.Domain, "*") || strings.ContainsAny(m.Domain, " /") {
		return errors.New("请填写域名，例如 example.com")
	}
	if m.PublicKey == "" {
		return nil
	}
	if len(m.CertInfo.DNSNames) == 0 {
		return errors.New("证书没有域名")
	}
	apex := m.Domain
	wild := "*." + apex
	for _, name := range m.CertInfo.DNSNames {
		name = strings.ToLower(name)
		if name != apex && name != wild {
			return errors.New("证书只能包含域名，或再加它的泛解析")
		}
	}
	if !m.Wildcard {
		for _, name := range m.CertInfo.DNSNames {
			if strings.EqualFold(name, wild) {
				return errors.New("证书含有泛解析，请打开泛解析")
			}
		}
		return nil
	}
	for _, name := range m.CertInfo.DNSNames {
		if strings.EqualFold(name, wild) {
			return nil
		}
	}
	return errors.New("证书未包含 " + wild)
}

func (m *Cert) covers(host string) bool {
	if err := m.ParseCertInfo(); err != nil {
		return false
	}
	host = strings.ToLower(strings.TrimSpace(host))
	if host == "" {
		return false
	}
	for _, pattern := range m.CertInfo.DNSNames {
		pattern = strings.ToLower(strings.TrimSpace(pattern))
		if pattern == "" {
			continue
		}
		if pattern == host {
			return true
		}
		if !strings.HasPrefix(pattern, "*.") || strings.HasPrefix(host, "*.") {
			continue
		}
		suffix := pattern[1:]
		if !strings.HasSuffix(host, suffix) {
			continue
		}
		label := strings.TrimSuffix(host, suffix)
		if label != "" && !strings.Contains(label, ".") {
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
