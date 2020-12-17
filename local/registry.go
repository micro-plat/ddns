package local

import (
	"net"
	"strings"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/types"
)

//Registry 注册中心
type Registry struct {
	r             registry.IRegistry
	root          string
	domains       cmap.ConcurrentMap
	rootWatcher   watcher.IChildWatcher
	domainWatcher cmap.ConcurrentMap
	notify        chan *watcher.ChildChangeArgs
	domainDetails cmap.ConcurrentMap
	log           logger.ILogger
	closeCh       chan struct{}
}

//newRegistry 创建注册中心
func newRegistry() *Registry {
	r := &Registry{
		root:          "/dns",
		log:           hydra.G.Log(),
		r:             registry.GetCurrent(),
		domainWatcher: cmap.New(6),
		domains:       cmap.New(6),
		domainDetails: cmap.New(6),
		closeCh:       make(chan struct{}),
	}
	return r
}

//Start 启动注册中心监控
func (r *Registry) Start() (err error) {
	r.rootWatcher, err = watcher.NewChildWatcherByRegistry(r.r, []string{r.root}, r.log)
	if err != nil {
		return err
	}
	r.notify, err = r.rootWatcher.Start()
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
			r.rootWatcher.Close()
			return
		case <-r.notify:
			if err := r.load(); err != nil {
				r.log.Error(err)
			}
		}
	}
}

//Lookup 查询域名解析结果
func (r *Registry) Lookup(domain string) ([]net.IP, bool) {
	v, ok := r.domains.Get(domain)
	if !ok {
		return nil, false
	}
	ips := v.([]net.IP)
	return ips, len(ips) > 0
}

//GetDomainDetails 获取域名详情信息，返回格式为map[string]interface{}{"ddns.com",[]byte("{....}")}
func (r *Registry) GetDomainDetails() map[string]interface{} {
	return r.domainDetails.Items()
}

//CreateOrUpdate 创建或设置域名的IP信息
func (r *Registry) CreateOrUpdate(domain string, ip string, value ...string) error {
	path := registry.Join(r.root, domain, ip)
	ok, err := r.r.Exists(path)
	if err != nil {
		return err
	}
	if !ok {
		return r.r.CreatePersistentNode(path, types.GetStringByIndex(value, 0, "{}"))
	}
	return r.r.Update(path, types.GetStringByIndex(value, 0, "{}"))
}

//Load 加载所有域名的IP信息
func (r *Registry) load() error {
	//拉取所有域名
	cdomains, err := r.getAllDomains()
	if err != nil {
		return err
	}

	//清理已删除的域名
	r.domains.RemoveIterCb(func(k string, v interface{}) bool {
		//不处理，直接返回
		if ok := cdomains[k]; ok {
			return false
		}
		//移除监听
		if w, ok := r.domainWatcher.Get(k); ok {
			wc := w.(watcher.IChildWatcher)
			wc.Close()
		}
		//从缓存列表移除
		return true

	})

	//添加不存在的域名
	for domain := range cdomains {
		r.domainWatcher.SetIfAbsentCb(domain, func(input ...interface{}) (interface{}, error) {
			domain := input[0].(string)
			path := registry.Join(r.root, domain)
			w, err := watcher.NewChildWatcherByRegistry(r.r, []string{path}, r.log)
			if err != nil {
				return nil, err
			}
			notify, err := w.Start()
			if err != nil {
				return nil, err
			}
			//处理域名监控
			recv := func(domain string, notify chan *watcher.ChildChangeArgs) {
				for {
					select {
					case <-r.closeCh:
						return
					case <-notify:
						//获取所有IP列表
						if err := r.loadIP(domain); err != nil {
							r.log.Error(err)
						}
						//获取域名下所有IP的详情信息
						if err := r.loadDetail(domain); err != nil {
							r.log.Error(err)
						}
					}
				}
			}
			go recv(domain, notify)
			return notify, nil

		}, domain)
	}
	return nil

}
func (r *Registry) getAllDomains() (map[string]bool, error) {
	paths, _, err := r.r.GetChildren(r.root)
	if err != nil {
		return nil, err
	}
	m := make(map[string]bool)
	for _, v := range paths {
		m[TrimDomain(v)] = true
	}
	return m, nil
}

func (r *Registry) loadIP(domain string) error {
	path := registry.Join(r.root, domain)
	ips, _, err := r.r.GetChildren(path)
	if err != nil {
		return nil
	}
	nips := unpack(ips)
	switch {
	case len(nips) == 0:
		r.domains.Remove(domain)
	default:
		r.domains.Set(domain, nips)
	}
	return nil
}

func (r *Registry) loadDetail(domain string) error {

	//拉取域名下所有IP列表
	path := registry.Join(r.root, domain)
	ips, _, err := r.r.GetChildren(path)
	if err != nil {
		return err
	}

	//拉取每个IP对应的值
	list := make([][]byte, 0, len(ips))
	for _, ip := range ips {
		buff, _, err := r.r.GetValue(registry.Join(path, ip))
		if err != nil {
			return err
		}
		list = append(list, buff)
	}
	//保存到域名列表
	r.domainDetails.Set(domain, list)
	return nil
}

//getIPs 转换字符串为ip地址
func unpack(lst []string) []net.IP {
	ips := make([]net.IP, 0, 1)
	for _, v := range lst {
		args := strings.SplitN(v, "_", 2)
		if ip := net.ParseIP(args[0]); ip != nil {
			ips = append(ips, ip)
		}
	}

	return ips
}

//Close 关闭当前服务
func (r *Registry) Close() error {
	close(r.closeCh)
	return nil
}

//Update 更新域名的ip列表
func (r *Registry) Update(domain string, ips ...string) error {
	domainPath := registry.Join(r.root, domain)
	b, err := r.r.Exists(domainPath)
	if err != nil {
		return err
	}
	if b {
		if err := r.r.Delete(domainPath); err != nil {
			return err
		}
	}

	for _, ip := range ips {
		ippath := registry.Join(domainPath, ip)
		if err := r.r.CreatePersistentNode(ippath, "{}"); err != nil {
			return err
		}
	}
	return nil
}
