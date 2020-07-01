package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/micro-plat/ddns/dns"
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/cron"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

var app = hydra.NewApp(
	hydra.WithPlatName("ddns"),
	hydra.WithServerTypes(dns.DDNS, http.API, cron.CRON))

func main() {
	app.Start()
}
