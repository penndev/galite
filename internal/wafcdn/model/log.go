package model

import "github.com/penndev/galite/pkg/orm"

type Log struct {
	orm.ModelBase
	SiteID        uint   `json:"site_id" form:"site_id" binding:"gte=0"`
	Host          string `json:"host"`
	RemoteAddr    string `json:"remote_addr"`
	HTTPReferer   string `json:"http_referer"`
	HTTPUserAgent string `json:"http_user_agent"`
	Request       string `json:"request"`
	RequestMethod string `json:"request_method"`
	RequestTime   string `json:"request_time"`
	Status        int    `json:"status"`
	BytesReceived int64  `json:"bytes_received"`
	BytesSent     int64  `json:"bytes_sent"`
}
