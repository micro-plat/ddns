package dns

import (
	"sync"
	"time"

	"github.com/micro-plat/hydra/components/queues/mq"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/miekg/dns"
)

//Processor cron管理程序，用于管理多个任务的执行，暂停，恢复，动态添加，移除
type Processor struct {
	*dispatcher.Engine
	lock      sync.Mutex
	done      bool
	closeChan chan struct{}
	queues    cmap.ConcurrentMap
	startTime time.Time
	customer  mq.IMQC
}

//NewProcessor 创建processor
func NewProcessor() (p *Processor) {
	p = &Processor{
		closeChan: make(chan struct{}),
		startTime: time.Now(),
		queues:    cmap.New(4),
	}
	p.Engine = dispatcher.New()
	p.Engine.Use(middleware.Recovery().DispFunc(DDNS))
	p.Engine.Use(middleware.Logging().DispFunc())
	p.Engine.Use(middleware.Recovery().DispFunc())
	p.Engine.Use(middleware.Trace().DispFunc()) //跟踪信息
	p.Engine.Handle("GET", "*", p.execute().DispFunc(DDNS))
	return p
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
		//lockup......
		//处理响应
		r, _ := ctx.Request().GetMap().Get("req")
		req := r.(*dns.Msg)

		w, _ := ctx.Request().GetMap().Get("w")
		writer := w.(dns.ResponseWriter)

		msg.SetReply(req)
		writer.WriteMsg(msg)
		ctx.Response().Write(msg)

	}
}
