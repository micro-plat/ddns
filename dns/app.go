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

	vchttp "github.com/micro-plat/hydra/conf/vars/http"
	c "github.com/micro-plat/hydra/hydra/servers/cron"
)

//App dns应用程序
var App = hydra.NewApp(
	hydra.WithPlatName("ddns", "DDNS"),
	hydra.WithSystemName("ddns", "域名解析"),
	hydra.WithUsage("DNS服务"),
	hydra.WithServerTypes(DDNS, c.CRON, http.Web),
	hydra.WithClusterName("prod"),
)

func init() {

	//初始化服务器配置
	hydra.Conf.Custom(DDNS, conf.New(conf.WithTimeout(10, 10))).
		Sub(conf.TypeNodeName, conf.NewNames("8.8.8.8", "114.114.114.114"))

	hydra.Conf.Web(":80", api.WithTimeout(300, 300), api.WithDNS("ddns.com")).
		Static(static.WithArchive(web.Archive)).
		Header(header.WithCrossDomain())

	hydra.Conf.CRON(cron.WithMasterSlave())
	hydra.Conf.Vars().HTTP("http", vchttp.WithRequestTimeout(30), vchttp.WithConnTimeout(30))

	//注册服务
	App.Micro("/ddns/*", services.NewDdnsHandler())
	App.Micro("/github/ip/*", services.NewGithubHandler())
	App.CRON("/github/ip/*", services.NewGithubHandler(), "@midnight", "@now")

}
