package wafcdn

import "github.com/penndev/galite/model/orm"

type Domain struct {
	orm.Model
	Name   string `json:"name"`
	SiteID *uint  `json:"SiteId"` // 必须用指针因为外键关联问题 foreign key constraint
	Site   Site   `json:"Site"`
}
