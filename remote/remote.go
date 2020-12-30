package remote

import (
	"fmt"
	xnet "net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/micro-plat/ddns/names"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/logger"
	"github.com/miekg/dns"
)

type Remote struct {
	names   *names.Names
	closeCh chan struct{}
}

//New 构建远程解析器
func New() (*Remote, error) {
	names, err := names.New()
	if err != nil {
		return nil, err
	}
	rmt := &Remote{
		closeCh: make(chan struct{}),
		names:   names,
	}
	return rmt, nil
}

//Lookup 从远程服务器查询解析信息
func (r *Remote) Lookup(req *dns.Msg, net string) (message *dns.Msg, count int, err error) {

	b, mes := r.checkAnalyHost(req)
	if b {
		return mes, 0, nil
	}

	//查询名称服务器，并处理结果
	names := r.names.Lookup()

	response, count, errs := r.lookupByNames(net, names, req)

	//处理返回结果
	select {
	case re, ok := <-response:
		if ok {
			return re, count, nil
		}
		if len(errs) > 0 {
			qname := req.Question[0].Name
			return nil, count, fmt.Errorf("解析的域名出错:%s(%+v)", qname, errs)
		}
	default:
	}
	qname := req.Question[0].Name
	return nil, count, fmt.Errorf("无法解析的域名:%s[%v]", qname, names)
}

func (r *Remote) lookupByNames(net string, names []string, req *dns.Msg) (chan *dns.Msg, int, []error) {

	response := make(chan *dns.Msg, len(names))
	errList := make([]error, 0, 1)

	msgChan := make(chan struct{}, len(names))
	stopCh := make(chan struct{})
	var once sync.Once
	ticker := time.NewTicker(time.Millisecond * 300)
	var count int32
	var isClose = false

	stop := func() {
		once.Do(func() { //关闭消息队列，和时钟
			isClose = true
			close(msgChan)
			close(stopCh)
			ticker.Stop()
			close(response)
		})
	}

	//发送首个信号
	msgChan <- struct{}{}
	//启动指定协程，收到指令后启动任务
	for _, host := range names {
		go func(h string, logger logger.ILogger) {
			_, ok := <-msgChan //等待启动指令
			if !ok {
				stop()
				return
			}
			res, err := r.singleLookup(net, h, req, logger)
			if err != nil { //发生错误
				errList = append(errList, err)
			} else {
				if !isClose {
					response <- res
				}
				stop()
			}
			//所有任务已执行完成
			if atomic.AddInt32(&count, 1) == int32(len(names)) {
				stop()
			}
		}(host, context.Current().Log())
	}

	//timer 定时向消息队列放入一个任务
loop:
	for {
		select {
		case <-r.closeCh:
			stop()
			return nil, int(count), nil
		case <-stopCh:
			break loop
		case _, ok := <-ticker.C:
			if !ok {
				break loop
			}
			select {
			case msgChan <- struct{}{}:
			default:
			}

		}
	}
	return response, int(count), errList
}

func (r *Remote) singleLookup(net string, nameserver string, req *dns.Msg, log logger.ILogger) (res *dns.Msg, err error) {

	log.Info("  -->exchange.request:", req.Question[0].Name, nameserver)
	start := time.Now()
	defer func() {
		hasAnswer := res != nil && len(res.Answer) > 0

		timerange := time.Since(start)
		if err != nil {
			log.Error("  -->exchange.response:", hasAnswer, timerange, req.Question[0].Name, nameserver, "err:", err)
			return
		}
		rcode := 999
		if res != nil {
			rcode = res.Rcode
		}
		log.Info("  -->exchange.response:", hasAnswer, timerange, req.Question[0].Name, nameserver, "OK", rcode)

	}()
	res, err = dns.Exchange(req, nameserver)
	if err != nil {
		return nil, err
	}
	if len(res.Answer) > 0 {
		//当存在结果时，异步更新rtt
		go r.names.UpdateRTT(nameserver, time.Since(start))
	}
	if res != nil {
		if res.Rcode == dns.RcodeServerFailure {
			return nil, fmt.Errorf("请求失败[%s]:%d %s", nameserver, res.Rcode, dns.RcodeToString[res.Rcode])
		}
	}
	return res, nil
}

//如果被解析的地址就是本地ip  那么就直接返回本机ip作为解析结果
func (r *Remote) checkAnalyHost(req *dns.Msg) (b bool, message *dns.Msg) {
	b = false
	localIP := global.LocalIP()
	if strings.HasPrefix(req.Question[0].Name, localIP) {
		b = true
		message = &dns.Msg{}
		message.Id = req.Id
		message.Question = req.Question
		header := dns.RR_Header{
			Name:   req.Question[0].Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    600,
		}
		ip := xnet.ParseIP(localIP)
		if ip != nil {
			message.Answer = append(message.Answer, &dns.A{header, ip})
		}
		return
	}

	return
}

//Close 关闭服务
func (r *Remote) Close() {
	close(r.closeCh)
	if r.names != nil {
		r.names.Close()
	}
}
