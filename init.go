package main

import (
	"sync"

	"github.com/micro-plat/ddns/dns"
	"github.com/micro-plat/ddns/services"
	"github.com/micro-plat/hydra/component"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func init() {

	var server *dns.Server
	var once sync.Once
	app.Initializing(func(c component.IContainer) error {
		once.Do(func() {
			//初始化注册中查询服务
			server = dns.NewServer(c)
			if err := server.Start(); err != nil {
				panic(err)
			}
		})
		return nil
	})
	app.Closing(func(c component.IContainer) error {

		server.Shutdown()
		return nil
	})
	app.Micro("/ddns/request", services.NewDdnsHandler)
	app.CRON("/github/ip/check", services.NewGithubHandler)
}
