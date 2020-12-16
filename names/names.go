package names

import (
	"fmt"
	"net"
	"sync"

<<<<<<< HEAD:names/names.go
	"github.com/micro-plat/ddns/conf"
	"github.com/micro-plat/ddns/pkgs"
=======
	"github.com/micro-plat/ddns/dns/conf"
	"github.com/micro-plat/ddns/dns/pkgs"
>>>>>>> dev1.2-ljy:dns/names.go
	"github.com/micro-plat/hydra/conf/app"

	"github.com/micro-plat/lib4go/logger"
)

var defNames = []string{"127.0.1.1"}

//Names 本地name server读取配置
type Names struct {
	closeCh       chan struct{}
	namesCh       chan []string
	log           logger.ILogger
	names         []string
	localnames    []string
	registrynames []string
	lk            sync.RWMutex
}

//NewNames 创建本地host文件读取对象
func NewNames() *Names {
	names := &Names{
		closeCh:       make(chan struct{}),
		namesCh:       make(chan []string),
<<<<<<< HEAD:names/names.go
		log:           logger.New("names"),
=======
		log:           log,
>>>>>>> dev1.2-ljy:dns/names.go
		localnames:    make([]string, 0),
		registrynames: make([]string, 0),
	}
	return names
}

//Start 启动服务，进行本地文件读取与文件变动重新加载
func (f *Names) Start() (err error) {
	err = f.load()
	if err != nil {
		return fmt.Errorf("加载配置失败:%w", err)
	}
	f.log.Infof("[启用 NAMES,%d条]", f.len())
	go f.loopWatch()
	return nil
}

//Lookup 查询域名解析结果
func (f *Names) Lookup() []string {
	f.lk.RLock()
	defer f.lk.RUnlock()
	return f.names
}
func (f *Names) len() int {
	f.lk.RLock()
	defer f.lk.RUnlock()
	return len(f.names)
}

//Close 关闭服务
func (f *Names) Close() error {
	close(f.closeCh)
	return nil
}

func (f *Names) loopWatch() {
	go pkgs.WatchNameFile(f.closeCh, f.namesCh)
	for {
		select {
		case <-f.closeCh:
			return
		case names := <-f.namesCh:
			f.localnames = names
			f.refresh()
		}
	}
}

func (f *Names) loadlocal() error {
	fnames, err := pkgs.GetNameServers()
	if err != nil {
		return err
	}
	f.lk.Lock()
	defer f.lk.Unlock()
	f.localnames = fnames
	return nil
}

func (f *Names) loadregistry() error {
	rnames, err := f.loadRgt()
	if err != nil {
		return err
	}
	f.lk.Lock()
	defer f.lk.Unlock()
	f.registrynames = rnames
	return nil
}

func (f *Names) load() error {
	if err := f.loadlocal(); err != nil {
		return err
	}
	if err := f.loadregistry(); err != nil {
		return err
	}
	return f.refresh()
}
func (f *Names) refresh() error {
	result := make([]string, 0, len(f.localnames)+len(f.registrynames))
	result = append(result, f.localnames...)
	result = append(result, f.registrynames...)

	result = pkgs.Distinct(result)
	nNames := f.sortByTTL(result)
	for i, nm := range nNames {
		nNames[i] = net.JoinHostPort(nm, "53")
	}
	f.lk.Lock()
	defer f.lk.Unlock()
	f.names = nNames
	return nil
}
func (f *Names) sortByTTL(names []string) []string {
	sorted, err := getSortedServer(names...)
	if err != nil {
		return names
	}
	return sorted
}

//loadRgt 加载注册中心配置dbs列表信息
func (f *Names) loadRgt() ([]string, error) {
<<<<<<< HEAD:names/names.go
	//@todo: 此处引用 dns.DDNS 会循环引用
	ddnsConf, err := app.Cache.GetAPPConf("ddns")
=======
	ddnsConf, err := app.Cache.GetAPPConf(DDNS)
>>>>>>> dev1.2-ljy:dns/names.go
	if err != nil {
		return nil, fmt.Errorf("加载注册中心dns配置信息失败:%w", err)
	}

	var dnslist *conf.Names
	_, err = ddnsConf.GetServerConf().GetSubObject(conf.TypeNodeName, dnslist)
	if err != nil {
		return nil, fmt.Errorf("获取[%s]注册中心dns配置信息失败:%w", conf.TypeNodeName, err)
	}
	names := []string{}
<<<<<<< HEAD:names/names.go
	for _, str := range dnslist.IPS {
=======
	for _, str := range dnslist.Dnss {
>>>>>>> dev1.2-ljy:dns/names.go
		ip := net.ParseIP(str)
		if ip == nil {
			continue //ip格式错误
		}
		names = append(names, str)
	}
	return names, nil
}
<<<<<<< HEAD:names/names.go
=======

/*
func (f *Names) checkAndCreateConf() error {
	_, err := os.Stat(pkgs.NAME_FILE)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	fwriter, err := file.CreateFile(pkgs.NAME_FILE)
	if err != nil {
		return fmt.Errorf("创建文件:%s失败 %w", pkgs.NAME_FILE, err)
	}

	defer fwriter.Close()
	_, err = fwriter.Write([]byte(strings.Join(defNames, "\n")))
	if err != nil {
		return fmt.Errorf("写入文件:%s失败:%s", defNames, err)
	}
	return nil
}
*/
>>>>>>> dev1.2-ljy:dns/names.go
