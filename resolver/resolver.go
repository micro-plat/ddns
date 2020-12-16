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
	log logger.ILogger
}

func New() *Resolver {
	return &Resolver{
		log: hydra.G.Log(),
	}
}

//Lookup 循环所有名称服务器，以最快速度拿取解析信息，所有名称服务器都未能成功,再次从缓存中获取
func (r *Resolver) Lookup(net string, req *dns.Msg) (message *dns.Msg, cache bool, err error) {

	//查询本地缓存
	cmsg, ok := local.Lookup(req)
	if ok {
		return cmsg, true, nil
	}

	//查询远程服务
	rmsg, ok := remote.Lookup(req)
	if !ok {
		return nil, false, fmt.Errorf("未获取到解析结果")
	}
	//保存缓存
	if len(rmsg.Answer) > 0 {
		local.Save2Cache(req)
		return rmsg, false, nil
	}
	//再次从缓存中拉取，解决并发请求时部分请求未能从名称服务器中获取到结果的问题
	cmsg, ok = local.Lookup(req)
	if ok {
		return cmsg, true, nil
	}
	return nil, true, fmt.Errorf("未获取到解析结果")
}
