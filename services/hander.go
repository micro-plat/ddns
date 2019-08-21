package services

import (
	"fmt"
	"net"
	"time"

	"github.com/micro-plat/lib4go/logger"
	"github.com/miekg/dns"
)

//DNSHandler 处理DNS域名解析
type DNSHandler struct {
	cache    ICache
	resolver IResolver
	log      logger.ILogger
}

//NewHandler 构建DNS处理程序
func NewHandler(log logger.ILogger) *DNSHandler {
	return &DNSHandler{
		resolver: NewResolver("114.114.114.114:53", "8.8.8.8:53"),
		cache:    NewCache(),
		log:      log,
	}
}

//lookupFromCache 从缓存中查询DNS服务
func (h *DNSHandler) lookupFromCache(question *Question, q dns.Question) (*dns.Msg, error) {
	ips := h.cache.Lookup(question.qname)
	if len(ips) == 0 {
		return nil, fmt.Errorf("缓存中未查询到:%s", question.qname)
	}
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

func (h *DNSHandler) do(net string, w dns.ResponseWriter, req *dns.Msg) string {
	//从缓存中查询
	from := "cache"
	msg, err := h.lookupFromCache(NewQuestion(req.Question[0]), req.Question[0])
	if err != nil {
		from = "remote"
		msg, err = h.resolver.Lookup(net, req) //从名称服务中查询
	}
	if err != nil {
		panic(err)
	}

	//处理响应
	msg.SetReply(req)
	w.WriteMsg(msg)
	return from
}

//Do 处理请求
func (h *DNSHandler) Do(proto string) func(w dns.ResponseWriter, req *dns.Msg) {

	return func(w dns.ResponseWriter, req *dns.Msg) {
		start := time.Now()
		log := logger.New("ddns_" + proto)
		defer recovery(log, w, req)
		log.Info("dns.request", proto, req.Question[0].Name, "from", w.RemoteAddr().(*net.UDPAddr).IP)

		from := h.do(proto, w, req)
		log.Info("dns.response", proto, from, time.Since(start))
	}

}
