package main

import (
	"github.com/micro-plat/ddns/services"
	"github.com/micro-plat/hydra"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func init() {

	app.Micro("/ddns/request", services.NewDdnsHandler())
	app.Micro("/github/ip/check", services.NewGithubHandler())
	app.CRON("/github/ip/check", services.NewGithubHandler(), "@every 50s")

	hydra.Conf.OnReady(func() {
		hydra.CRON.Add("@now", "/github/ip/check")
	})

}
