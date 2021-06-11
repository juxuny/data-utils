package proxy

import (
	"github.com/juxuny/data-utils/proxy/dt"
	"github.com/juxuny/data-utils/proxy/fetcher"
)

type Fetcher interface {
	Len() int
	Page(page int, pageSize int) (ret dt.ServerList, err error)
	Init() (err error)
	Reset()
}

type Provider string

const (
	// http://www.66ip.cn/
	Provider66Ip = Provider("66ip")
)

var providerMap = map[Provider]Fetcher{
	Provider66Ip: fetcher.New66Ip(),
}

func Fetch(provider Provider) (ret dt.ServerList, err error) {
	return
}
