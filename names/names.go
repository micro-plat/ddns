package names

import (
	"time"

	"github.com/micro-plat/ddns/pkgs"
)

var defNames = []string{"8.8.8.8:53"}

type Names struct {
	r      *ResolveConf
	sorter *Sorter
}

//New 构建本地服务
func New() (*Names, error) {
	l := &Names{
		r:      NewResolveConf(),
		sorter: newSorter(),
	}
	if err := l.r.Start(); err != nil {
		return nil, err
	}
	return l, nil
}

//Lookup 获取可用的名称服务器
func (l *Names) Lookup() []string {

	//从注册中心拉取
	names := DefRegistry.Lookup()

	//从本地拉取
	names = append(names, l.r.Lookup()...)

	//追加默认服务
	names = append(names, defNames...)

	//服务去重,并排序
	return l.sorter.Sort(pkgs.Distinct(pkgs.Filte(names...))...)
}

//UpdateRTT 更新请求时长
func (l *Names) UpdateRTT(name string, t time.Duration) {
	l.sorter.UpdateRTT(name, t)
}

//Close 关闭服务
func (l *Names) Close() {
	if l.r != nil {
		l.r.Close()
	}
	if l.sorter != nil {
		l.sorter.Close()
	}
}
