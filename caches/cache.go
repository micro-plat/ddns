package caches

import (
	"net"

	"github.com/micro-plat/lib4go/logger"
)

//ICache 缓存DNS
type ICache interface {
	Lookup(string) []net.IP
	Close() error
}

//Cache 缓存信息
type Cache struct {
}

//NewCache 创建缓存对象
func NewCache(log logger.ILogger) (*Cache, error) {
	if err := Start(log); err != nil {
		return nil, err
	}
	return &Cache{}, nil
}

//Lookup 查询域名解析
func (c *Cache) Lookup(name string) []net.IP {
	ips := Lookup(name)
	return ips
}

//Close 关闭缓存服务
func (c *Cache) Close() error {
	return Close()
}
