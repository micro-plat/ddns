package main

import (
	"github.com/micro-plat/ddns/conf"
	"github.com/micro-plat/ddns/dns"
	"github.com/micro-plat/ddns/services"
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/static"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func init() {
	hydra.Conf.Custom(dns.DDNS, conf.New(conf.WithTimeout(10, 10))).
		Sub(conf.TypeNodeName, conf.NewNames("114.114.114.114", "8.8.8.8"))
	dns.App.Micro("/ddns", services.NewDdnsHandler())
	dns.App.Micro("/github/ip/*", services.NewGithubHandler())
	dns.App.CRON("/github/ip/*", services.NewGithubHandler(), "@midnight")

	hydra.Conf.API(":8081").Static(static.WithArchive("./assert/assert.so"), static.WithRoot("./assert"), static.WithEnable())
	hydra.Conf.CRON(cron.WithMasterSlave())
	hydra.CRON.Add("@now", "/github/ip/request")

}
