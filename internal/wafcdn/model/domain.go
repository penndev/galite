package model

import "github.com/penndev/galite/pkg/orm"

type Domain struct {
	orm.ModelBase
	Name       string `json:"name"`
	SiteID     *uint  `json:"SiteId"` // 必须用指针因为外键关联问题 foreign key constraint
	Site       Site   `json:"Site"`
	SSL        bool   `json:"ssl"`
	SSLForce   bool   `json:"sslForce"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}
