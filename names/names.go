package names

import (
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/lib4go/logger"
)

var defNames = []string{"127.0.1.1"}

//Names 本地name server读取配置
type Names struct {
	closeCh  chan struct{}
	syncChan chan string
	watcher  *fsnotify.Watcher
	log      logger.ILogger
	names    []string
	lk       sync.RWMutex
}

//New 创建本地host文件读取对象
func New() *Names {
	names := &Names{
		closeCh:  make(chan struct{}),
		syncChan: make(chan string, 100),
		log:      hydra.G.Log(),
	}
	return names
}

//Start 启动服务，进行本地文件读取与文件变动重新加载
func (f *Names) Start() (err error) {
	// if err := f.checkAndCreateConf(); err != nil {
	// 	return err
	// }
	// f.watcher, err = fsnotify.NewWatcher()
	// if err != nil {
	// 	return fmt.Errorf("构建文件监控器失败:%w", err)
	// }
	// if err := f.watcher.Add(pkgs.NAME_FILE); err != nil {
	// 	return fmt.Errorf("添加监控文件%s失败 %w", pkgs.NAME_FILE, err)
	// }
	// err = f.reload()
	// if err != nil {
	// 	return fmt.Errorf("加载配置失败:%w", err)
	// }
	// f.log.Infof("[启用 NAMES,%d条]", f.len())
	// go f.loopWatch()
	return nil
}

//Lookup 查询域名解析结果
func (f *Names) Lookup() []string {
	return []string{"8.8.8.8:53"}
	// f.lk.RLock()
	// defer f.lk.RUnlock()
	// return f.names
}

//Close 关闭服务
func (f *Names) Close() error {
	close(f.closeCh)
	if f.watcher != nil {
		f.watcher.Close()
	}
	return nil
}

// func (f *Names) loopWatch() {
// 	go f.watchChange()
// 	go f.syncFileChange()

// }

// func (f *Names) watchChange() {

// 	for {
// 		select {
// 		case <-f.closeCh:
// 			return
// 		case event := <-f.watcher.Events:
// 			fmt.Println("xff", event.Name, event.Op)
// 			switch event.Op {
// 			case fsnotify.Write, fsnotify.Remove:
// 				f.syncChan <- event.Name
// 			default:
// 			}
// 		}
// 	}
// }

// func (f *Names) syncFileChange() {
// 	period := time.Second * 5
// 	ticker := time.NewTicker(period)
// 	for {
// 		select {
// 		case <-ticker.C:
// 			ticker.Stop()
// 			files := pkgs.GetSyncData(f.syncChan)
// 			if len(files) > 0 {
// 				files = pkgs.RemoveRepeat(files)
// 			}

// 			for i := range files {
// 				info, err := os.Stat(files[i])
// 				if err != nil {
// 					if os.IsNotExist(err) {
// 						f.log.Infof("文件%s已删除", files[i])
// 					}
// 					continue
// 				}
// 				if info.IsDir() {
// 					continue
// 				}

// 				f.log.Infof("文件%s发生变更", files[i])
// 				f.reload()
// 			}
// 			ticker.Reset(period)
// 		}
// 	}
// }

func (f *Names) reload() error {

	// names := types.XMap{}
	// fnames, err := f.load(pkgs.NAME_FILE)
	// if err != nil {
	// 	return err
	// }
	// names.Merge(fnames)
	// rnames, err := f.loadRgt()
	// if err != nil {
	// 	return err
	// }
	// names.Merge(rnames)
	// nNames := f.sortByTTL(names.Keys())
	// for i, nm := range nNames {
	// 	nNames[i] = net.JoinHostPort(nm, "53")
	// }
	// f.lk.Lock()
	// defer f.lk.Unlock()
	// f.names = nNames
	return nil

}

// func (f *Names) sortByTTL(names []string) []string {
// 	sorted, err := getSortedServer(names...)
// 	if err != nil {
// 		return names
// 	}
// 	return sorted
// }

// //load 加载配置文件，并读取指定的文件内容
// func (f *Names) load(path string) (types.XMap, error) {
// 	buf, err := os.Open(path)
// 	if err != nil {
// 		return nil, fmt.Errorf("无法打开文件%s %w", path, err)
// 	}
// 	defer buf.Close()

// 	names := types.XMap{}

// 	scanner := bufio.NewScanner(buf)
// 	for scanner.Scan() {

// 		line := pkgs.PrepareLine(scanner.Text())
// 		if strings.HasPrefix(line, "#") || line == "" {
// 			continue
// 		}

// 		ip := net.ParseIP(line)
// 		if ip == nil {
// 			continue //ip格式错误
// 		}
// 		if _, ok := names[line]; !ok {
// 			names[line] = line
// 		}
// 	}
// 	return names, nil
// }

// //loadRgt 加载注册中心配置dbs列表信息
// func (f *Names) loadRgt() ([]string, error) {
// 	appConf, err := app.Cache.GetAPPConf(dns.DDNS)
// 	if err != nil {
// 		return nil, fmt.Errorf("从缓存中获取:%w", err)
// 	}
// 	names, err := conf.GetNamesConf(appConf.GetServerConf())
// 	if err != nil {
// 		return nil, fmt.Errorf("获取[%s]注册中心dns配置信息失败:%w", conf.TypeNodeName, err)
// 	}
// 	return names.IPS, nil
// }

// func (f *Names) checkAndCreateConf() error {
// 	_, err := os.Stat(pkgs.NAME_FILE)
// 	if err == nil {
// 		return nil
// 	}
// 	if !os.IsNotExist(err) {
// 		return err
// 	}
// 	fwriter, err := file.CreateFile(pkgs.NAME_FILE)
// 	if err != nil {
// 		return fmt.Errorf("创建文件:%s失败 %w", pkgs.NAME_FILE, err)
// 	}

// 	defer fwriter.Close()
// 	_, err = fwriter.Write([]byte(strings.Join(defNames, "\n")))
// 	if err != nil {
// 		return fmt.Errorf("写入文件:%s失败:%s", defNames, err)
// 	}
// 	return nil
// }
