package query

import (
	"fmt"
	"net"

	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/lib4go/logger"
)

//IQuery 查询服务
type IQuery interface {
	Lookup(name string) []net.IP
	Close() error
}

//localQueries 本地查询服务
var localQueries = make(map[string]IQuery)

//Register 注册查询服务
func Register(name string, query IQuery) {
	if _, ok := localQueries[name]; ok {
		panic(fmt.Sprintf("%s:重复注册查询服务", name))
	}
	localQueries[name] = query
}
func Lookup(name string) []net.IP {
	for _, q := range localQueries {
		if lst := q.Lookup(name); len(lst) > 0 {
			return lst
		}
	}
	return nil
}

//Start 启动查询注册
func Start(c component.IContainer, log logger.ILogger) error {
	registry := NewRegistry(c.GetRegistry(), log)
	if err := registry.Start(); err != nil {
		return err
	}
	Register("register", registry)
	hosts := NewHosts(log)
	if err := hosts.Start(); err != nil {
		return err
	}
	Register("hosts", hosts)
	return nil
}

//Close 关闭服务
func Close() error {
	for _, q := range localQueries {
		if err := q.Close(); err != nil {
			return err
		}
	}
	return nil
}
