package weibo

import "github.com/juxuny/env"

type cookieManager struct {
}

func NewCookieManager() *cookieManager {
	return &cookieManager{}
}

func (t *cookieManager) GetCookies(domain string) string {
	return env.GetString("COOKIE")
}
