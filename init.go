package main

import (
	"github.com/micro-plat/ddns/services"
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/cron"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func init() {

	app.Micro("/ddns", services.NewDdnsHandler())
	app.Micro("/github/ip/check", services.NewGithubHandler())
	app.CRON("/github/ip/check", services.NewGithubHandler(), "@every 5m")
	hydra.OnReady(func() {
		hydra.Conf.API(":8081", api.WithDNS("www.ddns.com"))
		hydra.Conf.CRON(cron.WithMasterSlave())
		hydra.CRON.Add("@now", "/github/ip/check")
	})
}
