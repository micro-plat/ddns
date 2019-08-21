package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/miekg/dns"
)

type RResp struct {
	msg        *dns.Msg
	nameserver string
	rtt        time.Duration
}

type IResolver interface {
	Lookup(net string, req *dns.Msg) (message *dns.Msg, err error)
}

type Resolver struct {
	nameServers []string
}

func NewResolver(server ...string) *Resolver {
	return &Resolver{nameServers: server}
}

//Lookup 循环所有名称服务器，以最快速度拿取解析信息，所有名称服务器都未能成功拿到解析信息则返回失败
func (r *Resolver) Lookup(net string, req *dns.Msg) (message *dns.Msg, err error) {
	c := &dns.Client{
		Net:          net,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}
	if net == "udp" {
		req = req.SetEdns0(65535, true)
	}

	//查询名称服务器，并处理结果
	qname := req.Question[0].Name
	res := make(chan *RResp, 1)
	var wg sync.WaitGroup
	lookup := func(nameserver string) {
		defer wg.Done()
		r, rtt, err := c.Exchange(req, nameserver)
		if err != nil {
			return
		}
		if r != nil && r.Rcode != dns.RcodeSuccess {
			if r.Rcode == dns.RcodeServerFailure {
				return
			}
		}
		select {
		case res <- &RResp{r, nameserver, rtt}:
		default:
		}
	}

	//循环所有名称服务器，每个服务器等待200毫秒，未拿到解析结果则发起下一个名称解析
	ticker := time.NewTicker(time.Millisecond * 200)
	defer ticker.Stop()
	for _, nameserver := range r.nameServers {
		wg.Add(1)
		go lookup(nameserver)
		select {
		case re := <-res:
			return re.msg, nil
		case <-ticker.C:
			continue
		}
	}
	wg.Wait()

	//处理返回结果
	select {
	case re := <-res:
		return re.msg, nil
	default:
		return nil, fmt.Errorf("无法解析的域名:%s", qname)
	}
}
