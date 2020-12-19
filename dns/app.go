package dns

import (
	"github.com/micro-plat/ddns/conf"
	"github.com/micro-plat/ddns/services"
	"github.com/micro-plat/ddns/web"
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/hydra/servers/http"

	c "github.com/micro-plat/hydra/hydra/servers/cron"
)

//App dns应用程序
var App = hydra.NewApp(
	hydra.WithPlatName("ddns"),
	hydra.WithSystemName("ddnsserver"),
	hydra.WithUsage("DNS服务"),
	hydra.WithServerTypes(DDNS, http.API, c.CRON, http.Web),
	hydra.WithClusterName("dns-1.2"),
	hydra.WithRunFlag("dnsroot", "DNS的跟节点名称"),
	hydra.WithRegistry("zk://192.168.0.101"),
)

func init() {

	//初始化服务器配置
	hydra.Conf.Custom(DDNS, conf.New(conf.WithTimeout(10, 10), conf.WithOnlyUseRemote())).
		Sub(conf.TypeNodeName, conf.NewNames("8.8.8.8", "114.114.114.114"))
	hydra.Conf.Web(":80", api.WithDNS("ddns.com")).Static(static.WithArchive(web.Archive))
	hydra.Conf.API(":8081", api.WithDNS("ddns.com"), api.WithTimeout(300, 300)).
		Header(header.WithCrossDomain())
	hydra.Conf.CRON(cron.WithMasterSlave())

	//注册服务
	App.Micro("/ddns/*", services.NewDdnsHandler())
	App.Micro("/github/ip/*", services.NewGithubHandler())
	App.CRON("/github/ip/*", services.NewGithubHandler())
	hydra.CRON.Add("@midnight", "/github/ip/request")
	hydra.CRON.Add("@now", "/github/ip/request")

}
