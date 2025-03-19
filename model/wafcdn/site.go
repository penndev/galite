package wafcdn

import "github.com/penndev/galite/model/orm"

type SiteSecurity struct {
	Limit struct {
		Status  bool `json:"status"`
		Rate    int  `json:"rate"`
		Queries int  `json:"queries"`
		Seconds int  `json:"seconds"`
	} `json:"limit"`
	Sign struct {
		Status     bool   `json:"status"`
		Method     string `json:"method"`
		Key        string `json:"key"`
		ExpireArgs string `json:"expire_args"`
		SignArgs   string `json:"sign_args"`
	} `json:"sign"`
}

type SiteStatic struct {
	Root string `json:"root"`
}

type SiteProxy struct {
	Server            string            `json:"server"`
	Host              string            `json:"host"`
	KeepaliveTimeout  int               `json:"keepalive_timeout"`
	KeepaliveRequests int               `json:"keepalive_requests"`
	Header            map[string]string `json:"header"`
	Cache             []struct {
		Ruth   string   `json:"ruth"`
		Time   int      `json:"time"`
		Args   bool     `json:"args"`
		Method []string `json:"method"`
		Status []int    `json:"status"`
	} `json:"cache"`
}

type Site struct {
	orm.Model
	Security SiteSecurity      `gorm:"serializer:json" json:"security"`
	Static   SiteStatic        `gorm:"serializer:json" json:"static"`
	Proxy    SiteProxy         `gorm:"serializer:json" json:"proxy"`
	Header   map[string]string `gorm:"serializer:json" json:"header"`
}

type Domain struct {
	orm.Model
	Name string `json:"name"`
}
