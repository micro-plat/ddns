package query

import (
	"net"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/lib4go/logger"
)

//Registry 注册中心
type Registry struct {
	r       registry.IRegistry
	root    string
	domain  map[string][]net.IP
	watcher *watcher.Watcher
	notify  chan *watcher.ContentChangeArgs
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
		watcher: watcher.NewWatcher("/dns", time.Second*10, r, log),
	}
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
		case <-r.notify:
			r.log.Infof("%s发生变更", r.root)
			if err := r.load(); err != nil {
				r.log.Error(err)
			}
			r.log.Infof("[启用 注册中心,%d条]", len(r.domain))
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
func (r *Registry) load() error {
	ndomain := make(map[string][]net.IP)
	domains, _, err := r.r.GetChildren(r.root)
	if err != nil {
		return err
	}
	for _, d := range domains {
		ips, _, err := r.r.GetChildren(filepath.Join(r.root, d))
		if err != nil {
			return err
		}
		ndomain[d] = getIPs(ips)
	}
	//修改本地域名缓存
	r.lk.Lock()
	r.domain = ndomain
	r.lk.Unlock()
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
