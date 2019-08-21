package main

import "github.com/micro-plat/hydra/hydra"
import _ "github.com/go-sql-driver/mysql"

type ddns struct {
	*hydra.MicroApp
}

func main() {
	app := &ddns{
		hydra.NewApp(
			hydra.WithPlatName("ddns"),
			hydra.WithSystemName("ddns"),
			hydra.WithServerTypes("api")),
	}

	app.init()
	app.Start()
}
