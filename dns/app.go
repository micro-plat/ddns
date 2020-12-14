package dns

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/cron"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

var App = hydra.NewApp(
	hydra.WithPlatName("ddns"),
	hydra.WithUsage("DNS服务"),
	hydra.WithServerTypes(DDNS, http.API, cron.CRON),
	hydra.WithClusterName("dns-1.2"))
