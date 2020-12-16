package remote

import (
	"fmt"
	"time"

	"github.com/micro-plat/ddns/names"
	"github.com/miekg/dns"
)

type Remote struct {
	names *names.Names
}

//New 构建远程解析器
func New() *Remote {
	return &Remote{
		names: names.New(),
	}
}

//Lookup 从远程服务器查询解析信息
func (r *Remote) Lookup(req *dns.Msg) (message *dns.Msg, err error) {
	//查询名称服务器，并处理结果
	c := &dns.Client{
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}
	qname := req.Question[0].Name
	names := r.names.Lookup()
	response := make(chan *dns.Msg, len(names))
	errChan := make(chan error, len(names))
	lookup := func(nameserver string) {
		res, _, err1 := c.Exchange(req, nameserver)
		fmt.Println("result:", res, err1)
		if err1 != nil {
			errChan <- fmt.Errorf("%v,%v", err, err1)
			return
		}
		if res != nil {
			if res.Rcode == dns.RcodeServerFailure {
				errChan <- fmt.Errorf("请求失败")
				return
			}
		}
		select {
		case response <- res:
		default:
		}
	}

	//循环所有名称服务器，每个服务器等待500毫秒，未拿到解析结果则发起下一个名称解析
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()

	respChan := make(chan *dns.Msg, len(names))
	for _, host := range names {
		go lookup(host)
		select {
		case re := <-response:
			respChan <- re
		case <-ticker.C: //1秒没返回则同步查询
			continue
		}
	}

	//处理返回结果
	select {
	case re := <-respChan:
		return re, nil
	case <-time.After(time.Second * 10):
		select {
		case err := <-errChan:
			return nil, err
		default:
			return nil, fmt.Errorf("无法解析的域名:%s[%v](%v)", qname, names, err)
		}
	}
}
