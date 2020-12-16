package remote

import (
	"fmt"
	"sync"
	"time"

	"github.com/micro-plat/ddns/names"
	"github.com/miekg/dns"
)

type Remote struct {
	names *names.Names
}

//New 构建远程解析器
func New() (*Remote, error) {
	names, err := names.New()
	if err != nil {
		return nil, err
	}
	rmt := &Remote{
		names: names,
	}
	return rmt, nil
}

//Lookup 从远程服务器查询解析信息
func (r *Remote) Lookup(req *dns.Msg) (message *dns.Msg, err error) {
	//查询名称服务器，并处理结果
	c := &dns.Client{
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	names := r.names.Lookup()
	response := make(chan *dns.Msg, len(names))
	errChan := make(chan error, 1)
	wait := sync.WaitGroup{}
	lookup := func(nameserver string) {
		wait.Add(1)
		defer wait.Done()
		res, _, err1 := c.Exchange(req, nameserver)
		if err1 != nil {
			select {
			case errChan <- err1:
			default:
			}
			return
		}
		if res != nil {
			if res.Rcode == dns.RcodeServerFailure {
				select {
				case errChan <- fmt.Errorf("请求失败"):
				default:
				}
				return
			}
		}
		response <- res
	}

	//循环所有名称服务器，每个服务器等待500毫秒，未拿到解析结果则发起下一个名称解析
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()

LOOP:
	for _, host := range names {
		go lookup(host)
		select {
		case re := <-response:
			select {
			case response <- re:
			default:
			}
			break LOOP
		case <-ticker.C:
			continue
		}
	}

	//处理返回结果
	timeout := make(chan struct{})
	go func() {
		wait.Wait()
		close(timeout)
	}()
	for {
		select {
		case re := <-response:
			if len(re.Answer) > 0 {
				return re, nil
			}
		case <-timeout:
			select {
			case re := <-response:
				if len(re.Answer) > 0 {
					return re, nil
				}
			case err := <-errChan:
				return nil, err
			default:
				qname := req.Question[0].Name
				return nil, fmt.Errorf("无法解析的域名:%s[%v](%v)", qname, names, err)
			}

		}
	}
}

//Close 关闭服务
func (r *Remote) Close() {
	if r.names != nil {
		r.names.Close()
	}

}
