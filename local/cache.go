package local

import (
	"time"

	"github.com/miekg/dns"
	"github.com/zkfy/go-cache"
)

var defCache = newCache()

//Cache 缓存信息
type Cache struct {
	cache *cache.Cache
}

//newCache 创建缓存对象
func newCache() *Cache {
	return &Cache{
		cache: cache.New(time.Minute, time.Minute),
	}
}

//Lookup 查询域名解析
func (c *Cache) Lookup(name string) (*dns.Msg, bool) {
	return nil, false
}

//Set 保存到缓存
func (c *Cache) Set(name string, msg *dns.Msg) {
	c.cache.Set(name, msg, time.Second*60)
}
