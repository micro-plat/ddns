package main

import (
	"github.com/micro-plat/ddns/dns"
	"github.com/micro-plat/hydra/component"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func (app *ddns) init() {

	app.install()
	app.handling()

	var server *dns.Server
	app.Initializing(func(c component.IContainer) error {
		//初始化注册中查询服务
		server = dns.NewServer(c)
		return server.Start()
	})
	app.Closing(func(c component.IContainer) error {

		server.Shutdown()
		return nil
	})
}
