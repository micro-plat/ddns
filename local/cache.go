package local

import (
	"time"

	"github.com/miekg/dns"
	"github.com/zkfy/go-cache"
)

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
func (c *Cache) Lookup(req *dns.Msg) (*dns.Msg, bool) {
	if v, ok := c.cache.Get(req.Question[0].Name); ok {
		msg := v.(*dns.Msg)
		m := dns.Msg{}
		m = *msg
		m.Id = req.Id
		return &m, true
	}
	return nil, false
}

//Set 保存到缓存
func (c *Cache) Set(msg *dns.Msg) {
	c.cache.Set(msg.Question[0].Name, msg, time.Second*60)
}
