package names

import "sync/atomic"

type NameRTT struct {
	name       string
	maxRequest int64
	avgRTT     int64
}

//Update 更新请求时间
func (n *NameRTT) Update(t int64) {
	atomic.CompareAndSwapInt64(&n.maxRequest, 1<<30, 1000)
	max := atomic.AddInt64(&n.maxRequest, 1)
	n.avgRTT = (n.avgRTT*(max-1) + t) / max
}

type RTTS []*NameRTT

func (s RTTS) Len() int { return len(s) }

func (s RTTS) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s RTTS) Less(i, j int) bool {
	left := s[i].avgRTT
	right := s[j].avgRTT
	return left < right
}
func (s RTTS) ToList() []string {
	list := make([]string, 0, len(s))
	for _, v := range s {
		list = append(list, v.name)
	}
	return list
}
