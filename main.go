package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/micro-plat/hydra/hydra"
)

type ddns struct {
	*hydra.MicroApp
}

func main() {
	app := &ddns{
		hydra.NewApp(
			hydra.WithPlatName("ddns"),
			hydra.WithSystemName("ddns"),
			hydra.WithServerTypes("api-cron"),
			hydra.WithClusterName("dns")),
	}

	app.init()
	app.Start()
}
