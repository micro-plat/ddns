package local

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/pub"
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
	lock          sync.RWMutex
	plats         map[string][]*Plat
	lazyClock     *time.Ticker
	lastStart     time.Time
	maxWait       time.Duration
	onceWait      time.Duration
	log           logger.ILogger
	closeCh       chan struct{}
}

//newRegistry 创建注册中心
func newRegistry() *Registry {
	r := &Registry{
		root:          hydra.G.GetDNSRoot(), //  "/dns",
		log:           hydra.G.Log(),
		r:             registry.GetCurrent(),
		plats:         make(map[string][]*Plat),
		lazyClock:     time.NewTicker(time.Second * 10),
		lastStart:     time.Now(),
		maxWait:       time.Hour,
		onceWait:      time.Second * 10,
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
	go r.lazyBuild()
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
func (r *Registry) GetDomainDetails() map[string][]*Plat {
	r.lock.RLock()
	defer r.lock.RUnlock()
	fmt.Println("plat:", len(r.plats))
	fmt.Println("raw:", r.domainDetails.Count())
	fmt.Println("domain:", r.domains.Count())
	return r.plats
}

//CreateOrUpdate 创建或设置域名的IP信息
func (r *Registry) CreateOrUpdate(domain string, ip string, value ...string) error {
	domain = TrimDomain(domain)
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
	fmt.Println("all.domains:", cdomains)
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
		fmt.Println("remove:", k)
		r.domainDetails.Remove(k)
		//从缓存列表移除
		return true
	})

	//添加不存在的域名
	for domain := range cdomains {
		_, _, err = r.domainWatcher.SetIfAbsentCb(domain, func(input ...interface{}) (interface{}, error) {
			d := input[0].(string)
			path := registry.Join(r.root, d)

			//获取所有IP列表
			fmt.Println("------load:", path)
			if err := r.loadIP(d); err != nil {
				r.log.Error(err)
			}
			//获取域名下所有IP的详情信息
			if err := r.loadDetail(d); err != nil {
				r.log.Error(err)
			}

			w, err := watcher.NewChildWatcherByRegistry(r.r, []string{path}, r.log)
			if err != nil {
				return nil, err
			}
			notify, err := w.Start()
			if err != nil {
				return nil, err
			}
			//处理域名监控
			recv := func(d string, notify chan *watcher.ChildChangeArgs) {
				for {
					select {
					case <-r.closeCh:
						return
					case <-notify:
						//获取所有IP列表
						if err := r.loadIP(d); err != nil {
							r.log.Error(err)
						}
						//获取域名下所有IP的详情信息
						if err := r.loadDetail(d); err != nil {
							r.log.Error(err)
						}
					}
				}
			}
			go recv(d, notify)
			return notify, nil

		}, domain)
		if err != nil {
			r.log.Error(err)
		}
		fmt.Println("watch:", domain)
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
		if HasWWW(v) {
			continue
		}
		m[v] = true
	}
	return m, nil
}

func (r *Registry) loadIP(domain string) error {
	path := registry.Join(r.root, domain)
	ips, _, err := r.r.GetChildren(path)
	if err != nil {
		return err
	}
	nips := unpack(ips)
	switch {
	case len(nips) == 0:
		fmt.Println("remove.domain:", domain, ips, nips, r.domains.Count())
		r.domains.Remove(domain)
	default:
		fmt.Println("add.domain:", domain, ips, nips, r.domains.Count())
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
		fmt.Println("get.detail:", registry.Join(path, ip), string(buff))
		list = append(list, buff)
	}
	//保存到域名列表
	r.domainDetails.Set(domain, list)

	//通过延迟加载的方式更新平台分组数据
	r.lastStart = time.Now()      //时钟重置
	r.lazyClock.Reset(r.onceWait) //时间置为等待一个周期
	return nil
}

//getIPs 转换字符串为ip地址
func unpack(lst []string) []net.IP {
	ips := make([]net.IP, 0, 1)
	for _, v := range lst {
		args := strings.Split(v, ":")
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
func (r *Registry) lazyBuild() {
	for {
		select {
		case <-r.closeCh:
			return
		case <-r.lazyClock.C:

			//重新构建平台分组数据
			col := make(platCollection)

			items := r.domainDetails.Items()
			fmt.Println("lazy.build:", items)
			for k, v := range items {
				fmt.Println("lazy.for:", k, v)
				list := v.([][]byte)
				for _, buff := range list {
					if err := col.append(k, buff); err != nil {
						r.log.Error(err)
					}
				}

			}
			//1个周期内没有变化，则一直等待
			if time.Since(r.lastStart) >= r.onceWait {
				r.lazyClock.Reset(r.maxWait)
			}
			r.lock.Lock()
			r.plats = col
			r.lock.Unlock()

		}
	}

}

type Plat struct {
	PlatName   string             `json:"plat_name,omitempty"`
	PlatCNName string             `json:"plat_cn_name,omitempty"`
	Clusters   map[string]*System `json:"clusters"`
}
type System struct {
	SystemName     string `json:"system_name,omitempty"`
	SystemCNName   string `json:"system_cn_name,omitempty"`
	ServerType     string `json:"server_type,omitempty"`
	ServerName     string `json:"server_name,omitempty"`
	ServiceAddress string `json:"service_address,omitempty"`
	IPAddress      string `json:"ip,omitempty"`
	URL            string `json:"url"`
}

func toPlat(r *pub.DNSConf, domain string) *Plat {
	p := &Plat{}
	p.PlatName = r.PlatName
	p.PlatCNName = r.PlatCNName
	p.Clusters = map[string]*System{
		r.ClusterName: &System{
			SystemName:     r.SystemName,
			SystemCNName:   r.SystemCNName,
			ServerType:     r.ServerType,
			ServerName:     r.ServerName,
			ServiceAddress: r.ServiceAddress,
			IPAddress:      r.IPAddress,
			URL:            GetURL(r.Proto, domain, r.Port),
		},
	}
	return p
}

//Append 将域名信息添加到列表
var defTag = "-"

type platCollection map[string][]*Plat

func (r platCollection) append(domain string, buff []byte) error {
	fmt.Println("append:", domain, types.BytesToString(buff), len(r))
	//外部注册域名
	if len(buff) == 0 || types.BytesToString(buff) == "{}" {
		if _, ok := r[defTag]; !ok {
			r[defTag] = make([]*Plat, 0, 1)
		}
		r[defTag] = append(r[defTag], &Plat{Clusters: map[string]*System{"-": &System{URL: domain}}})
		return nil
	}
	//转换服务对应的详情信息
	raw, err := pub.GetDNSConf(buff)
	if err != nil {
		return err
	}
	if raw.ServerType != http.API && raw.ServerType != http.Web {
		return nil
	}
	plat := toPlat(raw, domain)
	plats, ok := r[raw.ServerType]
	if !ok {
		r[raw.ServerType] = []*Plat{plat}
		return nil
	}

	//平台存时，将当前信息添加到指定集群
	for k, v := range plats {
		if v.PlatName == plat.PlatName {
			r[raw.ServerType][k].Clusters[raw.ClusterName] = plat.Clusters[raw.ClusterName]
			return nil
		}
	}
	//没有同名的平台，直接追加
	r[raw.ServerType] = append(r[raw.ServerType], plat)
	return nil

}
