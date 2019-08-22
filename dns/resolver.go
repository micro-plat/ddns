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
}

func NewResolver() (*Resolver, error) {
	name := NewNames(logger.New("ddns_names"))
	if err := name.Start(); err != nil {
		return nil, err
	}
	return &Resolver{
		cache: cache.New(5*time.Minute, 10*time.Minute),
		names: name,
	}, nil
}

//Lookup 循环所有名称服务器，以最快速度拿取解析信息，所有名称服务器都未能成功拿到解析信息则返回失败
func (r *Resolver) Lookup(net string, req *dns.Msg) (message *dns.Msg, cache bool, err error) {

	//查询本地缓存
	msg, ok := r.lookupFromCache(req)
	if ok {
		return msg, true, nil
	}

	//查询远程服务
	msg, err = r.lookupFromRemote(net, req)
	if err != nil {
		return nil, false, err
	}

	//保存缓存
	r.save2Cache(req.Question[0].Name, msg)
	return msg, false, nil

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
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}
	if net == "udp" {
		req = req.SetEdns0(65535, true)
	}
	qname := req.Question[0].Name
	res := make(chan *dns.Msg, 1)
	var wg sync.WaitGroup
	lookup := func(nameserver string) {
		defer wg.Done()
		r, _, err := c.Exchange(req, nameserver)
		if err != nil {
			return
		}
		if r != nil && r.Rcode != dns.RcodeSuccess {
			if r.Rcode == dns.RcodeServerFailure {
				return
			}
		}
		select {
		case res <- r:
		default:
		}
	}

	//循环所有名称服务器，每个服务器等待200毫秒，未拿到解析结果则发起下一个名称解析
	ticker := time.NewTicker(time.Millisecond * 200)
	defer ticker.Stop()
	names := r.names.Lookup()
	for _, host := range names {
		wg.Add(1)
		go lookup(host)
		select {
		case re := <-res:
			return re, nil
		case <-ticker.C:
			continue
		}
	}
	wg.Wait()

	//处理返回结果
	select {
	case re := <-res:
		return re, nil
	default:
		return nil, fmt.Errorf("无法解析的域名:%s", qname)
	}
}
