package model

import "time"

type Provider string

const (
	// http://www.66ip.cn/
	Provider66Ip = Provider("66ip")
)

type Proxy struct {
	Id           int64      `json:"id" gorm:"TYPE:BIGINT(20);PRIMARY_KEY;AUTO_INCREMENT"`
	Ip           string     `json:"ip"`
	Port         int64      `json:"port"`
	ProviderName Provider   `json:"provider_name"`
	Socks5       bool       `json:"socks5" gorm:"COLUMN:socks5"`
	Http         bool       `json:"http"`
	Https        bool       `json:"https"`
	CreatedAt    *time.Time `json:"created_at" gorm:"TYPE:TIMESTAMP;DEFAULT"`
	UpdatedAt    *time.Time `json:"updated_at" gorm:"TYPE:TIMESTAMP;DEFAULT"`
	Latency      float64    `json:"latency"`
	UsedCount    int64      `json:"used_count"`
}

func (Proxy) TableName() string {
	return "proxy"
}

type ProxyList []Proxy
