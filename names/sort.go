package names

import (
	"sync/atomic"
	"time"

	"github.com/micro-plat/lib4go/concurrent/cmap"
)

type NameRTT struct {
	name       string
	maxRequest int64
	avgRTT     int64
}

//Update 更新请求时间
func (n *NameRTT) Update(t int64) {
	max := atomic.AddInt64(&n.maxRequest, 1)
	n.avgRTT = (n.avgRTT*(max-1) + t) / max
}

type Sorter struct {
	sorted []string
	rtts   cmap.ConcurrentMap
}

func newSorter() *Sorter {
	return &Sorter{
		sorted: make([]string, 0, 1),
		rtts:   cmap.New(3),
	}
}

//Sort 获取根据访问速率排序的服务列表
func (s *Sorter) Sort(names ...string) []string {
	if len(s.sorted) == 0 {
		return names
	}
	return names
}

//UpdateRTT 更新请求时长
func (s *Sorter) UpdateRTT(name string, t time.Duration) {
	ok, rtt := s.rtts.SetIfAbsent(name, func(...interface{}) (interface{}, error) {
		return &NameRTT{name: name, maxRequest: 1, avgRTT: int64(t)}, nil
	})
	nrtt := rtt.(*NameRTT)
	if !ok {
		nrtt.Update(int64(t))
	}
}
func (s *Sorter) updateList() {
	items := s.rtts.Items()
	list := make([]*NameRTT, 0, len(items))
	for _, v := range items {
		list = append(list, v.(*NameRTT))
	}
}
