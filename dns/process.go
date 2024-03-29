package dns

import (
	"fmt"
	"sync"

	"github.com/micro-plat/ddns/resolver"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/miekg/dns"
)

//Processor cron管理程序，用于管理多个任务的执行，暂停，恢复，动态添加，移除
type Processor struct {
	*dispatcher.Engine
	resolver      *resolver.Resolver
	metric        *middleware.Metric
	once          sync.Once
	onlyUseRemote bool
}

//NewProcessor 创建processor
func NewProcessor(onlyUseRemote bool) (p *Processor, err error) {
	r, err := resolver.New()
	if err != nil {
		return nil, err
	}
	p = &Processor{
		onlyUseRemote: onlyUseRemote,
		resolver:      r,
		metric:        middleware.NewMetric(),
	}
	p.Engine = dispatcher.New()
	p.Engine.Use(middleware.Recovery().DispFunc(DDNS))
	p.Engine.Use(middleware.Logging().DispFunc())
	p.Engine.Use(middleware.Recovery().DispFunc())
	p.Engine.Use(middleware.Trace().DispFunc()) //跟踪信息
	p.Engine.Use(p.metric.Handle().DispFunc())  //生成metric报表
	p.Engine.Handle(DefMethod, "/*name", p.execute().DispFunc(DDNS))
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
		//处理响应
		r, _ := ctx.Request().GetMap().Get("request")
		req := r.(*dns.Msg)

		w, _ := ctx.Request().GetMap().Get("writer")
		writer := w.(dns.ResponseWriter)

		//解析域名
		msg, cache, count, err := p.resolver.Lookup(ctx.Request().Headers().GetString("net"), req, p.onlyUseRemote)
		if err != nil {
			ctx.Response().WriteAny(err)
			return
		}
		ctx.Response().AddSpecial(fmt.Sprint(count))
		if cache {
			ctx.Response().AddSpecial("C")
		}

		//处理响应结果
		msg.SetReply(req)
		writer.WriteMsg(msg)
		ctx.Response().WriteAny(msg)
	}
}

//Close 关闭上游服务
func (p *Processor) Close() {
	defer p.metric.Stop()
	p.once.Do(func() {
		if p.resolver != nil {
			p.resolver.Close()
		}
	})
}
