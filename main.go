package main

import (
	"github.com/micro-plat/ddns/dns"
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/global/compatible"
	"github.com/micro-plat/lib4go/logger"
)

func main() {
	defer logger.Close()
	if err := compatible.CheckPrivileges(); err != nil {
		hydra.G.Log().Error(err)
		return
	}
	dns.App.Start()
}
