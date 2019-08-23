module github.com/micro-plat/ddns

go 1.12

require (
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gin-gonic/gin v1.4.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/micro-plat/hydra v0.10.13
	github.com/micro-plat/lib4go v0.1.7
	github.com/miekg/dns v1.1.16
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/zkfy/log v0.0.0-20180312054228-b2704c3ef896
)

//replace github.com/micro-plat/hydra => ../../../github.com/micro-plat/hydra
