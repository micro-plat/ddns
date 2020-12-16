package names

import "github.com/micro-plat/ddns/pkgs"

var defNames = []string{"8.8.8.8:53"}

type Names struct {
	r *ResolveConf
}

//New 构建本地服务
func New() (*Names, error) {
	l := &Names{
		r: NewResolveConf(),
	}
	if err := l.r.Start(); err != nil {
		return nil, err
	}
	return l, nil
}

//Lookup 根据域名查询
func (l *Names) Lookup() []string {

	//从注册中心拉取
	names := DefRegistry.Lookup()

	//从本地拉取
	names = append(names, l.r.Lookup()...)

	//追加默认服务
	names = append(names, defNames...)

	//服务去重
	return pkgs.Distinct(names)
}

//Close 关闭服务
func (l *Names) Close() {
	if l.r != nil {
		l.r.Close()
	}
}
