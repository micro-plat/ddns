package dns

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/api"
	"github.com/micro-plat/hydra/conf/server/header"
)

func init() {
	hydra.Conf.API(":9090", api.WithTimeout(300, 300)).
		Header(header.WithCrossDomain())

	//_dns_bootstrap_js

}
