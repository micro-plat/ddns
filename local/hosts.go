package local

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/micro-plat/ddns/pkgs"
	"github.com/micro-plat/lib4go/logger"
)

//Hosts 本地Host读取配置
type Hosts struct {
	closeCh  chan struct{}
	syncChan chan string
	watcher  *fsnotify.Watcher
	log      logger.ILogger
	domain   map[string]map[string][]net.IP
	lk       sync.RWMutex
}

//NewHosts 创建本地host文件读取对象
func NewHosts(log logger.ILogger) *Hosts {

	hosts := &Hosts{
		closeCh:  make(chan struct{}),
		syncChan: make(chan string, 100),
		log:      log,
	}
	return hosts
}

//Start 启动服务，进行本地文件读取与文件变动重新加载
func (f *Hosts) Start() (err error) {
	f.watcher, err = fsnotify.NewWatcher()
	if err != nil {
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
	f.lk.RLock()
	defer f.lk.RUnlock()
	for _, domain := range f.domain {
		if ips, ok := domain[name]; ok {
			return ips
		}
	}
	return nil
}
func (f *Hosts) len() int {
	f.lk.RLock()
	defer f.lk.RUnlock()
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

func (f *Hosts) loadAll() error {
	files, err := filepath.Glob(pkgs.HOST_FILE)
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
	}
	f.lk.Lock()
	defer f.lk.Unlock()

	for file := range ndomain {
		f.watcher.Add(file)
	}
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
	if _, ok := f.domain[file]; !ok {
		f.watcher.Add(file)
	}
	f.domain[file] = domain
	return nil
}

//load 加载配置文件，并读取指定的文件内容
func (f *Hosts) load(path string) (map[string][]net.IP, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("path:%s is directory", path)
	}

	buf, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer buf.Close()

	hosts := make(map[string][]net.IP)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {

		line := pkgs.PrepareLine(scanner.Text())
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

func (f *Hosts) loopWatch() {
	go f.watchNewFile()
	go f.watchChange()
	go f.syncFileChange()
}

func (f *Hosts) watchNewFile() {
	period := time.Second * 5
	ticker := time.NewTicker(period)
	for {
		select {
		case <-f.closeCh:
			ticker.Stop()
			return
		case <-ticker.C:
			ticker.Stop()
			files, err := filepath.Glob(pkgs.HOST_FILE)
			if err != nil {
				f.log.Errorf("filepath.Glob:%s;%w", pkgs.HOST_FILE, err)
				break
			}
			for i := range files {
				if _, ok := f.domain[files[i]]; !ok {
					err = f.reloadOne(files[i])
					if err != nil {
						f.log.Errorf("reloadOne:%s;%w", files[i], err)
					}

				}
			}
			ticker.Reset(period)
		}
	}
}

func (f *Hosts) watchChange() {
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
				!strings.HasPrefix(filepath.Base(event.Name), "hosts") {
				continue
			}
			switch event.Op {
			case fsnotify.Write, fsnotify.Remove:
				f.syncChan <- event.Name
			default:
			}
		}
	}
}

func (f *Hosts) syncFileChange() {
	period := time.Second * 5
	ticker := time.NewTicker(period)
	for {
		select {
		case <-ticker.C:
			ticker.Stop()
			files := pkgs.GetSyncData(f.syncChan)
			if len(files) > 0 {
				files = pkgs.Distinct(files)
			}

			for i := range files {
				info, err := os.Stat(files[i])
				if err != nil {
					if os.IsNotExist(err) {
						f.log.Infof("文件%s已删除", files[i])
						f.removeOne(files[i])
						continue
					}
					continue
				}
				if info.IsDir() {
					continue
				}

				f.log.Infof("文件%s发生变更", files[i])
				f.reloadOne(files[i])
			}
			ticker.Reset(period)
		}
	}
}
