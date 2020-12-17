package main

import (
	"time"

	"github.com/micro-plat/ddns/dns"
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/global/compatible"
)

func main() {
	if err := compatible.CheckPrivileges(); err != nil {
		hydra.G.Log().Error(err)
		time.Sleep(time.Second)
		return
	}
	dns.App.Start()
}
