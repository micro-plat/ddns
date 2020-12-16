package local

import (
	"net"

	"github.com/micro-plat/hydra"
	"github.com/miekg/dns"
)

type Local struct {
	r *Registry
	h *Hosts
	c *Cache
}

//IQuery 查询服务
type IQuery interface {
	Lookup(*dns.Msg) ([]net.IP, bool)
}

//New 构建本地服务
func New() (*Local, error) {
	l := &Local{
		r: newRegistry(),
		h: NewHosts(hydra.G.Log()),
		c: newCache(),
	}
	if err := l.r.Start(); err != nil {
		return nil, err
	}
	return l, nil
}

//Lookup 根据域名查询
func (l *Local) Lookup(req *dns.Msg) (*dns.Msg, bool) {

	//从本地缓存获取
	domain := TrimDomain(req.Question[0].Name)
	if msg, ok := l.c.Lookup(domain, req); ok {
		return msg, ok
	}
	ips, ok := l.r.Lookup(domain)
	if !ok {
		ips, ok = l.h.Lookup(domain)
	}
	if !ok || len(ips) == 0 {
		return nil, false
	}
	return pack(ips, req), true
}

// func (l *Local) CacheItems() interface{} {
// 	return l.c.Items()
// }

// func (l *Local) RegistryItems() interface{} {
// 	return l.r.domains.Items()
// }

// func (l *Local) HostItems() interface{} {
// 	return l.h.domain
// }

//pack 对本地ip的包进行打包处理
func pack(ips []net.IP, req *dns.Msg) *dns.Msg {
	q := req.Question[0]
	question := NewQuestion(q)
	m := dns.Msg{}
	m.Id = req.Id
	switch question.QueryType() {
	case _IP4Query:
		header := dns.RR_Header{
			Name:   q.Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    600,
		}
		for _, ip := range ips {
			m.Answer = append(m.Answer, &dns.A{header, ip})
		}
	case _IP6Query:
		header := dns.RR_Header{
			Name:   q.Name,
			Rrtype: dns.TypeAAAA,
			Class:  dns.ClassINET,
			Ttl:    600,
		}
		for _, ip := range ips {
			m.Answer = append(m.Answer, &dns.AAAA{header, ip})
		}
	}
	return &m
}

//Save2Cache 保存到缓存
func (l *Local) Save2Cache(msg *dns.Msg) {
	name := TrimDomain(msg.Question[0].Name)
	l.c.Set(name, msg)
}

//Close 关闭服务
func (l *Local) Close() {
	if l.r != nil {
		l.r.Close()
	}
}
