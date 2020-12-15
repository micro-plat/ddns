package query

import (
	"bufio"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/micro-plat/lib4go/logger"
)

//Hosts 本地Host读取配置
type Hosts struct {
	closeCh chan struct{}
	watcher *fsnotify.Watcher
	log     logger.ILogger
	domain  map[string]map[string][]net.IP
	lk      sync.Mutex
}

//NewHosts 创建本地host文件读取对象
func NewHosts(log logger.ILogger) *Hosts {

	hosts := &Hosts{
		closeCh: make(chan struct{}),
		log:     log,
	}
	return hosts
}

//Start 启动服务，进行本地文件读取与文件变动重新加载
func (f *Hosts) Start() (err error) {
	f.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if err := f.watcher.Add(HOST_FILE); err != nil {
		return err
	}
	err = f.loadAll()
	if err != nil {
		return err
	}
	f.log.Infof("[启用 HOSTS,%d条]", f.len())
	go f.loopWatch()
	return nil
}

//Lookup 查询域名解析结果
func (f *Hosts) Lookup(name string) []net.IP {
	f.lk.Lock()
	defer f.lk.Unlock()
	for _, domain := range f.domain {
		if ips, ok := domain[name]; ok {
			return ips
		}
	}
	return nil
}
func (f *Hosts) len() int {
	f.lk.Lock()
	defer f.lk.Unlock()
	count := 0
	for _, domain := range f.domain {
		count += len(domain)
	}
	return count
}

//Close 关闭服务
func (f *Hosts) Close() error {
	close(f.closeCh)
	if f.watcher != nil {
		f.watcher.Close()
	}
	return nil
}

func (f *Hosts) loopWatch() {
	for {
		select {
		case <-f.closeCh:
			return
		case event := <-f.watcher.Events:
			if strings.HasSuffix(event.Name, ".swp") ||
				strings.HasSuffix(event.Name, ".swx") ||
				strings.HasSuffix(event.Name, ".swpx") ||
				strings.HasPrefix(event.Name, "~") ||
				strings.HasSuffix(event.Name, "~") ||
				!strings.Contains(event.Name, "hosts") {
				continue
			}
			switch event.Op {
			case fsnotify.Write:
				f.log.Infof("文件%s发生变更", event.Name)
				if err := f.reloadOne(event.Name); err != nil {
					f.log.Error(err)
				}
				f.log.Infof("[启用 HOSTS,%d]", f.len())
			case fsnotify.Remove:
				_, err := os.Stat(event.Name)
				if err == nil {
					continue
				}
				f.log.Infof("文件%s已删除", event.Name)
				if err := f.removeOne(event.Name); err != nil {
					f.log.Error(err)
				}
				f.log.Infof("[启用 HOSTS,%d]", f.len())
			default:

			}
		}
	}
}

func (f *Hosts) loadAll() error {
	files, err := filepath.Glob(HOST_FILE)
	if err != nil {
		return err
	}
	ndomain := make(map[string]map[string][]net.IP)
	for _, file := range files {
		domain, err := f.load(file)
		if err != nil {
			return err
		}
		ndomain[file] = domain
		f.watcher.Add(file)
	}
	f.lk.Lock()
	defer f.lk.Unlock()
	f.domain = ndomain
	return nil
}
func (f *Hosts) removeOne(file string) error {
	f.lk.Lock()
	defer f.lk.Unlock()
	delete(f.domain, file)
	f.watcher.Remove(file)
	return nil
}

func (f *Hosts) reloadOne(file string) error {
	domain, err := f.load(file)
	if err != nil {
		return err
	}
	f.lk.Lock()
	defer f.lk.Unlock()
	f.domain[file] = domain
	if _, ok := f.domain[file]; !ok {
		f.watcher.Add(file)
	}
	return nil
}

//load 加载配置文件，并读取指定的文件内容
func (f *Hosts) load(path string) (map[string][]net.IP, error) {
	buf, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer buf.Close()

	hosts := make(map[string][]net.IP)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {

		line := strings.Replace(strings.TrimSpace(scanner.Text()), "\t", " ", -1)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		sli := strings.Split(line, " ")
		if len(sli) < 2 {
			continue
		}

		ip := net.ParseIP(sli[0])
		if ip == nil {
			continue //ip格式错误
		}

		for i := 1; i <= len(sli)-1; i++ {
			domain := strings.ToLower(strings.TrimSpace(sli[i]))
			if domain == "" {
				continue
			}
			if _, ok := hosts[domain]; !ok {
				hosts[domain] = make([]net.IP, 0, 1)
			}
			hosts[domain] = append(hosts[domain], ip)
		}
	}
	return hosts, nil
}
