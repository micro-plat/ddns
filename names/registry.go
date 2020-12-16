package names

import (
	"net"
	"strings"
	"sync"

	"github.com/micro-plat/ddns/conf"
)

//DefRegistry 管理配置中心名称服务
var DefRegistry = &Registry{names: make([]string, 0, 1)}

type Registry struct {
	names []string
	lk    sync.RWMutex
}

//Lookup 查询域名解析结果
func (f *Registry) Lookup() []string {
	f.lk.RLock()
	defer f.lk.RUnlock()
	return f.names
}

//Notify 加载注册中心配置dbs列表信息
func (f *Registry) Notify(names *conf.Names) {
	nnames := joinHostPort(names.IPS)
	f.lk.Lock()
	defer f.lk.Unlock()
	f.names = nnames
}

func joinHostPort(ips []string) []string {
	nname := make([]string, 0, len(ips))
	for _, ip := range ips {
		if strings.Contains(ip, ":") {
			nname = append(nname, ip)
			continue
		}
		nname = append(nname, net.JoinHostPort(ip, "53"))
	}
	return nname
}
