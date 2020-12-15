package dns

import (
	"fmt"
	"time"

	dnsconf "github.com/micro-plat/ddns/dns/conf"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/hydra/servers"
	"github.com/micro-plat/hydra/registry/pub"
	"github.com/micro-plat/lib4go/logger"
	"github.com/miekg/dns"
)

//Server DNS服务器
type Server struct {
	conf     app.IAPPConf
	address  string
	servers  []*dns.Server
	rTimeout time.Duration
	wTimeout time.Duration
	log      logger.ILogger
	pub      pub.IPublisher
	hander   *DNSHandler
	comparer conf.IComparer
}

//NewServer 构建DNS服务器
func NewServer(cnf app.IAPPConf) (*Server, error) {
	h := &Server{
		conf:     cnf,
		log:      logger.New(cnf.GetServerConf().GetServerName()),
		pub:      pub.New(cnf.GetServerConf()),
		comparer: conf.NewComparer(cnf.GetServerConf(), api.MainConfName, api.SubConfName...),
	}
	servers, err := h.getServer(cnf)
	if err != nil {
		return nil, err
	}
	h.servers = servers
	app.Cache.Save(cnf)
	return h, nil
}

//Start 启动服务器
func (s *Server) Start() (err error) {
	s.log.Info("开始启动[DNS]服务...")
	if len(s.servers) == 0 {
		s.log.Warnf("开启[DNS]服务器失败，没有需要启动的服务器")
		return
	}
	errChan := make(chan error, 2)
	for _, server := range s.servers {
		go func() {
			if err := s.serve(server); err != nil {
				errChan <- err
			}
		}()
	}
	select {
	case <-time.After(time.Millisecond * 500):
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
		servers, err := s.getServer(c)
		if err != nil {
			return false, err
		}

		s.Shutdown()
		s.conf = c
		app.Cache.Save(c)
		if !c.GetServerConf().IsStarted() {
			s.log.Info("dns服务被禁用，不用重启")
			return true, nil
		}

		s.servers = servers
		if err = s.Start(); err != nil {
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
	if s.hander != nil {
		s.hander.Close()
	}
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

func (s *Server) getServer(cnf app.IAPPConf) (servers []*dns.Server, err error) {
	dnsConf, err := dnsconf.GetConf(cnf.GetServerConf())
	if err != nil {
		return nil, err
	}

	s.address = dnsConf.GetAddress()
	s.hander, err = NewHandler(s.log)
	if err != nil {
		App.Close()
		return nil, fmt.Errorf("构建DNS处理服务失败:%w", err)
	}

	tcpHandler := dns.NewServeMux()
	tcpHandler.HandleFunc(".", s.hander.Do("tcp"))

	udpHandler := dns.NewServeMux()
	udpHandler.HandleFunc(".", s.hander.Do("udp"))

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
	return servers, nil
}

func init() {
	fn := func(c app.IAPPConf) (servers.IResponsiveServer, error) {
		return NewServer(c)
	}
	servers.Register(DDNS, fn)
}

//DDNS 动态DNS服务器
const DDNS = "ddns"
