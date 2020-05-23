package main

import (
	"github.com/micro-plat/ddns/services"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func init() {

	app.Micro("/ddns/request", services.NewDdnsHandler())
	app.CRON("/github/ip/check", services.NewGithubHandler())
}
