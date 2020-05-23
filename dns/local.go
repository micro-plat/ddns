package dns

import (
	"net"

	"github.com/micro-plat/ddns/dns/query"
	"github.com/micro-plat/lib4go/logger"
)

//ILocal 缓存DNS
type ILocal interface {
	Lookup(string) []net.IP
	Close() error
}

//Local 缓存信息
type Local struct {
}

//NewLocal 创建缓存对象
func NewLocal(log logger.ILogger) (*Local, error) {
	if err := query.Start(log); err != nil {
		return nil, err
	}
	return &Local{}, nil
}

//Lookup 查询域名解析
func (c *Local) Lookup(name string) []net.IP {
	ips := query.Lookup(name)
	return ips
}

//Close 关闭缓存服务
func (c *Local) Close() error {
	return query.Close()
}
