package model

import (
	"strings"

	"gorm.io/gorm"
)

type legacyDomainCert struct {
	ID         uint
	Remark     string
	Domain     string
	SSLEmail   string `gorm:"column:ssl_email"`
	PublicKey  string `gorm:"column:public_key"`
	PrivateKey string `gorm:"column:private_key"`
}

// 注册表结构。
func Migration(db *gorm.DB) {
	// 后台系统用户权限管理模块
	db.AutoMigrate(&Log{})
	db.AutoMigrate(&Cache{})
	db.AutoMigrate(&Site{})
	db.AutoMigrate(&Cert{})
	db.AutoMigrate(&Domain{})
	db.AutoMigrate(&CacheDelete{})
	migrateLegacyDomainCerts(db)
	migrateWildcardDomains(db)
}

func migrateWildcardDomains(db *gorm.DB) {
	if !db.Migrator().HasTable("domains") || !db.Migrator().HasColumn("domains", "wildcard") {
		return
	}
	var rows []struct {
		ID     uint
		Domain string
	}
	db.Table("domains").Where("wildcard = ?", true).Find(&rows)
	for _, row := range rows {
		name := row.Domain
		if !strings.HasPrefix(name, "*.") {
			name = "*." + name
		}
		db.Table("domains").Where("id = ?", row.ID).Updates(map[string]any{
			"domain":   name,
			"wildcard": false,
		})
	}
}

func migrateLegacyDomainCerts(db *gorm.DB) {
	if !db.Migrator().HasTable("domains") || !db.Migrator().HasColumn("domains", "public_key") {
		return
	}
	var rows []legacyDomainCert
	db.Table("domains").Where("cert_id IS NULL AND public_key <> '' AND private_key <> ''").Find(&rows)
	for _, row := range rows {
		cert := Cert{}
		err := db.Where("public_key = ? AND private_key = ?", row.PublicKey, row.PrivateKey).First(&cert).Error
		if err != nil {
			remark := row.Remark
			if remark == "" {
				remark = row.Domain
			}
			cert = Cert{
				Remark:     remark,
				Email:      row.SSLEmail,
				PublicKey:  row.PublicKey,
				PrivateKey: row.PrivateKey,
			}
			if err = db.Create(&cert).Error; err != nil {
				continue
			}
		}
		db.Model(&Domain{}).Where("id = ?", row.ID).Update("cert_id", cert.ID)
	}
}
