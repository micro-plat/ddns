package remote

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/micro-plat/hydra/global"

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
func (r *Remote) Lookup(req *dns.Msg, net string) (message *dns.Msg, err error) {
	//查询名称服务器，并处理结果
	names := r.names.Lookup()
	response := make(chan *dns.Msg, len(names))
	errChan := make(chan error, 1)
	stopChan := make(chan struct{})
	finishChan := make(chan struct{})

	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()

	wait := &sync.WaitGroup{}
	localIP := global.LocalIP()
	allLookup := func() {
		for _, host := range names {
			if strings.HasPrefix(host, localIP) {
				continue
			}
			//fmt.Println("for:", host, net)
			go func() {
				wait.Add(1)
				defer wait.Done()
				res, err := r.singleLookup(net, host, req)
				if err != nil {
					select {
					case errChan <- err:
					default:
					}
				}
				if !reflect.ValueOf(res).IsNil() {
					response <- res
				}
			}()

			select {
			case <-stopChan:
				//fmt.Println("stopChan")
				return
			case <-ticker.C:
				//fmt.Println("next", host)
				continue
			}
		}
		wait.Wait()
		close(finishChan)
	}

	go allLookup()

	//处理返回结果
	select {
	case re := <-response:
		close(stopChan)
		return re, nil
	case <-finishChan:
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

func (r *Remote) singleLookup(net string, nameserver string, req *dns.Msg) (res *dns.Msg, err error) {
	start := time.Now()
	defer func() {
		issuc := false
		if res != nil {
			issuc = len(res.Answer) > 0
		}
		fmt.Println("singleLookup:timerange(ms):", time.Now().Sub(start).Milliseconds(), net, nameserver, req.Question[0].Name, issuc, err)
	}()
	c := &dns.Client{
		Net:          net,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	res, rtt, err := c.Exchange(req, nameserver)
	if err != nil {
		//fmt.Println("singleLookup-err:", nameserver, net, err, req.Question[0].Name)
		return nil, err
	}

	//异步更新rtt
	go r.names.UpdateRTT(nameserver, rtt)
	if res != nil {
		//bytes, _ := json.Marshal(res)
		//fmt.Println("singleLookup:", nameserver, net, string(bytes))
		if res.Rcode == dns.RcodeServerFailure {
			return nil, fmt.Errorf("请求失败")
		}
	}
	if len(res.Answer) > 0 {
		return res, nil
	}
	return nil, nil
}

//Close 关闭服务
func (r *Remote) Close() {
	if r.names != nil {
		r.names.Close()
	}

}
