package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/micro-plat/hydra/hydra"
)

var app = hydra.NewApp(
	hydra.WithPlatName("ddns"),
	hydra.WithSystemName("ddns"),
	hydra.WithServerTypes("api-cron"),
	hydra.WithClusterName("dns"))

func main() {
	app.Start()
}
