package remote

import (
	"fmt"
	"sync"
	"time"

	"github.com/micro-plat/ddns/names"
	"github.com/micro-plat/lib4go/logger"
	"github.com/miekg/dns"
)

var defRemote = New()

func Lookup(req *dns.Msg) (message *dns.Msg, b bool) {
	return nil, false
}

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
	logger := logger.New("ctx")
	qname := req.Question[0].Name
	response := make(chan *dns.Msg, 1)
	var wg sync.WaitGroup
	lookup := func(nameserver string) {
		defer wg.Done()
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

	//循环所有名称服务器，每个服务器等待500毫秒，未拿到解析结果则发起下一个名称解析
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	names := r.names.Lookup()
	for _, host := range names {
		wg.Add(1)
		go lookup(host)
		select {
		case re := <-response:
			return re, nil
		case <-ticker.C: //1秒没返回则同步查询
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
