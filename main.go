package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/micro-plat/ddns/dns"
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/global/compatible"
)

func main() {
	if err := compatible.CheckPrivileges(); err != nil {
		hydra.G.Log().Error(err)
		return
	}
	dns.App.Start()
}
