package query

import (
	"fmt"
	"net"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
)

//IQuery 查询服务
type IQuery interface {
	Lookup(name string) []net.IP
	Close() error
}

//localQueries 本地查询服务
var localQueries = make(map[string]IQuery)
var queries = []IQuery{}

//Register 注册查询服务
func Register(name string, query IQuery) {
	if _, ok := localQueries[name]; ok {
		panic(fmt.Sprintf("%s:重复注册查询服务", name))
	}
	localQueries[name] = query
	queries = append(queries, query)
}
func Lookup(name string) []net.IP {
	for _, q := range queries {
		if lst := q.Lookup(name); len(lst) > 0 {
			return lst
		}
	}
	return nil
}

//Start 启动查询注册
func Start(log logger.ILogger) error {
	r, err := registry.NewRegistry(hydra.G.RegistryAddr, log)
	if err != nil {
		return err
	}
	registry := NewRegistry(r, log)
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
