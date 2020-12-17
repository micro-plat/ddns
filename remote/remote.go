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
	ctrlChan := make(chan struct{})

	go r.allLookup(names, ctrlChan, req, response, errChan)

	//处理返回结果
	select {
	case re := <-response:
		close(ctrlChan)
		return re, nil
	case <-ctrlChan:
		select {
		case re := <-response:
			return re, nil
		case err := <-errChan:
			return nil, err
		default:
			qname := req.Question[0].Name
			return nil, fmt.Errorf("无法解析的域名:%s[%v](%v)", qname, names, err)
		}
	}
}

func (r *Remote) allLookup(names []string, ctrlChan chan struct{}, req *dns.Msg, response chan *dns.Msg, errChan chan error) {
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	wait := &sync.WaitGroup{}
	for _, host := range names {
		go r.singleLookup(wait, host, req, response, errChan)
		select {
		case <-ctrlChan:
			return
		case <-ticker.C:
			continue
		}
	}
	wait.Wait()
	close(ctrlChan)
}
func (r *Remote) singleLookup(wait *sync.WaitGroup, nameserver string, req *dns.Msg, response chan *dns.Msg, errChan chan error) {
	c := &dns.Client{
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	wait.Add(1)
	defer wait.Done()
	res, rtt, err := c.Exchange(req, nameserver)
	if err != nil {
		select {
		case errChan <- err:
		default:
		}
		return
	}

	//异步更新rtt
	go r.names.UpdateRTT(nameserver, rtt)
	if res != nil {
		if res.Rcode == dns.RcodeServerFailure {
			select {
			case errChan <- fmt.Errorf("请求失败"):
			default:
			}
			return
		}
	}
	if len(res.Answer) > 0 {
		response <- res
	}
	return
}

//Close 关闭服务
func (r *Remote) Close() {
	if r.names != nil {
		r.names.Close()
	}

}
