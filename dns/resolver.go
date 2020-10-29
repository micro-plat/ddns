package dns

import (
	"fmt"
	"sync"
	"time"

	"github.com/micro-plat/lib4go/logger"
	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
)

type IResolver interface {
	Lookup(net string, req *dns.Msg) (message *dns.Msg, cache bool, err error)
}

type Resolver struct {
	cache *cache.Cache
	names *Names
	log   logger.ILogger
}

func NewResolver(log logger.ILogger) (*Resolver, error) {
	name := NewNames(log)
	if err := name.Start(); err != nil {
		return nil, err
	}
	return &Resolver{
		cache: cache.New(5*time.Minute, 10*time.Minute),
		names: name,
		log:   log,
	}, nil
}

//Lookup 循环所有名称服务器，以最快速度拿取解析信息，所有名称服务器都未能成功,再次从缓存中获取
func (r *Resolver) Lookup(net string, req *dns.Msg) (message *dns.Msg, cache bool, err error) {

	//查询本地缓存
	cmsg, ok := r.lookupFromCache(req)
	if ok {
		return cmsg, true, nil
	}

	//查询远程服务
	rmsg, err := r.lookupFromRemote(net, req)
	if err != nil {
		return nil, false, err
	}

	//保存缓存
	if len(rmsg.Answer) > 0 {
		return rmsg, false, nil
	}
	//再次从缓存中拉取，解决并发请求时部分请求未能从名称服务器中获取到结果的问题
	cmsg, ok = r.lookupFromCache(req)
	if ok {
		return cmsg, true, nil
	}
	return rmsg, true, nil

}

//lookupFromCache 从缓存中获和取解析信息
func (r *Resolver) lookupFromCache(req *dns.Msg) (message *dns.Msg, f bool) {
	qname := req.Question[0].Name
	msg, ok := r.cache.Get(qname)
	if !ok {
		return nil, false
	}
	message = msg.(*dns.Msg)
	message.Id = req.Id
	return message, true
}

//save2Cache 保存到缓存
func (r *Resolver) save2Cache(name string, req *dns.Msg) {
	r.cache.Set(name, req, time.Second*60*5)
}

//lookupFromRemote 从远程服务器查询解析信息
func (r *Resolver) lookupFromRemote(net string, req *dns.Msg) (message *dns.Msg, err error) {
	//查询名称服务器，并处理结果
	c := &dns.Client{
		Net:          net,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}
	if net == "udp" {
		req = req.SetEdns0(65535, true)
	}
	logger := logger.New("ctx")
	qname := req.Question[0].Name
	response := make(chan *dns.Msg, 1)
	var wg sync.WaitGroup
	lookup := func(nameserver string) {
		defer wg.Done()
		logger.Debug("exchange.start:", req.Question[0].Name, "by", nameserver)
		defer logger.Debug("exchange.end:", req.Question[0].Name, "by", nameserver)
		res, _, err1 := c.Exchange(req, nameserver)
		if err1 != nil {
			err = fmt.Errorf("%v,%v", err, err1)
			return
		}
		if res != nil {
			if res.Rcode == dns.RcodeServerFailure {
				return
			}
		}
		select {
		case response <- res:
		default:
		}
	}

	//循环所有名称服务器，每个服务器等待200毫秒，未拿到解析结果则发起下一个名称解析
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	names := r.names.Lookup()
	logger.Debug("lookup:", req.Question[0].Name, "from", names)
	for _, host := range names {
		wg.Add(1)
		go lookup(host)
		select {
		case re := <-response:
			return re, nil
		case <-ticker.C:
			continue
		}
	}
	wg.Wait()

	//处理返回结果
	select {
	case re := <-response:
		return re, nil
	default:
		logger.Debugf("无法解析的域名:%s[%v](%v)", qname, names, err)
		return nil, fmt.Errorf("无法解析的域名:%s[%v](%v)", qname, names, err)
	}
}
