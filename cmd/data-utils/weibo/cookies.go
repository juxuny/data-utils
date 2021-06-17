package weibo

import (
	"fmt"
	"github.com/juxuny/data-utils/cache"
	"github.com/juxuny/data-utils/log"
)

var (
	cacheDir = "tmp/cache"
)

func SetCacheDir(dir string) {
	cacheDir = dir
}

type cookieManager struct {
	cache.Cache
}

func NewCookieManager() *cookieManager {
	return &cookieManager{
		Cache: cache.NewFileCache(cacheDir),
	}
}

func (t *cookieManager) GetCookies(domain string) string {
	k := fmt.Sprintf("cookie:%s", domain)
	v, err := t.Get(k)
	if err != nil {
		log.Warnf("get key('%s') failed: %v", k, err)
		return ""
	}
	return v
}
