package conf

import (
	"fmt"
	xnet "net"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/global"
)

const (
	//StartStatus 开启服务
	StartStatus = "start"
	//StartStop 停止服务
	StartStop = "stop"
)

//MainConfName 主配置中的关键配置名
var MainConfName = []string{"host", "status", "rTimeout", "wTimeout", "port", "dn"}

//SubConfName 子配置中的关键配置名
var SubConfName = []string{TypeNodeName}

//Server api server配置信息
type Server struct {
	Status   string `json:"status,omitempty" valid:"in(start|stop)" toml:"status,omitempty"`
	RTimeout int    `json:"rTimeout,omitempty" toml:"rTimeout,omitzero"` //单位秒
	WTimeout int    `json:"wTimeout,omitempty" toml:"wTimeout,omitzero"` //单位秒
	UDPSize  int    `json:"udpSize,omitempty" toml:"udpSize,omitzero"`   //udp协议传输mesgges大小 单位字节
	Trace    bool   `json:"trace,omitempty" toml:"trace,omitempty"`
}

//New 构建websocket server配置信息
func New(opts ...Option) *Server {
	a := &Server{
		Status:   StartStatus,
		RTimeout: 5,
		WTimeout: 5,
		UDPSize:  65535,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

//GetAddress 获取dns服务地址端口
func (s *Server) GetAddress() string {
	return xnet.JoinHostPort(global.LocalIP(), "53")
}

//GetRTimeout 获取读取超时时间
func (s *Server) GetRTimeout() time.Duration {
	if s.RTimeout <= 0 {
		return 5 * time.Second
	}

	return time.Duration(s.RTimeout) * time.Second
}

//GetWTimeout 获取写超时时间
func (s *Server) GetWTimeout() time.Duration {
	if s.WTimeout <= 0 {
		return 5 * time.Second
	}
	return time.Duration(s.WTimeout) * time.Second
}

//GetUDPSize 获取写超时时间
func (s *Server) GetUDPSize() int {
	if s.UDPSize <= 0 {
		return 65535
	}
	return s.UDPSize
}

//GetConf 获取主配置信息
func GetConf(cnf conf.IServerConf) (s *Server, err error) {
	_, err = cnf.GetMainObject(&s)
	if err == conf.ErrNoSetting {
		return nil, fmt.Errorf("/%s :%w", cnf.GetServerPath(), err)
	}
	if err != nil {
		return nil, err
	}
	if b, err := govalidator.ValidateStruct(s); !b {
		return nil, fmt.Errorf("dds服务器主配置数据有误:%v", err)
	}
	return s, nil
}
