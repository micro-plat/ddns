package names

import (
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/micro-plat/lib4go/concurrent/cmap"
)

type Sorter struct {
	sorted  []string
	rtts    cmap.ConcurrentMap
	lock    sync.RWMutex
	closeCh chan struct{}
	index   int32
}

func newSorter() *Sorter {
	r := &Sorter{
		sorted:  make([]string, 0, 1),
		rtts:    cmap.New(3),
		closeCh: make(chan struct{}),
	}
	go r.loopUpdate()
	return r
}

//Sort 获取根据访问速率排序的服务列表
func (s *Sorter) Sort(names ...string) []string {

	//初始前不进行排序
	if len(s.sorted) == 0 {
		return names
	}

	//将列表存入map用于后续检查名称服务器是否存在,空间换时间
	dist := make(map[string]bool)
	for _, v := range names {
		dist[v] = false
	}

	atomic.CompareAndSwapInt32(&s.index, 2<<20, 1)
	index := atomic.AddInt32(&s.index, 1)
	//锁定与解锁
	s.lock.RLock()
	defer s.lock.RUnlock()

	//以排序好的服务器控制顺序
	sorted := make([]string, 0, len(names))
	current := s.sorted[int(index%int32(len(s.sorted)))]
	sorted = append(sorted, current)
	for _, v := range s.sorted {
		if v != current {
			if _, ok := dist[v]; ok {
				sorted = append(sorted, v)
				dist[v] = true
			}
		}

	}

	//未列入排序列表的，直接加入返回列表
	for k, v := range dist {
		if !v {
			sorted = append(sorted, k)
		}
	}
	return sorted
}

//UpdateRTT 更新请求时长
func (s *Sorter) UpdateRTT(name string, t time.Duration) {
	ok, rtt, _ := s.rtts.SetIfAbsentCb(name, func(...interface{}) (interface{}, error) {
		return &NameRTT{name: name, maxRequest: 1, avgRTT: int64(t)}, nil
	})
	nrtt := rtt.(*NameRTT)
	if !ok {
		nrtt.Update(int64(t))
	}
}

func (s *Sorter) loopUpdate() {
	for {
		select {
		case <-s.closeCh:
			return
		case <-time.After(time.Minute * 5): //每5分钟更新一次列表
			s.updateList()
		}
	}

}

func (s *Sorter) updateList() {
	items := s.rtts.Items()
	list := make(RTTS, 0, len(items))
	for _, v := range items {
		list = append(list, v.(*NameRTT))
	}
	sort.Sort(list)
	s.lock.Lock()
	defer s.lock.Unlock()
	s.sorted = list.ToList()
}

//Close 关闭服务
func (s *Sorter) Close() {
	close(s.closeCh)
}
