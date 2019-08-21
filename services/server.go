package services

import (
	"net"
	"time"

	"github.com/micro-plat/lib4go/logger"
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
}

func NewServer() *Server {
	return &Server{
		servers:  make([]*dns.Server, 0, 2),
		host:     "127.0.0.1",
		port:     "53",
		rTimeout: 3 * time.Second,
		wTimeout: 3 * time.Second,
		log:      logger.New("ddns"),
	}
}

//Addr 获取服务器地址
func (s *Server) Addr() string {
	return net.JoinHostPort(s.host, s.port)
}

//Start 启动服务器
func (s *Server) Start() {
	handler := NewHandler(s.log)

	tcpHandler := dns.NewServeMux()
	tcpHandler.HandleFunc(".", handler.Do("tcp"))

	udpHandler := dns.NewServeMux()
	udpHandler.HandleFunc(".", handler.Do("udp"))

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

	go s.serve(udpServer)
	go s.serve(tcpServer)
	s.servers = append(s.servers, tcpServer)
	s.servers = append(s.servers, udpServer)

}

//Shutdown 关闭服务器
func (s *Server) Shutdown() {
	for _, server := range s.servers {
		server.Shutdown()
	}
}
func (s *Server) serve(ds *dns.Server) {
	s.log.Infof("启动%s服务:%s", ds.Net, ds.Addr)
	err := ds.ListenAndServe()
	if err != nil {
		s.log.Errorf("服务%s-%s启动失败:%v", ds.Net, ds.Addr, err.Error())
	}
}
