package dt

type SchemaType string

const (
	SchemaTypeSocks5 = "socks5"
	SchemaTypeHttp   = "http"
	SchemaTypeHttps  = "https"
)

type ServerItem struct {
	Schema SchemaType `json:"schema"`
	Ip     string     `json:"ip"`
	Port   int        `json:"port"`
}

type ServerList []ServerItem
