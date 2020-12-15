package main

import (
	"github.com/micro-plat/ddns/dns"
	"github.com/micro-plat/ddns/dns/conf"
	"github.com/micro-plat/ddns/services"
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/cron"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func init() {
	hydra.Conf.Custom(dns.DDNS, conf.New(conf.WithDisable(), conf.WithTimeout(10, 10))).Sub(conf.TypeNodeName, conf.NewDnss("114.114.114.114", "8.8.8.8"))
	dns.App.Micro("/ddns", services.NewDdnsHandler())
	dns.App.Micro("/github/ip/check", services.NewGithubHandler())
	dns.App.CRON("/github/ip/check", services.NewGithubHandler(), "@midnight")
	hydra.OnReady(func() {
		hydra.Conf.API(":8081", api.WithDNS("www.ddns.com"))
		hydra.Conf.CRON(cron.WithMasterSlave())
		hydra.CRON.Add("@now", "/github/ip/check")
	})
}
