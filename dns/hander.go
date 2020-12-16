package dns

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/types"
	"github.com/miekg/dns"
)

//DNSHandler 处理DNS域名解析
type DNSHandler struct {
	local  ILocal
	remote IResolver
	log    logger.ILogger
}

//NewHandler 构建DNS处理程序
func NewHandler(log logger.ILogger) (*DNSHandler, error) {
	local, err := NewLocal(log)
	if err != nil {
		return nil, err
	}
	resolver, err := NewResolver(log)
	if err != nil {
		return nil, fmt.Errorf("创建域名解析器失败 %w", err)
	}
	return &DNSHandler{
		remote: resolver,
		local:  local,
		log:    log,
	}, nil
}
func sortByLocalFirst(addr net.IP, ips []net.IP) []net.IP {
	for idx, ip := range ips {
		if ip.Equal(addr) {
			sips := make([]net.IP, 0, len(ips))
			sips = append(sips, addr)
			sips = append(sips, ips[0:idx]...)
			sips = append(sips, ips[idx+1:]...)
			return sips
		}
	}
	return ips
}

//lookupFromLocal 从缓存中查询DNS服务
func (h *DNSHandler) lookupFromLocal(ip net.IP, question *Question, q dns.Question) (*dns.Msg, error) {
	ips := h.local.Lookup(question.qname)
	if len(ips) == 0 {
		return nil, fmt.Errorf("缓存中未查询到:%s", question.qname)
	}
	ips = sortByLocalFirst(ip, ips)
	m := new(dns.Msg)
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
	return m, nil
}


func (h *DNSHandler) do(net string, fromIP net.IP, w dns.ResponseWriter, req *dns.Msg) (string, *dns.Msg, error) {
	//从缓存中查询
	cache := true
	msg, err := h.lookupFromLocal(fromIP, NewQuestion(req.Question[0]), req.Question[0])
	if err != nil {
		msg, cache, err = h.remote.Lookup(net, req) //从名称服务中查询
	}
	if err != nil {
		return "", nil, err
	}
	//处理响应
	msg.SetReply(req)
	w.WriteMsg(msg)
	return types.DecodeString(cache, true, "C", "R"), msg, nil
}


//Close 关闭处理程序
func (h *DNSHandler) Close() error {
	return h.local.Close()
}
