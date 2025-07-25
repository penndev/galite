package model

import "github.com/penndev/galite/pkg/orm"

type SiteSecurity struct {
	Limit struct {
		Status  bool `json:"status"`
		Rate    int  `json:"rate"`    // 限制下载速率 kb/s
		Queries int  `json:"queries"` // 请求次数
		Seconds int  `json:"seconds"` // 多少秒内允许的 Queries
	} `json:"limit"`
	Sign struct {
		Status     bool   `json:"status"`
		Method     string `json:"method"`
		Key        string `json:"key"`
		ExpireArgs string `json:"expire_args"`
		SignArgs   string `json:"sign_args"`
	} `json:"sign"`
	Cors struct {
		Status      bool   `json:"status"`
		Origin      string `json:"origin"`
		Method      string `json:"method"`
		Header      string `json:"header"`
		Credentials string `json:"credentials"`
		Age         int    `json:"age"`
	} `json:"cors"`
}

type SiteStatic struct {
	Root string `json:"root"`
}

type SiteHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SiteProxy struct {
	Server            string       `json:"server"`
	Host              string       `json:"host"`
	KeepaliveTimeout  int          `json:"keepalive_timeout"`
	KeepaliveRequests int          `json:"keepalive_requests"`
	Header            []SiteHeader `json:"header"`
	Cache             []struct {
		Ruth   string   `json:"ruth"`   // 缓存路径 正则表达式
		Time   int      `json:"time"`   // 缓存时间/秒
		Args   bool     `json:"args"`   // 是否忽略参数
		Method []string `json:"method"` // 缓存的http方法 GET POST
		Status []int    `json:"status"` // 缓存的http状态码
	} `json:"cache"`
}

type Site struct {
	orm.Model
	Type     string       `json:"type"`
	Remark   string       `json:"remark"`
	Security SiteSecurity `gorm:"serializer:json" json:"security"`
	Static   SiteStatic   `gorm:"serializer:json" json:"static"`
	Proxy    SiteProxy    `gorm:"serializer:json" json:"proxy"`
	Header   []SiteHeader `gorm:"serializer:json" json:"header"`

	// 逻辑关联
	Domains []Domain // 域名one to many
}
