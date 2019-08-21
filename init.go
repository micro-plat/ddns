package main

import (
	"github.com/micro-plat/ddns/services"
	"github.com/micro-plat/hydra/component"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func (app *ddns) init() {

	app.install()
	app.handling()
	dns := services.NewServer()
	app.Initializing(func(c component.IContainer) error {
		dns.Start()
		return nil
	})
}
