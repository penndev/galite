package model

import "github.com/penndev/galite/pkg/orm"

type Domain struct {
	orm.ModelBase
	Remark string `json:"remark"` // 备注

	Domain string `json:"domain"` // 站点域名，泛解析直接写 *.example.com

	SiteID *uint `json:"SiteId"` // 必须用指针因为外键关联问题 foreign key constraint
	Site   Site  `json:"Site"`

	SSLForce bool `json:"sslForce"` // 自动跳转 https
}
