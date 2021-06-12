package dt

import "github.com/juxuny/data-utils/model"

type SchemaType string

const (
	SchemaTypeSocks5 = "socks5"
	SchemaTypeHttp   = "http"
	SchemaTypeHttps  = "https"
)

type ServerItem struct {
	Schema   SchemaType     `json:"schema"`
	Ip       string         `json:"ip"`
	Port     int            `json:"port"`
	Provider model.Provider `json:"provider"`
}

type ServerList []ServerItem

type SaveHandler interface {
	SaveServerList(list ServerList) error
}
