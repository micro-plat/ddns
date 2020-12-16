package names

import (
	"sync"

	"github.com/micro-plat/ddns/pkgs"
	"github.com/micro-plat/hydra"

	"github.com/micro-plat/lib4go/logger"
)

//ResolveConf 本地name server读取配置
type ResolveConf struct {
	closeCh chan struct{}
	namesCh chan []string
	log     logger.ILogger
	names   []string
	lk      sync.RWMutex
}

//NewResolveConf 创建本地host文件读取对象
func NewResolveConf() *ResolveConf {
	names := &ResolveConf{
		closeCh: make(chan struct{}),
		namesCh: make(chan []string),
		log:     hydra.G.Log(),
		names:   make([]string, 0),
	}
	return names
}

//Start 启动服务，进行本地文件读取与文件变动重新加载
func (f *ResolveConf) Start() (err error) {
	if err = f.load(); err != nil {
		return err
	}
	go f.loopWatch()
	return nil
}

//Lookup 查询域名解析结果
func (f *ResolveConf) Lookup() []string {
	f.lk.RLock()
	defer f.lk.RUnlock()
	return f.names
}

//Close 关闭服务
func (f *ResolveConf) Close() error {
	close(f.closeCh)
	return nil
}

func (f *ResolveConf) loopWatch() {
	go pkgs.WatchNameFile(f.closeCh, f.namesCh)
	for {
		select {
		case <-f.closeCh:
			return
		case names := <-f.namesCh:
			f.refresh(names)
		}
	}
}

func (f *ResolveConf) load() error {
	names, err := pkgs.GetNameServers()
	if err != nil {
		return err
	}
	return f.refresh(names)
}
func (f *ResolveConf) refresh(result []string) error {
	names := joinHostPort(result)
	f.lk.Lock()
	defer f.lk.Unlock()
	f.names = names
	return nil
}
