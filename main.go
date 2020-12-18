package main

import (
	"github.com/micro-plat/ddns/dns"
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/global/compatible"
	"github.com/micro-plat/lib4go/logger"
)

func main() {
	if err := compatible.CheckPrivileges(); err != nil {
		defer logger.Close()
		hydra.G.Log().Error(err)
		return
	}
	dns.App.Start()
}
