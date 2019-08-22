package dns

import (
	"fmt"
	"net"
	"time"

	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/lib4go/logger"
	xnet "github.com/micro-plat/lib4go/net"
	"github.com/miekg/dns"
)

//Server DNS服务器
type Server struct {
	host     string
	port     string
	servers  []*dns.Server
	rTimeout time.Duration
	wTimeout time.Duration
	log      logger.ILogger
	c        component.IContainer
	hander   *DNSHandler
}

func NewServer(c component.IContainer) *Server {
	return &Server{
		c:        c,
		servers:  make([]*dns.Server, 0, 2),
		host:     xnet.GetLocalIPAddress(),
		port:     "53",
		rTimeout: 5 * time.Second,
		wTimeout: 5 * time.Second,
		log:      logger.New("ddns"),
	}
}

//Addr 获取服务器地址
func (s *Server) Addr() string {
	return net.JoinHostPort(s.host, s.port)
}

//Start 启动服务器
func (s *Server) Start() (err error) {
	s.log.Info("开始启动[DNS]服务...")
	s.hander, err = NewHandler(s.c, s.log)
	if err != nil {
		return err
	}

	tcpHandler := dns.NewServeMux()
	tcpHandler.HandleFunc(".", s.hander.Do("tcp"))

	udpHandler := dns.NewServeMux()
	udpHandler.HandleFunc(".", s.hander.Do("udp"))

	tcpServer := &dns.Server{Addr: s.Addr(),
		Net:          "tcp",
		Handler:      tcpHandler,
		ReadTimeout:  s.rTimeout,
		WriteTimeout: s.wTimeout}

	udpServer := &dns.Server{Addr: s.Addr(),
		Net:          "udp",
		Handler:      udpHandler,
		UDPSize:      65535,
		ReadTimeout:  s.rTimeout,
		WriteTimeout: s.wTimeout}

	if err := s.serve(udpServer); err != nil {
		return err
	}
	s.log.Infof("服务启动成功(DNS,udp://%s)", s.Addr())
	if err := s.serve(tcpServer); err != nil {
		return err
	}
	s.log.Infof("服务启动成功(DNS,tcp://%s)", s.Addr())
	s.servers = append(s.servers, tcpServer)
	s.servers = append(s.servers, udpServer)
	return nil

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
