package dns

import (
	"github.com/micro-plat/ddns/conf"
	"github.com/micro-plat/ddns/services"
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/static"
)

func init() {

	//初始化服务器配置
	hydra.Conf.Custom(DDNS, conf.New(conf.WithTimeout(10, 10))).
		Sub(conf.TypeNodeName, conf.NewNames("8.8.8.8"))
	hydra.Conf.Web(":80", api.WithDNS("ddns.com")).Static(static.WithArchive(services.Archive))
	hydra.Conf.API(":8081", api.WithDNS("ddns.com"))
	hydra.Conf.CRON(cron.WithMasterSlave())

	//注册服务
	App.Micro("/ddns/*", services.NewDdnsHandler())
	App.Micro("/github/ip/*", services.NewGithubHandler())
	App.CRON("/github/ip/*", services.NewGithubHandler(), "@midnight", "@now")

}
