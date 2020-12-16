package names

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

	//从本地缓存获取
	if names := DefRegistry.Lookup(); len(names) > 0 {
		return names
	}
	if names := l.r.Lookup(); len(names) > 0 {
		return names
	}
	return defNames

}

//Close 关闭服务
func (l *Names) Close() {
	l.r.Close()
}
