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
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}

	response := make(chan *dns.Msg, 1)
	errChan := make(chan error, 1)
	lookup := func(nameserver string) {
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
		select {
		case response <- res:
		default:
		}

	}

	//循环所有名称服务器，每个服务器等待500毫秒，未拿到解析结果则发起下一个名称解析
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	names := r.names.Lookup()
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
	select {
	case re := <-response:
		return re, nil
	case <-time.After(time.Second * 30):
		select {
		case err := <-errChan:
			return nil, err
		default:
			qname := req.Question[0].Name
			return nil, fmt.Errorf("无法解析的域名:%s[%v](%v)", qname, names, err)
		}
	}
}
