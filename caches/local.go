package caches

import (
	"fmt"
	"net"
)

//localQueries 本地查询服务
var qmaps = make(map[string]IQuery)
var queries = []IQuery{}

//IQuery 查询服务
type IQuery interface {
	Lookup(name string) []net.IP
}

//Register 注册查询服务
func Register(name string, query IQuery) {
	if _, ok := qmaps[name]; ok {
		panic(fmt.Sprintf("%s:重复注册查询服务", name))
	}
	qmaps[name] = query
	queries = append(queries, query)
}
