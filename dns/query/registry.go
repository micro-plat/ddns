package query

import (
	"net"
	"strings"
	"sync"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/lib4go/logger"
)

//Registry 注册中心
type Registry struct {
	r       registry.IRegistry
	root    string
	domain  map[string][]net.IP
	watcher watcher.IChildWatcher
	notify  chan *watcher.ChildChangeArgs
	log     logger.ILogger
	closeCh chan struct{}
	lk      sync.Mutex
}

//NewRegistry 创建注册中心
func NewRegistry(r registry.IRegistry, log logger.ILogger) *Registry {
	registry := &Registry{
		r:       r,
		root:    "/dns",
		domain:  make(map[string][]net.IP),
		log:     log,
		closeCh: make(chan struct{}),
	}
	registry.watcher, _ = watcher.NewChildWatcherByRegistry(r, []string{"/dns"}, log)
	return registry
}

//Start 启动注册中心监控
func (r *Registry) Start() (err error) {
	r.notify, err = r.watcher.Start()
	if err != nil {
		return err
	}
	go r.loopWatch()
	return nil
}
func (r *Registry) loopWatch() {
	for {
		select {
		case <-r.closeCh:
			return
		case n := <-r.notify:
			if err := r.load(n.Parent, n.Name); err != nil {
				r.log.Error(err)
			}

		}
	}
}

//Lookup 查询域名解析结果
func (r *Registry) Lookup(name string) []net.IP {
	r.lk.Lock()
	defer r.lk.Unlock()
	return r.domain[name]
}

//Close 关闭当前服务
func (r *Registry) Close() error {
	close(r.closeCh)
	r.watcher.Close()
	return nil
}

//Load 加载所有域名的IP信息
func (r *Registry) load(path string, name string) error {
	if b, err := r.r.Exists(path); !b && err == nil {
		r.lk.Lock()
		delete(r.domain, name)
		r.lk.Unlock()
		r.log.Infof("[缓存:%s,0条]", path)
		return nil
	}
	ips, _, err := r.r.GetChildren(path)
	if err != nil {
		return nil
	}
	nips := getIPs(ips)
	//修改本地域名缓存
	r.lk.Lock()
	switch {
	case len(nips) == 0:
		delete(r.domain, name)
	default:
		r.domain[name] = nips
	}
	r.lk.Unlock()
	r.log.Infof("[缓存:%s,%d条]", name, len(nips))
	return nil
}

//getIPs 转换字符串为ip地址
func getIPs(lst []string) []net.IP {
	ips := make([]net.IP, 0, len(lst))
	for _, v := range lst {
		args := strings.SplitN(v, "_", 2)
		ips = append(ips, net.ParseIP(args[0]))
	}
	return ips
}
func (r *Registry) len() int {
	r.lk.Lock()
	defer r.lk.Unlock()
	count := 0
	for _, domain := range r.domain {
		count += len(domain)
	}
	return count
}
