module github.com/micro-plat/ddns

go 1.15

replace github.com/micro-plat/hydra => ../../../github.com/micro-plat/hydra

replace github.com/micro-plat/lib4go => ../../../github.com/micro-plat/lib4go

require (
	github.com/PuerkitoBio/goquery v1.6.0
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/fsnotify/fsnotify v1.4.9
	github.com/micro-plat/hydra v1.0.2
	github.com/micro-plat/lib4go v1.0.9
	github.com/miekg/dns v1.1.35
	github.com/zkfy/go-cache v2.1.0+incompatible
	golang.org/x/sys v0.0.0-20201221093633-bc327ba9c2f0
	gopkg.in/fsnotify.v1 v1.4.7
)