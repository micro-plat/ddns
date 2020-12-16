package dns

import (
	"github.com/micro-plat/ddns/resolver"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/miekg/dns"
)

//Processor cron管理程序，用于管理多个任务的执行，暂停，恢复，动态添加，移除
type Processor struct {
	*dispatcher.Engine
	resolver *resolver.Resolver
}

//NewProcessor 创建processor
func NewProcessor() (p *Processor, err error) {
	r, err := resolver.New()
	if err != nil {
		return nil, err
	}
	p = &Processor{
		resolver: r,
	}
	p.Engine = dispatcher.New()
	p.Engine.Use(middleware.Recovery().DispFunc(DDNS))
	p.Engine.Use(middleware.Logging().DispFunc())
	p.Engine.Use(middleware.Recovery().DispFunc())
	p.Engine.Use(middleware.Trace().DispFunc()) //跟踪信息
	p.Engine.Handle(DefMethod, "/", p.execute().DispFunc(DDNS))
	return p, nil
}

//TCP 处理用户请求
func (p *Processor) TCP() func(w dns.ResponseWriter, req *dns.Msg) {
	return p.Handle("tcp")
}

//UDP 处理用户请求
func (p *Processor) UDP() func(w dns.ResponseWriter, req *dns.Msg) {
	return p.Handle("udp")
}

//Handle 处理用户请求
func (p *Processor) Handle(proto string) func(w dns.ResponseWriter, req *dns.Msg) {
	return func(w dns.ResponseWriter, req *dns.Msg) {
		p.Engine.HandleRequest(NewRequest(proto, w, req))
	}
}

//ExecuteHandler 业务处理Handler
func (p *Processor) execute() middleware.Handler {
	return func(ctx middleware.IMiddleContext) {

		ctx.Log().Debug("excute")
		//处理响应
		r, _ := ctx.Request().GetMap().Get("request")
		req := r.(*dns.Msg)

		w, _ := ctx.Request().GetMap().Get("writer")
		writer := w.(dns.ResponseWriter)

		//解析域名
		msg, cache, err := p.resolver.Lookup(ctx.Request().Headers().GetString("net"), req)
		if err != nil {
			ctx.Response().WriteAny(err)
			return
		}
		if cache {
			ctx.Response().AddSpecial("C")
		}

		//处理响应结果
		msg.SetReply(req)
		writer.WriteMsg(msg)
		ctx.Response().WriteAny(msg)

	}
}
