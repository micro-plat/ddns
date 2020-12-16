package names

import (
	"sort"
	"sync"
	"time"

	"github.com/sparrc/go-ping"
)

type pingStat struct {
	server string
	stats  *ping.Statistics
	err    error
}

type pingStats []*pingStat

func (s pingStats) Len() int { return len(s) }

func (s pingStats) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s pingStats) Less(i, j int) bool {
	left := time.Minute
	right := time.Minute
	if s[i].err == nil {
		left = s[i].stats.AvgRtt
	}
	if s[j].err == nil {
		right = s[j].stats.AvgRtt
	}
	return left < right
}

//Sort 获取根据访问速率排序的服务列表
func Sort(server ...string) ([]string, error) {
	lst := make([]*pingStat, 0, len(server))
	var wg sync.WaitGroup
	p := func(s string) {
		defer wg.Done()
		pinger, err := ping.NewPinger(s)
		if err != nil {
			lst = append(lst, &pingStat{server: s, err: err})
			return
		}
		pinger.Timeout = time.Second
		pinger.SetPrivileged(true)
		pinger.Count = 3
		pinger.Run() // blocks until finished
		stats := pinger.Statistics()
		lst = append(lst, &pingStat{server: s, stats: stats})
	}

	for _, s := range server {
		wg.Add(1)
		go p(s)
	}
	wg.Wait()
	sort.Sort(pingStats(lst))
	nlst := make([]string, 0, len(server))
	for _, stat := range lst {
		nlst = append(nlst, stat.server)
	}
	return nlst, nil
}
