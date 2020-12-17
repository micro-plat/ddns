package dns

import (
	"net"

	"github.com/miekg/dns"
)

const DefMethod = "GET"

//Request 处理任务请求
type Request struct {
	method string
	w      dns.ResponseWriter
	req    *dns.Msg
	form   map[string]interface{}
	header map[string]string
}

//NewRequest 构建任务请求
func NewRequest(proto string, w dns.ResponseWriter, req *dns.Msg) (r *Request) {
	r = &Request{
		method: proto,
		w:      w,
		req:    req,
		form:   make(map[string]interface{}),
		header: make(map[string]string),
	}
	r.header["net"] = proto
	r.header["Content-Type"] = "__raw__"
	if proto == "udp" {
		r.header["Client-IP"] = w.RemoteAddr().(*net.UDPAddr).IP.String()
	} else {
		r.header["Client-IP"] = w.RemoteAddr().(*net.TCPAddr).IP.String()
	}
	if proto == "udp" {
		req = req.SetEdns0(65535, true)
	}
	r.form["request"] = req
	r.form["writer"] = w
	return r
}

//GetName 获取任务名称
func (m *Request) GetName() string {
	return m.GetService()
}

//GetService 服务名
func (m *Request) GetService() string {
	if len(m.req.Question) > 0 {
		return "/" + m.req.Question[0].Name
	}
	return "/"
}

//GetMethod 方法名
func (m *Request) GetMethod() string {
	return DefMethod
}

//GetForm 输入参数
func (m *Request) GetForm() map[string]interface{} {
	return m.form
}

//GetHeader 头信息
func (m *Request) GetHeader() map[string]string {
	return m.header
}
