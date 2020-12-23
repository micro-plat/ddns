package resolver

import (
	"fmt"

	"github.com/micro-plat/ddns/local"
	"github.com/micro-plat/ddns/remote"
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/lib4go/logger"
	"github.com/miekg/dns"
)

type IResolver interface {
	Lookup(net string, req *dns.Msg) (message *dns.Msg, cache bool, err error)
}

type Resolver struct {
	log    logger.ILogger
	local  *local.Local
	remote *remote.Remote
}

func New() (*Resolver, error) {

	l, err := local.New()
	if err != nil {
		return nil, err
	}
	r, err := remote.New()
	if err != nil {
		return nil, err
	}
	re := &Resolver{
		log:    hydra.G.Log(),
		local:  l,
		remote: r,
	}
	return re, nil
}

//Lookup 从缓存和远程
func (r *Resolver) Lookup(net string, req *dns.Msg, onlyUseRemote bool) (message *dns.Msg, cache bool, count int, err error) {
	if onlyUseRemote {
		return r.LookupFromRemote(net, req)
	}
	return r.LookupFromCacheAndRemote(net, req)
}

//LookupFromRemote 从远程服务器获得解析结果
func (r *Resolver) LookupFromRemote(net string, req *dns.Msg) (message *dns.Msg, cache bool, count int, err error) {

	//查询本地缓存
	cmsg, ok := r.local.Lookup(req, false)
	if ok {
		return cmsg, true, 1, nil
	}

	//远程查询
	rmsg, count, err := r.remote.Lookup(req, net)
	if err != nil {
		return nil, false, count, err
	}

	//数据正确则保存到缓存
	if len(rmsg.Answer) > 0 {
		r.local.Save2Cache(rmsg)
	}
	return rmsg, false, count, nil
}

//LookupFromCacheAndRemote 从缓存和远程
func (r *Resolver) LookupFromCacheAndRemote(net string, req *dns.Msg) (message *dns.Msg, cache bool, count int, err error) {
	//查询本地缓存
	cmsg, ok := r.local.Lookup(req)
	if ok {
		return cmsg, true, 1, nil
	}

	//查询远程服务
	rmsg, count, err := r.remote.Lookup(req, net)
	if err != nil {
		return nil, false, count, err
	}

	//数据正确则保存到缓存
	if len(rmsg.Answer) > 0 {
		r.local.Save2Cache(rmsg)
		return rmsg, false, count, nil
	}

	//再次从缓存中拉取，解决并发请求时部分请求未能从名称服务器中获取到结果的问题
	cmsg, ok = r.local.Lookup(req)
	if ok {
		return cmsg, true, count + 1, nil
	}
	return nil, false, count + 1, fmt.Errorf("未获取到解析结果")
}

//Close 关闭上游服务
func (r *Resolver) Close() {
	if r.remote != nil {
		r.remote.Close()
	}
}
