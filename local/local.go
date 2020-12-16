package local

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
)

//localQueries 本地查询服务
var qmaps = make(map[string]IQuery)
var queries = make([]IQuery, 0, 1)

//IQuery 查询服务
type IQuery interface {
	Lookup(*dns.Msg) ([]net.IP, bool)
}

//Register 注册查询服务
func Register(name string, query IQuery) {
	if _, ok := qmaps[name]; ok {
		panic(fmt.Sprintf("%s:重复注册查询服务", name))
	}
	qmaps[name] = query
	queries = append(queries, query)
}

//Lookup 根据域名查询
func Lookup(req *dns.Msg) (*dns.Msg, bool) {
	var ips []net.IP
	var ok bool
	for _, q := range queries {
		if ips, ok = q.Lookup(req); ok {
			break
		}
	}
	if !ok {
		return nil, false
	}

	//处理响应包
	var msg *dns.Msg
	q := req.Question[0]
	question := NewQuestion(q)
	m := dns.Msg{}
	m = *msg
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

	return &m, true
}

//Save2Cache 保存到缓存
func Save2Cache(msg *dns.Msg) {
	defCache.Set(msg.Question[0].Name, msg)
}
