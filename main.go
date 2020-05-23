package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/micro-plat/ddns/dns"
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/servers/http"
)

var app = hydra.NewApp(
	hydra.WithPlatName("ddns"),
	hydra.WithSystemName("ddns"),
	hydra.WithServerTypes(http.API, dns.DDNS),
	hydra.WithClusterName("dns"))

func main() {
	app.Start()
}
