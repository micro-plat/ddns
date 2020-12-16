package conf

//Option 配置选项
type Option func(*Server)

//WithTrace 构建api server配置信息
func WithTrace() Option {
	return func(a *Server) {
		a.Trace = true
	}
}

//WithTimeout 。
func WithTimeout(rtimeout, wtimeout int) Option {
	return func(a *Server) {
		a.RTimeout = rtimeout
		a.WTimeout = wtimeout
	}
}

//WithDisable 禁用任务
func WithDisable() Option {
	return func(a *Server) {
		a.Status = StartStop
	}
}

//WithEnable 启用任务
func WithEnable() Option {
	return func(a *Server) {
		a.Status = StartStatus
	}
}

//WithUDPSize 设置udp包体大小
func WithUDPSize(UDPSize int) Option {
	return func(a *Server) {
		a.UDPSize = UDPSize
	}
}

//WithHost 设置host
func WithHost(host string) Option {
	return func(a *Server) {
		a.Host = host
	}
}

//WithPort 设置端口号
func WithPort(port string) Option {
	return func(a *Server) {
		a.Port = port
	}
}
