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

//Lookup 循环所有名称服务器，以最快速度拿取解析信息，所有名称服务器都未能成功,再次从缓存中获取
func (r *Resolver) Lookup(net string, req *dns.Msg) (message *dns.Msg, cache bool, err error) {
	//查询本地缓存
	cmsg, ok := r.local.Lookup(req)
	if ok {
		return cmsg, true, nil
	}

	//查询远程服务
	rmsg, err := r.remote.Lookup(req, net)
	if err != nil {
		return nil, false, fmt.Errorf("未获取到解析结果:%w", err)
	}

	//数据正确则保存到缓存
	if len(rmsg.Answer) > 0 {
		r.local.Save2Cache(rmsg)
		return rmsg, false, nil
	}

	//再次从缓存中拉取，解决并发请求时部分请求未能从名称服务器中获取到结果的问题
	cmsg, ok = r.local.Lookup(req)
	if ok {
		return cmsg, true, nil
	}
	return nil, true, fmt.Errorf("未获取到解析结果")
}

//Close 关闭上游服务
func (r *Resolver) Close() {
	if r.local != nil {
		r.local.Close()
	}
	if r.remote != nil {
		r.remote.Close()
	}
}
