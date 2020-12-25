package names

import (
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/micro-plat/lib4go/concurrent/cmap"
)

type Sorter struct {
	fastList []string
	slowList []string
	rtts     cmap.ConcurrentMap
	lock     sync.RWMutex
	closeCh  chan struct{}
	maxRTT   int64
	index    int32
}

func newSorter() *Sorter {
	r := &Sorter{
		fastList: make([]string, 0, 1),
		slowList: make([]string, 0, 1),
		rtts:     cmap.New(3),
		maxRTT:   int64(time.Millisecond * 30),
		closeCh:  make(chan struct{}),
	}
	go r.loopUpdate()
	return r
}

//Sort 获取根据访问速率排序的服务列表
func (s *Sorter) Sort(names ...string) []string {

	//初始前不进行排序
	if len(s.fastList) == 0 {
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

	//较快的DNS优先使用
	sorted := make([]string, 0, len(names))
	cindex := int(index % int32(len(s.fastList)))
	current := s.fastList[cindex]
	sorted = append(sorted, current)

	//加入较快的DNS
	for i, v := range s.fastList {
		if i != cindex {
			if _, ok := dist[v]; ok {
				sorted = append(sorted, v)
				dist[v] = true
			}
		}
	}

	//加入较慢的DNS
	for _, v := range s.slowList {
		if _, ok := dist[v]; ok {
			sorted = append(sorted, v)
			dist[v] = true
		}
	}

	//加入未列入排序列表的DNS
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
		case <-time.After(time.Minute * 2): //每2分钟更新一次列表
			s.updateList()
		}
	}
}

func (s *Sorter) updateList() {
	items := s.rtts.Items()
	fastList := make(RTTS, 0, len(items))
	slowList := make(RTTS, 0, 1)
	for _, v := range items {
		rtt := v.(*NameRTT)
		if rtt.avgRTT < s.maxRTT {
			fastList = append(fastList, rtt)
		} else {
			slowList = append(slowList, rtt)
		}
	}
	sort.Sort(fastList)
	sort.Sort(slowList)
	s.lock.Lock()
	defer s.lock.Unlock()
	s.fastList = fastList.ToList()
	s.slowList = slowList.ToList()
}

//Close 关闭服务
func (s *Sorter) Close() {
	close(s.closeCh)
}
