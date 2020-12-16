package dns

import (
	"fmt"
	"time"

	dnsconf "github.com/micro-plat/ddns/conf"
	"github.com/micro-plat/ddns/names"
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/hydra/servers"
	"github.com/micro-plat/hydra/hydra/servers/cron"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/registry/pub"
	"github.com/micro-plat/lib4go/logger"
	"github.com/miekg/dns"
)

//App dns应用程序
var App = hydra.NewApp(
	hydra.WithPlatName("ddns-test"),
	hydra.WithSystemName("ddnsserver"),
	hydra.WithUsage("DNS服务"),
	hydra.WithServerTypes(DDNS, http.API, cron.CRON, http.Web),
	hydra.WithClusterName("dns-1.2"),
	// hydra.WithRegistry("zk://192.168.0.101"),
)

func init() {
	hydra.Conf.API(":9090", api.WithTimeout(300, 300)).
		Header(header.WithCrossDomain())
}

//Server DNS服务器
type Server struct {
	conf     app.IAPPConf
	p        *Processor
	address  string
	servers  []*dns.Server
	log      logger.ILogger
	pub      pub.IPublisher
	comparer conf.IComparer
}

//NewServer 构建DNS服务器
func NewServer(cnf app.IAPPConf) (*Server, error) {
	h := &Server{
		servers:  make([]*dns.Server, 2),
		conf:     cnf,
		log:      logger.New(cnf.GetServerConf().GetServerName()),
		pub:      pub.New(cnf.GetServerConf()),
		comparer: conf.NewComparer(cnf.GetServerConf(), api.MainConfName, api.SubConfName...),
	}
	servers, process, err := h.getServer(cnf)
	if err != nil {
		return nil, err
	}
	h.servers = servers
	h.p = process
	app.Cache.Save(cnf)
	return h, nil
}

//Start 启动服务器
func (s *Server) Start() (err error) {
	s.log.Info("开始启动[DNS]服务...")
	if !s.conf.GetServerConf().IsStarted() {
		s.log.Warnf("%s被禁用，未启动", s.conf.GetServerConf().GetServerType())
		return
	}

	if len(s.servers) == 0 {
		s.log.Warnf("开启[DNS]服务器失败，没有需要启动的服务器")
		return
	}

	errChan := make(chan error, 2)
	for _, server := range s.servers {
		go func(serv *dns.Server) {
			if err := s.serve(serv); err != nil {
				errChan <- err
			}
		}(server)
	}
	select {
	case <-time.After(time.Millisecond * 500):
		if err = s.publish(); err != nil {
			err = fmt.Errorf("%s服务发布失败 %w", s.conf.GetServerConf().GetServerType(), err)
			s.Shutdown()
			return err
		}
		nnames, err := dnsconf.GetNamesConf(s.conf.GetServerConf())
		if err != nil {
			err = fmt.Errorf("获取配置中心的名称服务器失败:%v", err)
			return err
		}
		names.DefRegistry.Notify(nnames)
		s.log.Infof("服务启动成功(DNS,udp://%s)", s.address)
		s.log.Infof("服务启动成功(DNS,tcp://%s)", s.address)
		return nil
	case err := <-errChan:
		if err != nil {
			App.Close()
		}
		return err
	}

}

//Notify 配置变更后重启
func (s *Server) Notify(c app.IAPPConf) (bool, error) {
	s.comparer.Update(c.GetServerConf())
	if !s.comparer.IsChanged() {
		return false, nil
	}
	if s.comparer.IsValueChanged() || s.comparer.IsSubConfChanged() {
		s.log.Info("关键配置发生变化，准备重启服务器")
		servers, process, err := s.getServer(c)
		if err != nil {
			return false, err
		}
		s.Shutdown()
		s.conf = c
		s.p = process
		app.Cache.Save(c)
		if !c.GetServerConf().IsStarted() {
			process.Close()
			s.log.Warn("dns服务被禁用，不用重启")
			return true, nil
		}

		s.servers = servers
		if err = s.Start(); err != nil {
			process.Close()
			return false, err
		}
		return true, nil
	}
	app.Cache.Save(c)
	s.conf = c
	return true, nil
}

//Shutdown 关闭服务器
func (s *Server) Shutdown() {
	for _, server := range s.servers {
		server.Shutdown()
	}
	if s.p != nil {
		s.p.Close()
	}
	s.pub.Clear()
}

//publish 将当前服务器的节点信息发布到注册中心
func (s *Server) publish() (err error) {
	if err := s.pub.Publish(s.address, "tcp-udp://"+s.address, s.conf.GetServerConf().GetServerID()); err != nil {
		return err
	}
	return
}

func (s *Server) serve(ds *dns.Server) error {
	errChan := make(chan error, 1)
	go func(ch chan error) {
		if err := ds.ListenAndServe(); err != nil {
			ch <- err
		}
	}(errChan)
	select {
	case <-time.After(time.Millisecond * 500):
		return nil
	case err := <-errChan:
		return fmt.Errorf("DNS服务%s://%s启动失败:%v", ds.Net, ds.Addr, err)
	}
}

func (s *Server) getServer(cnf app.IAPPConf) (servers []*dns.Server, p *Processor, err error) {
	p, err = NewProcessor()
	if err != nil {
		return nil, nil, err
	}
	dnsConf, err := dnsconf.GetConf(cnf.GetServerConf())
	if err != nil {
		p.Close()
		return nil, nil, err
	}

	s.address = dnsConf.GetAddress()
	tcpHandler := dns.NewServeMux()
	tcpHandler.HandleFunc(".", p.TCP())

	udpHandler := dns.NewServeMux()
	udpHandler.HandleFunc(".", p.UDP())

	tcpServer := &dns.Server{Addr: s.address,
		Net:          "tcp",
		Handler:      tcpHandler,
		ReadTimeout:  dnsConf.GetRTimeout(),
		WriteTimeout: dnsConf.GetWTimeout()}

	udpServer := &dns.Server{Addr: s.address,
		Net:          "udp",
		Handler:      udpHandler,
		UDPSize:      dnsConf.GetUDPSize(),
		ReadTimeout:  dnsConf.GetRTimeout(),
		WriteTimeout: dnsConf.GetWTimeout()}

	servers = append(servers, tcpServer)
	servers = append(servers, udpServer)
	return servers, p, nil
}

func init() {
	fn := func(c app.IAPPConf) (servers.IResponsiveServer, error) {
		return NewServer(c)
	}
	servers.Register(DDNS, fn)
}

//DDNS 动态DNS服务器
const DDNS = "ddns"
