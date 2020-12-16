package pkgs

import (
 	"os"
	"time"

	"github.com/micro-plat/lib4go/types"
	"gopkg.in/fsnotify.v1"
)

type FileWatcherCallback func(string) error

type FileWatcher struct {
	closeCh   chan struct{}
	syncChan  chan string
	tickerSec int
	Change    FileWatcherCallback
	Deleted   FileWatcherCallback
	watcher   *fsnotify.Watcher
	files     types.XMap
}

func NewFileWatcher(tickerSec int) *FileWatcher {
	fw := &FileWatcher{
		tickerSec: tickerSec,
		closeCh:   make(chan struct{}),
		syncChan:  make(chan string, 100),
		files:     types.XMap{},
	}
	fw.watcher, _ = fsnotify.NewWatcher()
	go fw.loopwatch()
	return fw
}

func (fw *FileWatcher) loopwatch() {
	go fw.watchFile()
	go fw.loadFileChange()
}

func (fw *FileWatcher) Add(file string) {
	if _, ok := fw.files[file]; ok {
		return
	}
	fw.files[file] = true
	fw.watcher.Add(file)

}
func (fw *FileWatcher) Remove(file string) {
	delete(fw.files, file)
	fw.watcher.Remove(file)
}

func (fw *FileWatcher) WatchFiles() []string {
	return fw.files.Keys()
}

func (fw *FileWatcher) Close() {
	close(fw.closeCh)
	if fw.watcher != nil {
		fw.watcher.Close()
	}
}

func (fw *FileWatcher) loadFileChange() {
	period := time.Second * time.Duration(fw.tickerSec)
	ticker := time.NewTicker(period)
	for {
		select {
		case <-fw.closeCh:
			return
		case <-ticker.C:
			ticker.Stop()
			files := GetSyncData(fw.syncChan)
			if len(files) > 0 {
				files = Distinct(files)
			}

			for i := range files {
				info, err := os.Stat(files[i])
				if err != nil {
					if os.IsNotExist(err) {
						if fw.Deleted != nil {
							fw.Deleted(files[i])
						}
					}
					continue
				}
				if info.IsDir() {
					continue
				}
				if fw.Change != nil {
					fw.Change(files[i])
				}
			}
			ticker.Reset(period)
		}
	}
}

func (fw *FileWatcher) watchFile() {
	for {
		select {
		case <-fw.closeCh:
			return
		case event := <-fw.watcher.Events:
			switch event.Op {
			case fsnotify.Write, fsnotify.Remove:
				fw.syncChan <- event.Name
			default:
			}
		}
	}
}
